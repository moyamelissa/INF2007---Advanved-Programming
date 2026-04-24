package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// ========== Tests unitaires ==========

// TestCountWordsHTMLSimple vérifie le comptage de mots dans un HTML simple.
func TestCountWordsHTMLSimple(t *testing.T) {
	htmlContent := `<html><body><p>Bonjour le monde</p></body></html>`
	result := countWordsHTML(htmlContent)
	if result != 3 {
		t.Errorf("attendu 3 mots, obtenu %d", result)
	}
}

// TestCountWordsHTMLMultipleTags vérifie le comptage à travers plusieurs balises.
func TestCountWordsHTMLMultipleTags(t *testing.T) {
	htmlContent := `<html><body>
		<h1>Titre principal</h1>
		<p>Premier paragraphe avec cinq mots ici.</p>
		<p>Deuxième paragraphe.</p>
	</body></html>`
	result := countWordsHTML(htmlContent)
	// "Titre principal" = 2, "Premier paragraphe avec cinq mots ici." = 6, "Deuxième paragraphe." = 2
	if result != 10 {
		t.Errorf("attendu 10 mots, obtenu %d", result)
	}
}

// TestCountWordsHTMLIgnoreScript vérifie que le contenu de <script> est ignoré.
func TestCountWordsHTMLIgnoreScript(t *testing.T) {
	htmlContent := `<html><body>
		<p>Texte visible</p>
		<script>var x = "code ignoré ici";</script>
		<p>Autre texte</p>
	</body></html>`
	result := countWordsHTML(htmlContent)
	// "Texte visible" = 2, "Autre texte" = 2
	if result != 4 {
		t.Errorf("attendu 4 mots (script ignoré), obtenu %d", result)
	}
}

// TestCountWordsHTMLIgnoreStyle vérifie que le contenu de <style> est ignoré.
func TestCountWordsHTMLIgnoreStyle(t *testing.T) {
	htmlContent := `<html><head><style>body { color: red; }</style></head>
		<body><p>Seul texte</p></body></html>`
	result := countWordsHTML(htmlContent)
	if result != 2 {
		t.Errorf("attendu 2 mots (style ignoré), obtenu %d", result)
	}
}

// TestCountWordsHTMLEmpty vérifie le comptage pour un HTML vide.
func TestCountWordsHTMLEmpty(t *testing.T) {
	result := countWordsHTML("")
	if result != 0 {
		t.Errorf("attendu 0 mots pour HTML vide, obtenu %d", result)
	}
}

// TestCountWordsHTMLOnlyTags vérifie le comptage pour un HTML avec balises mais sans texte.
func TestCountWordsHTMLOnlyTags(t *testing.T) {
	htmlContent := `<html><body><div><span></span></div></body></html>`
	result := countWordsHTML(htmlContent)
	if result != 0 {
		t.Errorf("attendu 0 mots pour HTML sans texte, obtenu %d", result)
	}
}

// TestFetchPageSuccess vérifie la récupération d'une page via un serveur de test.
func TestFetchPageSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<html><body><p>Hello World</p></body></html>`)
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	content, err := fetchPage(server.URL, client)
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}
	if !strings.Contains(content, "Hello World") {
		t.Errorf("contenu attendu non trouvé dans la réponse")
	}
}

// TestFetchPageInvalidURL vérifie la gestion d'erreur pour une URL invalide.
func TestFetchPageInvalidURL(t *testing.T) {
	client := &http.Client{Timeout: 2 * time.Second}
	_, err := fetchPage("http://url-invalide-qui-nexiste-pas.xyz", client)
	if err == nil {
		t.Fatal("une erreur était attendue pour une URL invalide, mais aucune erreur retournée")
	}
}

// TestFetchPageTimeout vérifie la gestion du délai d'expiration.
func TestFetchPageTimeout(t *testing.T) {
	// Serveur qui ne répond jamais
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	}))
	defer server.Close()

	client := &http.Client{Timeout: 1 * time.Second}
	_, err := fetchPage(server.URL, client)
	if err == nil {
		t.Fatal("une erreur de timeout était attendue, mais aucune erreur retournée")
	}
}

// TestFetchPage404 vérifie la gestion d'un code HTTP 404.
func TestFetchPage404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	_, err := fetchPage(server.URL, client)
	if err == nil {
		t.Fatal("une erreur était attendue pour HTTP 404, mais aucune erreur retournée")
	}
}

// TestCheckRobotsAllowed vérifie que robots.txt est respecté.
func TestCheckRobotsAllowed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			fmt.Fprint(w, "User-agent: *\nDisallow: /private/\nAllow: /public/\n")
			return
		}
		fmt.Fprint(w, "<html><body>OK</body></html>")
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	// /public/ devrait être autorisé
	if !checkRobotsAllowed(server.URL+"/public/page", client) {
		t.Error("/public/page devrait être autorisé")
	}

	// /private/ devrait être interdit
	if checkRobotsAllowed(server.URL+"/private/secret", client) {
		t.Error("/private/secret devrait être interdit par robots.txt")
	}

	// / (racine) devrait être autorisé
	if !checkRobotsAllowed(server.URL+"/", client) {
		t.Error("/ devrait être autorisé")
	}
}

// TestCheckRobotsNoFile vérifie qu'on autorise si robots.txt n'existe pas.
func TestCheckRobotsNoFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, "<html><body>OK</body></html>")
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	if !checkRobotsAllowed(server.URL+"/anything", client) {
		t.Error("devrait être autorisé quand robots.txt n'existe pas")
	}
}

// TestCrawlURLsIntegration teste le crawl complet avec un serveur local.
func TestCrawlURLsIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			fmt.Fprint(w, "User-agent: *\nAllow: /\n")
			return
		}
		w.Header().Set("Content-Type", "text/html")
		switch r.URL.Path {
		case "/page1":
			fmt.Fprint(w, `<html><body><p>Un deux trois</p></body></html>`)
		case "/page2":
			fmt.Fprint(w, `<html><body><p>Quatre cinq</p></body></html>`)
		default:
			fmt.Fprint(w, `<html><body><p>Mot</p></body></html>`)
		}
	}))
	defer server.Close()

	urls := []string{server.URL + "/page1", server.URL + "/page2"}
	results, total, errs := CrawlURLs(urls, 4)

	if len(errs) > 0 {
		t.Fatalf("erreurs inattendues : %v", errs)
	}

	if results[server.URL+"/page1"] != 3 {
		t.Errorf("page1 : attendu 3 mots, obtenu %d", results[server.URL+"/page1"])
	}
	if results[server.URL+"/page2"] != 2 {
		t.Errorf("page2 : attendu 2 mots, obtenu %d", results[server.URL+"/page2"])
	}
	if total != 5 {
		t.Errorf("total : attendu 5 mots, obtenu %d", total)
	}
}

// TestCrawlURLsRobotsBlocked vérifie que le crawl respecte robots.txt.
func TestCrawlURLsRobotsBlocked(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			fmt.Fprint(w, "User-agent: *\nDisallow: /blocked\n")
			return
		}
		fmt.Fprint(w, `<html><body><p>Contenu secret</p></body></html>`)
	}))
	defer server.Close()

	urls := []string{server.URL + "/blocked"}
	results, total, errs := CrawlURLs(urls, 2)

	if len(errs) == 0 {
		t.Error("une erreur était attendue pour URL bloquée par robots.txt")
	}
	if total != 0 {
		t.Errorf("total : attendu 0 mots pour URL bloquée, obtenu %d", total)
	}
	if len(results) != 0 {
		t.Errorf("aucun résultat attendu pour URL bloquée, obtenu %d", len(results))
	}
}

// ========== Benchmarks ==========

// setupBenchServer crée un serveur de test avec N pages.
func setupBenchServer(nPages int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			fmt.Fprint(w, "User-agent: *\nAllow: /\n")
			return
		}
		w.Header().Set("Content-Type", "text/html")
		// Générer une page avec ~100 mots
		var sb strings.Builder
		sb.WriteString("<html><body>")
		for i := 0; i < 10; i++ {
			sb.WriteString("<p>Lorem ipsum dolor sit amet consectetur adipiscing elit sed do</p>")
		}
		sb.WriteString("</body></html>")
		fmt.Fprint(w, sb.String())
	}))
}

// BenchmarkCrawlGoroutines compare les performances avec 1, 2, 4 et 8 goroutines.
func BenchmarkCrawlGoroutines(b *testing.B) {
	server := setupBenchServer(8)
	defer server.Close()

	// Créer 8 URLs pointant vers le serveur de test
	urls := make([]string, 8)
	for i := range urls {
		urls[i] = fmt.Sprintf("%s/page%d", server.URL, i)
	}

	goroutineCounts := []int{1, 2, 4, 8}

	for _, g := range goroutineCounts {
		b.Run(fmt.Sprintf("%d_goroutines", g), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				CrawlURLs(urls, g)
			}
		})
	}
}

// BenchmarkCountWordsHTML mesure la performance du parsing HTML.
func BenchmarkCountWordsHTML(b *testing.B) {
	var sb strings.Builder
	sb.WriteString("<html><body>")
	for i := 0; i < 100; i++ {
		sb.WriteString("<p>Lorem ipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labore et dolore magna aliqua</p>")
	}
	sb.WriteString("<script>var x = 'ignored code here';</script>")
	sb.WriteString("</body></html>")
	htmlContent := sb.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		countWordsHTML(htmlContent)
	}
}

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/temoto/robotstxt"
	"golang.org/x/net/html"
)

// CrawlResult contient le résultat de l'exploration d'une URL.
type CrawlResult struct {
	URL       string
	WordCount int
	Err       error
}

// newHTTPClient crée un client HTTP avec un délai d'expiration de 10 secondes.
func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

// checkRobotsAllowed vérifie si le chemin d'une URL est autorisé par robots.txt.
// Retourne true si l'exploration est permise, false sinon.
// En cas d'erreur de récupération de robots.txt, on autorise l'exploration par défaut
// (comportement standard des robots d'exploration).
func checkRobotsAllowed(targetURL string, client *http.Client) bool {
	parsed, err := url.Parse(targetURL)
	if err != nil {
		return false
	}

	robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsed.Scheme, parsed.Host)

	resp, err := client.Get(robotsURL)
	if err != nil {
		// Si robots.txt n'est pas accessible, on autorise par défaut
		return true
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Si robots.txt n'existe pas (404, etc.), tout est autorisé
		return true
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return true
	}

	robots, err := robotstxt.FromBytes(body)
	if err != nil {
		return true
	}

	// Vérifier les règles pour User-agent "*" (robot générique)
	group := robots.FindGroup("*")
	return group.Test(parsed.Path)
}

// fetchPage récupère le contenu HTML d'une URL.
func fetchPage(targetURL string, client *http.Client) (string, error) {
	resp, err := client.Get(targetURL)
	if err != nil {
		return "", fmt.Errorf("échec de la requête pour %s : %w", targetURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("statut HTTP %d pour %s", resp.StatusCode, targetURL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("échec de lecture du corps pour %s : %w", targetURL, err)
	}

	return string(body), nil
}

// countWordsHTML analyse le contenu HTML et compte les mots visibles (texte hors
// balises <script>, <style>, etc.). Utilise le tokenizer golang.org/x/net/html
// pour un parsing robuste du HTML.
func countWordsHTML(htmlContent string) int {
	tokenizer := html.NewTokenizer(strings.NewReader(htmlContent))
	count := 0
	// skipTags contient les balises dont le contenu textuel doit être ignoré
	skipTags := map[string]bool{
		"script":   true,
		"style":    true,
		"noscript": true,
	}
	skip := false

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return count
		case html.StartTagToken:
			tn, _ := tokenizer.TagName()
			tagName := string(tn)
			if skipTags[tagName] {
				skip = true
			}
		case html.EndTagToken:
			tn, _ := tokenizer.TagName()
			tagName := string(tn)
			if skipTags[tagName] {
				skip = false
			}
		case html.TextToken:
			if !skip {
				text := strings.TrimSpace(tokenizer.Token().Data)
				if text != "" {
					words := strings.Fields(text)
					count += len(words)
				}
			}
		}
	}
}

// crawlURL explore une URL unique : vérifie robots.txt, récupère la page,
// et compte les mots. Envoie le résultat sur le canal ch.
func crawlURL(targetURL string, client *http.Client, ch chan<- CrawlResult) {
	// Vérifier robots.txt avant d'explorer
	if !checkRobotsAllowed(targetURL, client) {
		ch <- CrawlResult{
			URL: targetURL,
			Err: fmt.Errorf("exploration interdite par robots.txt pour %s", targetURL),
		}
		return
	}

	content, err := fetchPage(targetURL, client)
	if err != nil {
		ch <- CrawlResult{URL: targetURL, Err: err}
		return
	}

	wordCount := countWordsHTML(content)
	ch <- CrawlResult{URL: targetURL, WordCount: wordCount}
}

// CrawlURLs explore une liste d'URL de manière concurrente en limitant le nombre
// de goroutines actives simultanément à maxGoroutines.
//
// Paramètres :
//   - urls : liste des URL à explorer.
//   - maxGoroutines : nombre maximal de goroutines concurrentes (1, 2, 4, 8, etc.).
//
// Retour :
//   - map[string]int : nombre de mots par URL.
//   - int : total global de mots.
//   - []error : liste des erreurs rencontrées.
func CrawlURLs(urls []string, maxGoroutines int) (map[string]int, int, []error) {
	if maxGoroutines <= 0 {
		maxGoroutines = 1
	}

	client := newHTTPClient()
	results := make(map[string]int)
	var totalWords int
	var mu sync.Mutex
	var errs []error

	ch := make(chan CrawlResult, len(urls))
	// Sémaphore pour limiter le nombre de goroutines concurrentes
	semaphore := make(chan struct{}, maxGoroutines)

	var wg sync.WaitGroup

	for _, u := range urls {
		wg.Add(1)
		go func(targetURL string) {
			defer wg.Done()
			// Acquérir le sémaphore (bloque si maxGoroutines goroutines sont actives)
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			crawlURL(targetURL, client, ch)
		}(u)
	}

	// Fermer le canal quand toutes les goroutines ont terminé
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collecter les résultats depuis le canal et agréger avec le mutex
	for result := range ch {
		mu.Lock()
		if result.Err != nil {
			errs = append(errs, result.Err)
		} else {
			results[result.URL] = result.WordCount
			totalWords += result.WordCount
		}
		mu.Unlock()
	}

	return results, totalWords, errs
}

func main() {
	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://go.dev",
		"https://en.wikipedia.org/wiki/Go_(programming_language)",
	}

	fmt.Println("=== Robot d'exploration Web concurrent ===")
	fmt.Printf("URLs à explorer : %d\n", len(urls))
	fmt.Printf("Goroutines max  : 8\n\n")

	start := time.Now()
	results, total, errs := CrawlURLs(urls, 8)
	elapsed := time.Since(start)

	for urlStr, count := range results {
		fmt.Printf("  %-60s : %d mots\n", urlStr, count)
	}

	if len(errs) > 0 {
		fmt.Printf("\nErreurs (%d) :\n", len(errs))
		for _, e := range errs {
			fmt.Printf("  - %v\n", e)
		}
	}

	fmt.Printf("\nTotal global    : %d mots\n", total)
	fmt.Printf("Temps d'exécution : %v\n", elapsed)
}

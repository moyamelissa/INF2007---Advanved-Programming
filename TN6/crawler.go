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

// politenessDelay est le délai de politesse entre les requêtes vers un même serveur.
// Modifiable dans les tests et bancs d'essai pour mesurer la vraie performance du crawl.
var politenessDelay = 100 * time.Millisecond

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

// fetchRobots récupère et analyse robots.txt pour un hôte donné.
// Retourne nil si robots.txt est inaccessible, introuvable ou invalide,
// ce qui signifie que l'exploration est autorisée par défaut.
func fetchRobots(scheme, host string, client *http.Client) *robotstxt.RobotsData {
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", scheme, host)
	resp, err := client.Get(robotsURL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	robots, err := robotstxt.FromBytes(body)
	if err != nil {
		return nil
	}
	return robots
}

// checkRobotsAllowed vérifie si le chemin d'une URL est autorisé par robots.txt.
// Retourne true si l'exploration est permise, false sinon.
// Si cache et cacheMu sont fournis (non nil), robots.txt est mis en cache par hôte
// pour éviter les requêtes répétées. En cas d'erreur, on autorise par défaut.
func checkRobotsAllowed(targetURL string, client *http.Client, cache map[string]*robotstxt.RobotsData, cacheMu *sync.RWMutex) bool {
	parsed, err := url.Parse(targetURL)
	if err != nil {
		return false
	}
	host := parsed.Host

	var robots *robotstxt.RobotsData
	if cache != nil && cacheMu != nil {
		cacheMu.RLock()
		r, found := cache[host]
		cacheMu.RUnlock()

		if !found {
			r = fetchRobots(parsed.Scheme, host, client)
			cacheMu.Lock()
			if _, alreadyCached := cache[host]; !alreadyCached {
				cache[host] = r
			} else {
				r = cache[host]
			}
			cacheMu.Unlock()
		}
		robots = r
	} else {
		robots = fetchRobots(parsed.Scheme, host, client)
	}

	if robots == nil {
		return true
	}
	return robots.FindGroup("*").Test(parsed.Path)
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
// Un délai de politesse configurable (politenessDelay) est appliqué après la
// vérification robots.txt pour limiter la fréquence d'exploration.
func crawlURL(targetURL string, client *http.Client, ch chan<- CrawlResult, cache map[string]*robotstxt.RobotsData, cacheMu *sync.RWMutex) {
	// Vérifier robots.txt avant d'explorer
	if !checkRobotsAllowed(targetURL, client, cache, cacheMu) {
		ch <- CrawlResult{
			URL: targetURL,
			Err: fmt.Errorf("exploration interdite par robots.txt pour %s", targetURL),
		}
		return
	}

	// Délai de politesse configurable : limite la fréquence des requêtes vers chaque serveur
	if politenessDelay > 0 {
		time.Sleep(politenessDelay)
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
	var errs []error

	// Cache robots.txt par hôte pour éviter les requêtes répétées
	robotsCache := make(map[string]*robotstxt.RobotsData)
	var cacheMu sync.RWMutex

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

			crawlURL(targetURL, client, ch, robotsCache, &cacheMu)
		}(u)
	}

	// Fermer le canal quand toutes les goroutines ont terminé
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collecter les résultats depuis le canal (goroutine unique : pas de mutex nécessaire)
	for result := range ch {
		if result.Err != nil {
			errs = append(errs, result.Err)
		} else {
			results[result.URL] = result.WordCount
			totalWords += result.WordCount
		}
	}

	return results, totalWords, errs
}

// run contient la logique principale du programme, extraite de main pour
// permettre les tests unitaires.
func run(urls []string, maxGoroutines int) {
	fmt.Println("=== Robot d'exploration Web concurrent ===")
	fmt.Printf("URLs à explorer : %d\n", len(urls))
	fmt.Printf("Goroutines max  : %d\n\n", maxGoroutines)

	start := time.Now()
	results, total, errs := CrawlURLs(urls, maxGoroutines)
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

// mainURLs contient la liste des URLs explorées par défaut.
// Variable exportée pour permettre l'injection dans les tests.
var mainURLs = []string{
	"https://www.google.com",
	"https://www.github.com",
	"https://go.dev",
	"https://en.wikipedia.org/wiki/Go_(programming_language)",
}

func main() {
	run(mainURLs, 8)
}

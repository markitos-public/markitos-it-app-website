package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"markitos-it-app-website/internal/faqs"
)

type PageModel struct {
	Title       string
	Description string
	Faqs        []FAQView
	Tags        []string
	Section     string
	Status      string
}

type FAQView struct {
	Number  int
	Title   string
	Content string
	Tags    []string
	TagText string
}

func parseTemplates() (*template.Template, error) {
	return template.ParseFiles(
		"index.html",
		"faqs.html",
		"articles.html",
		"videos.html",
		"git.html",
	)
}

func IndexHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		render(w, templates, "index.html", PageModel{
			Title:       "Markitos MDK | DevSecOps Kulture",
			Description: "Markitos MDK: artículos, FAQs, vídeos y recursos de Git sobre DevSecOps.",
		})
	}
}

func FaqsHandler(templates *template.Template, client *faqs.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/faqs" {
			http.NotFound(w, r)
			return
		}

		model := PageModel{
			Title:       "FAQs | Markitos MDK",
			Description: "Preguntas frecuentes sobre DevSecOps, cloud y herramientas de seguridad.",
		}

		items, err := client.List()
		if err != nil {
			log.Printf("list FAQs: %v", err)
			render(w, templates, "faqs.html", model)
			return
		}

		model.Faqs, model.Tags = buildFAQView(items)
		render(w, templates, "faqs.html", model)
	}
}

func buildFAQView(items []faqs.FAQ) ([]FAQView, []string) {
	views := make([]FAQView, 0, len(items))
	tags := make(map[string]struct{})

	for index, item := range items {
		views = append(views, FAQView{
			Number:  index + 1,
			Title:   item.Title,
			Content: item.Content,
			Tags:    item.Tags,
			TagText: strings.Join(item.Tags, " "),
		})
		for _, tag := range item.Tags {
			tags[tag] = struct{}{}
		}
	}

	uniqueTags := make([]string, 0, len(tags))
	for tag := range tags {
		uniqueTags = append(uniqueTags, tag)
	}
	sort.Strings(uniqueTags)

	return views, uniqueTags
}

func ArticlesHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/articles" {
			http.NotFound(w, r)
			return
		}

		render(w, templates, "articles.html", PageModel{Title: "Artículos | Markitos MDK", Description: "Artículos de Markitos MDK sobre DevSecOps, cloud y seguridad.", Section: "Artículos", Status: "Lecturas del taller"})
	}
}

func VideosHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/videos" {
			http.NotFound(w, r)
			return
		}

		render(w, templates, "videos.html", PageModel{Title: "Vídeos | Markitos MDK", Description: "Vídeos técnicos de Markitos MDK sobre DevSecOps y seguridad.", Section: "Vídeos", Status: "Sesiones del taller"})
	}
}

func GitHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/git" {
			http.NotFound(w, r)
			return
		}

		render(w, templates, "git.html", PageModel{Title: "Recursos de Git | Markitos MDK", Description: "Repositorios y recursos de Git de Markitos MDK.", Section: "Recursos de Git", Status: "Código del taller"})
	}
}

func render(w http.ResponseWriter, templates *template.Template, name string, model PageModel) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, name, model); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func main() {
	faqsEndpoint := os.Getenv("FAQS_API_ENDPOINT")
	if faqsEndpoint == "" {
		log.Fatal("FAQS_API_ENDPOINT must be set")
	}

	templates, err := parseTemplates()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexHandler(templates))
	mux.HandleFunc("/faqs", FaqsHandler(templates, faqs.NewClient(faqsEndpoint)))
	mux.HandleFunc("/articles", ArticlesHandler(templates))
	mux.HandleFunc("/videos", VideosHandler(templates))
	mux.HandleFunc("/git", GitHandler(templates))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           loggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Markitos MDK listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

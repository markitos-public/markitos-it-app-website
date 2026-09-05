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
}

type IndexModel struct {
	PageModel
	ResourceCount int
}

type FaqsModel struct {
	PageModel
	Faqs []FAQView
	Tags []string
}

type FAQView struct {
	Number  int
	Title   string
	Content string
	Tags    []string
	TagText string
}

type SectionModel struct {
	PageModel
	Section string
	Status  string
}

type App struct {
	templates *template.Template
	faqs      *faqs.Client
}

func NewApp(faqsClient *faqs.Client) (*App, error) {
	templates, err := template.ParseFiles(
		"index.html",
		"faqs.html",
		"articles.html",
		"videos.html",
		"git.html",
	)
	if err != nil {
		return nil, err
	}

	return &App{templates: templates, faqs: faqsClient}, nil
}

func (app *App) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	app.render(w, "index.html", IndexModel{
		PageModel: PageModel{
			Title:       "Markitos MDK | DevSecOps Kulture",
			Description: "Markitos MDK: artículos, FAQs, vídeos y recursos de Git sobre DevSecOps.",
		},
		ResourceCount: 4,
	})
}

func (app *App) FaqsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/faqs" {
		http.NotFound(w, r)
		return
	}

	items, err := app.faqs.List()
	if err != nil {
		log.Printf("list FAQs: %v", err)
		http.Error(w, "unable to load FAQs", http.StatusBadGateway)
		return
	}

	model := FaqsModel{
		PageModel: PageModel{
			Title:       "FAQs | Markitos MDK",
			Description: "Preguntas frecuentes sobre DevSecOps, cloud y herramientas de seguridad.",
		},
	}
	model.Faqs, model.Tags = buildFAQView(items)
	app.render(w, "faqs.html", model)
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

func (app *App) ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/articles" {
		http.NotFound(w, r)
		return
	}

	app.render(w, "articles.html", SectionModel{PageModel: PageModel{Title: "Artículos | Markitos MDK", Description: "Artículos de Markitos MDK sobre DevSecOps, cloud y seguridad."}, Section: "Artículos", Status: "Lecturas del taller"})
}

func (app *App) VideosHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/videos" {
		http.NotFound(w, r)
		return
	}

	app.render(w, "videos.html", SectionModel{PageModel: PageModel{Title: "Vídeos | Markitos MDK", Description: "Vídeos técnicos de Markitos MDK sobre DevSecOps y seguridad."}, Section: "Vídeos", Status: "Sesiones del taller"})
}

func (app *App) GitHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/git" {
		http.NotFound(w, r)
		return
	}

	app.render(w, "git.html", SectionModel{PageModel: PageModel{Title: "Recursos de Git | Markitos MDK", Description: "Repositorios y recursos de Git de Markitos MDK."}, Section: "Recursos de Git", Status: "Código del taller"})
}

func (app *App) render(w http.ResponseWriter, name string, model any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := app.templates.ExecuteTemplate(w, name, model); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func main() {
	faqsEndpoint := os.Getenv("FAQS_API_ENDPOINT")
	if faqsEndpoint == "" {
		log.Fatal("FAQS_API_ENDPOINT must be set")
	}

	app, err := NewApp(faqs.NewClient(faqsEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.IndexHandler)
	mux.HandleFunc("/faqs", app.FaqsHandler)
	mux.HandleFunc("/articles", app.ArticlesHandler)
	mux.HandleFunc("/videos", app.VideosHandler)
	mux.HandleFunc("/git", app.GitHandler)
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

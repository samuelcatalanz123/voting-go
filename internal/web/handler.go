// Package web sirve la app de votaciones.
package web

import (
	"embed"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/samuelcatalanz123/voting-go/internal/poll"
)

//go:embed templates/*.html static/*
var files embed.FS

// Handler sirve la web.
type Handler struct {
	tmpl  *template.Template
	store *poll.Store
}

// New crea el Handler con un store vacío.
func New() (*Handler, error) {
	tmpl, err := template.ParseFS(files, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Handler{tmpl: tmpl, store: poll.New()}, nil
}

// Routes monta las rutas.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(files))
	mux.HandleFunc("GET /{$}", h.home)
	mux.HandleFunc("POST /create", h.create)
	mux.HandleFunc("GET /poll/{id}", h.poll)
	mux.HandleFunc("POST /vote/{id}", h.vote)
	return mux
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.ExecuteTemplate(w, "home.html", nil); err != nil {
		http.Error(w, "error del servidor", http.StatusInternalServerError)
	}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	question := strings.TrimSpace(r.FormValue("question"))

	// Cada opción va en una línea; quitamos las vacías.
	var options []string
	for _, line := range strings.Split(r.FormValue("options"), "\n") {
		if o := strings.TrimSpace(line); o != "" {
			options = append(options, o)
		}
	}

	if question == "" || len(options) < 2 {
		http.Redirect(w, r, "/", http.StatusSeeOther) // falta pregunta o muy pocas opciones
		return
	}
	id := h.store.Create(question, options)
	http.Redirect(w, r, "/poll/"+id, http.StatusSeeOther)
}

// optionView es una opción con su porcentaje, para mostrarla.
type optionView struct {
	Index   int
	Text    string
	Votes   int
	Percent int
}

// pollView son los datos de la página de una encuesta.
type pollView struct {
	ID       string
	Question string
	URL      string
	Total    int
	Options  []optionView
}

func (h *Handler) poll(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, ok := h.store.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	total := p.Total()
	options := make([]optionView, len(p.Options))
	for i, o := range p.Options {
		percent := 0
		if total > 0 {
			percent = o.Votes * 100 / total
		}
		options[i] = optionView{Index: i, Text: o.Text, Votes: o.Votes, Percent: percent}
	}

	data := pollView{
		ID:       p.ID,
		Question: p.Question,
		URL:      "http://" + r.Host + "/poll/" + p.ID,
		Total:    total,
		Options:  options,
	}
	if err := h.tmpl.ExecuteTemplate(w, "poll.html", data); err != nil {
		http.Error(w, "error del servidor", http.StatusInternalServerError)
	}
}

func (h *Handler) vote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if idx, err := strconv.Atoi(r.FormValue("option")); err == nil {
		_ = h.store.Vote(id, idx)
	}
	http.Redirect(w, r, "/poll/"+id, http.StatusSeeOther)
}

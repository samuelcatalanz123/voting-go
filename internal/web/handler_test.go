package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func postForm(h *Handler, path string, form url.Values) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, req)
	return rec
}

// TestCrearYVotar crea una votación por la web y emite un voto, comprobando
// que el voto queda contado. No toca la red.
func TestCrearYVotar(t *testing.T) {
	h, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	rec := postForm(h, "/create", url.Values{
		"question": {"¿Color favorito?"},
		"options":  {"Rojo\nAzul"},
	})
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: código %d, esperaba 303", rec.Code)
	}
	loc := rec.Header().Get("Location") // /poll/{id}
	id := strings.TrimPrefix(loc, "/poll/")
	if id == "" || id == loc {
		t.Fatalf("create no redirigió a /poll/{id}, fue a %q", loc)
	}

	p, ok := h.store.Get(id)
	if !ok {
		t.Fatalf("la votación %q no se creó", id)
	}
	if len(p.Options) != 2 {
		t.Fatalf("esperaba 2 opciones, hay %d", len(p.Options))
	}

	// Votar por la opción 0 y comprobar que se sumó.
	postForm(h, "/vote/"+id, url.Values{"option": {"0"}})
	p, _ = h.store.Get(id)
	if p.Options[0].Votes != 1 {
		t.Errorf("esperaba 1 voto en la opción 0, hay %d", p.Options[0].Votes)
	}
}

// TestCreateRechazaIncompleta: una votación con menos de 2 opciones no se crea
// y redirige a la página de inicio.
func TestCreateRechazaIncompleta(t *testing.T) {
	h, _ := New()
	rec := postForm(h, "/create", url.Values{
		"question": {"¿Solo una?"},
		"options":  {"Única"},
	})
	if loc := rec.Header().Get("Location"); loc != "/" {
		t.Errorf("esperaba redirección a '/', fue a %q", loc)
	}
}

// TestPollDesconocido devuelve 404 para un id que no existe.
func TestPollDesconocido(t *testing.T) {
	h, _ := New()
	req := httptest.NewRequest(http.MethodGet, "/poll/noexiste", nil)
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("código %d, esperaba 404", rec.Code)
	}
}

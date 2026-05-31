// Package poll crea encuestas, registra votos y devuelve resultados.
package poll

import (
	"crypto/rand"
	"errors"
	"sync"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// ErrNotFound se devuelve cuando una encuesta no existe.
var ErrNotFound = errors.New("encuesta no encontrada")

// Option es una opción de la encuesta con sus votos.
type Option struct {
	Text  string
	Votes int
}

// Poll es una encuesta.
type Poll struct {
	ID       string
	Question string
	Options  []Option
}

// Total devuelve el número total de votos de la encuesta.
func (p Poll) Total() int {
	t := 0
	for _, o := range p.Options {
		t += o.Votes
	}
	return t
}

// Store guarda las encuestas en memoria.
type Store struct {
	mu    sync.Mutex
	polls map[string]*Poll
}

// New crea un Store vacío.
func New() *Store {
	return &Store{polls: make(map[string]*Poll)}
}

// Create crea una encuesta y devuelve su código.
func (s *Store) Create(question string, options []string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := newID()
	for _, existe := s.polls[id]; existe; _, existe = s.polls[id] {
		id = newID()
	}
	opts := make([]Option, len(options))
	for i, o := range options {
		opts[i] = Option{Text: o}
	}
	s.polls[id] = &Poll{ID: id, Question: question, Options: opts}
	return id
}

// Get devuelve una copia de la encuesta (segura para leer).
func (s *Store) Get(id string) (Poll, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.polls[id]
	if !ok {
		return Poll{}, false
	}
	cp := Poll{ID: p.ID, Question: p.Question, Options: append([]Option{}, p.Options...)}
	return cp, true
}

// Vote suma un voto a la opción indicada (por su índice).
func (s *Store) Vote(id string, optionIndex int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.polls[id]
	if !ok {
		return ErrNotFound
	}
	if optionIndex < 0 || optionIndex >= len(p.Options) {
		return errors.New("opción inválida")
	}
	p.Options[optionIndex].Votes++
	return nil
}

// newID genera un código aleatorio de 6 caracteres base62.
func newID() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b)
}

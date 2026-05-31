package poll

import "testing"

func TestCreateVoteCount(t *testing.T) {
	s := New()
	id := s.Create("¿Lenguaje favorito?", []string{"Go", "Python", "JS"})

	// Votamos: Go 2 veces, Python 1 vez
	_ = s.Vote(id, 0)
	_ = s.Vote(id, 0)
	_ = s.Vote(id, 1)

	p, ok := s.Get(id)
	if !ok {
		t.Fatal("la encuesta debería existir")
	}
	if p.Options[0].Votes != 2 {
		t.Errorf("Go = %d votos, esperaba 2", p.Options[0].Votes)
	}
	if p.Options[1].Votes != 1 {
		t.Errorf("Python = %d votos, esperaba 1", p.Options[1].Votes)
	}
	if p.Total() != 3 {
		t.Errorf("Total = %d, esperaba 3", p.Total())
	}
}

func TestVoteInvalid(t *testing.T) {
	s := New()
	id := s.Create("¿Sí o no?", []string{"Sí", "No"})
	if err := s.Vote(id, 5); err == nil {
		t.Error("una opción inválida debería dar error")
	}
	if err := s.Vote("noexiste", 0); err == nil {
		t.Error("una encuesta inexistente debería dar error")
	}
}

func TestGetUnknown(t *testing.T) {
	s := New()
	if _, ok := s.Get("noexiste"); ok {
		t.Error("una encuesta inexistente no debería encontrarse")
	}
}

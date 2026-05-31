// voting-go: una app web para crear votaciones y ver los resultados con barras.
package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/samuelcatalanz123/voting-go/internal/web"
)

func main() {
	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	h, err := web.New()
	if err != nil {
		log.Fatalf("no se pudo crear el handler: %v", err)
	}

	slog.Info("servidor iniciado", "abre", "http://localhost"+addr)
	if err := http.ListenAndServe(addr, h.Routes()); err != nil {
		log.Fatal(err)
	}
}

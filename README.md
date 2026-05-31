# App de votaciones (Go)

App web para crear encuestas (pregunta + opciones), votar y ver los resultados
con **barras de porcentaje**. Hecha en **Go**. Compartes un enlace y todos votan.

## Uso

```bash
go run .
```

Abre **http://localhost:8080**, escribe una pregunta y sus opciones (una por
línea), y pulsa **Crear votación**. Te da un enlace `/poll/{id}`; ábrelo, vota y
mira las barras crecer. Comparte el enlace para que otros voten.

## Cómo funciona

- `internal/poll`: guarda cada encuesta con un **código** corto. `Vote` suma un
  voto a la opción elegida; los resultados se calculan en porcentaje. Protegido
  con un `sync.Mutex` para que varios votos a la vez no choquen.
- `internal/web`: crear (`POST /create`), ver/votar (`GET /poll/{id}`,
  `POST /vote/{id}`). Las barras se dibujan con el porcentaje de cada opción.

## Estructura

```
main.go                 arranque del servidor
internal/poll/          encuestas: crear, votar, resultados + pruebas
internal/web/           crear y votar (handlers + plantillas)
```

## Pruebas

```bash
go test ./...
```

La prueba comprueba que los votos se cuentan en la opción correcta y que una
opción o encuesta inválida da error.

## Stack

Go (net/http, html/template, go:embed, crypto/rand, sync).

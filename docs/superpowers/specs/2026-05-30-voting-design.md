# Diseño: App de votaciones (Go)

**Fecha:** 2026-05-30 · **Estado:** Aprobado · **Autor:** Samuel (16º proyecto)

## Objetivo

Una app web para crear encuestas (pregunta + opciones), votar y ver los
resultados con **barras de porcentaje**. Objetivo de aprendizaje: contar votos y
calcular porcentajes; servir contenido por un código (como pastebin).

## Pantalla y rutas

- **Inicio** — `GET /`: formulario para crear una encuesta (pregunta + opciones,
  una por línea).
- **Crear** — `POST /create`: crea la encuesta y redirige a `/poll/{id}`.
- **Ver/Votar** — `GET /poll/{id}`: muestra la pregunta, botones para votar y los
  resultados con barras. Incluye el enlace para compartir. 404 si no existe.
- **Votar** — `POST /vote/{id}`: suma un voto a la opción elegida; redirige a `/poll/{id}`.

## Arquitectura

```
voting-go/
  main.go                 arranca el servidor (:8080 o PORT)
  internal/poll/
    poll.go               Option, Poll, Store: Create, Get, Vote; código base62
    poll_test.go          prueba: crear → votar → contar bien; opción inválida → error
  internal/web/
    handler.go            GET / (crear), POST /create, GET /poll/{id}, POST /vote/{id}
    templates/home.html   formulario para crear
    templates/poll.html   votar + resultados con barras
    static/style.css
  README.md
```

- **poll.go:** `Option{ Text string; Votes int }`. `Poll{ ID, Question string;
  Options []Option }` con método `Total()`. `Store{ mu sync.Mutex; polls map[string]*Poll }`,
  `New()`, `Create(question, options) string` (código único), `Get(id) (Poll, bool)`
  (devuelve una copia), `Vote(id, optionIndex) error`.
- **handler.go:** `create` valida (pregunta no vacía, ≥2 opciones); `poll` arma la
  vista con porcentajes (votos*100/total, cuidando división por cero) y el enlace;
  `vote` lee el índice de opción y llama a `Vote`.

## Pruebas

- **poll_test.go:** `Create` + `Vote` cuenta el voto en la opción correcta;
  `Vote` con índice inválido → error; `Get` de un id inexistente → false.
- `go build/vet/test` limpios.

## Seguridad

HTML escapado con `html/template`. Consultas/índices validados.

## Fuera de alcance (YAGNI)

Evitar votos repetidos (un voto por persona), cuentas, base de datos, cerrar encuesta.

## Criterios de éxito

1. `go run .` sirve en http://localhost:8080.
2. Crear una encuesta da un enlace `/poll/{id}`; votar suma y se ven las barras.
3. Los porcentajes son correctos.
4. La prueba de `poll` pasa.

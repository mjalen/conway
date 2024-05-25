package main

import (
	"flag"
	"log"
	"net/http"
	_ "embed"

	"conway-http/life"
	"conway-http/handler"
)

var (
	//go:embed frontend/index.html
	index string
	//go:embed frontend/game.html
	game string
)

func main() {
	sse := &handler.SSE{
		Rules: life.Rules{
			Birth:   []int{3},
			Survive: []int{2, 3},
		},
	}
	flag.IntVar(&(sse.Speed), "speed", 500, "Speed of the game.")
	flag.IntVar(&(sse.Size), "size", 32, "Size of the system.")
	flag.Parse()

	http.Handle("/game/connection", sse)
	http.Handle("/", &handler.HTML{ Content: index })
	http.Handle("/game/start", &handler.HTML{ Content: game })

	log.Fatal("HTTP server error: ", http.ListenAndServe(":8080", nil))
}

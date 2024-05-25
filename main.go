package main

import (
	"flag"
	"log"
	"fmt"
	"net/http"
	_ "embed"

	"github.com/a-h/templ"

	"github.com/mjalen/conway/life"
	"github.com/mjalen/conway/handler"
	"github.com/mjalen/conway/frontend"
)

func main() {
	sse := &handler.SSE{
		Rules: life.Rules{
			Birth:   []int{3},
			Survive: []int{2, 3},
		},
	}
	flag.IntVar(&(sse.Speed), "speed", 500, "Speed of the system.")
	flag.IntVar(&(sse.Size), "size", 32, "Size of the system.")
	flag.Int64Var(&(sse.Seed), "seed", 0, "Seed of the random system.")

	var port int
	flag.IntVar(&port, "port", 8080, "Port to serve to.")
	flag.Parse()

	http.Handle("/life/connection", sse)
	http.Handle("/", templ.Handler(frontend.Index()))
	http.Handle("/life/start", templ.Handler(frontend.Life()))

	portS := fmt.Sprintf(":%v", port)
	log.Printf("Serving to %v", portS)
	log.Fatal("HTTP server error: ", http.ListenAndServe(portS, nil))
}

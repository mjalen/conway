package main

import (
	"conway-http/life"
	"embed"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type SSEHandler struct {
	speed int
	size  int
	rules life.Rules
	Seed  int64
}

//go:embed index.html game.html
var filesystem embed.FS

var glider = []life.Pair{
	{X: 3, Y: 1},
	{X: 1, Y: 2},
	{X: 3, Y: 2},
	{X: 2, Y: 3},
	{X: 3, Y: 3},
}

var square = func(size int) []life.Pair {
	half := int(size / 2)
	return []life.Pair{
		{X: half - 1, Y: half - 1},
		{X: half - 1, Y: half},
		{X: half, Y: half - 1},
		{X: half, Y: half},
	}
}

func (h *SSEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	log.Printf("Client has connected.")
	render := make(chan *life.System)

	defer func() {
		if render != nil {
			close(render)
			render = nil
		}
		log.Printf("Client connection is lost.")
	}()

	go func() {
		for s := range render {
			fmt.Fprintf(w, "data: %s\n\n", s.ToHTML())
			w.(http.Flusher).Flush()
			time.Sleep(time.Duration(h.speed) * time.Millisecond)
		}
	}()

	s := life.System{
		Rules: h.rules,
		Size:  h.size,
	}
	h.Seed = rand.Int63n(9999999999999999)
	log.Printf("Seed: %v", h.Seed)
	for {
		if len(s.Alive) == 0 {
			s = s.Random(h.Seed)
		}
		render <- &s
		s = s.Next()
	}
}

type HTMLHandler struct {
	file string
}

func (h *HTMLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	content, err := filesystem.ReadFile(h.file)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
	}

	fmt.Fprintf(w, "%s", string(content))
}

func main() {
	sse := &SSEHandler{
		rules: life.Rules{
			Birth:   []int{3},
			Survive: []int{2, 3},
		},
	}
	flag.IntVar(&(sse.speed), "speed", 500, "Speed of the game.")
	flag.IntVar(&(sse.size), "size", 32, "Size of the system.")
	flag.Parse()

	http.Handle("/game/connection", sse)
	http.Handle("/", &HTMLHandler{"index.html"})
	http.Handle("/game/start", &HTMLHandler{"game.html"})

	log.Fatal("HTTP server error: ", http.ListenAndServe(":8080", nil))
}

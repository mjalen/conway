package handler

import (
	"net/http"
	"log"
	"time"
	"fmt"
	"context"

	"github.com/mjalen/conway/frontend"
	"github.com/mjalen/conway/life"
)

type SSE struct {
	Speed int
	Size  int
	Rules life.Rules
	Seed  int64
}

func (h *SSE) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
			fmt.Fprintf(w, "data: ")
			frontend.Board(s).Render(context.Background(), w)
			fmt.Fprintf(w, "\n\n")

			w.(http.Flusher).Flush()
			time.Sleep(time.Duration(h.Speed) * time.Millisecond)
		}
	}()

	s := life.System{
		Rules: h.Rules,
		Size:  h.Size,
		Seed: h.Seed,
	}
	for {
		if len(s.Alive) == 0 {
			s = s.Random()
			log.Printf("Seed: %v", s.Seed)
		}
		render <- &s
		s = s.Next()
	}
}

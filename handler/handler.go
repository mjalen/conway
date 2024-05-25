package handler

import (
	"net/http"
	"log"
	"time"
	"fmt"
	"math/rand"

	"conway-http/life"
)

type SSE struct {
	Speed int
	Size  int
	Rules life.Rules
	Seed  int64
}

type HTML struct {
	Content string
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
			fmt.Fprintf(w, "data: %s\n\n", s.ToHTML())
			w.(http.Flusher).Flush()
			time.Sleep(time.Duration(h.Speed) * time.Millisecond)
		}
	}()

	s := life.System{
		Rules: h.Rules,
		Size:  h.Size,
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

func (h *HTML) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	/*
	content, err := h.Filesystem.ReadFile(fmt.Sprintf("../%s", h.File))

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
	}
	*/

	if len(h.Content) == 0 {
		w.WriteHeader(404)
	}

	fmt.Fprintf(w, "%s", h.Content)
}

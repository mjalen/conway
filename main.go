package main

import (
	"fmt";
	"time";
	"flag";
	"net/http";
	"conway-http/conway";
	"embed";
	"log"
)

type Game struct {
	System conway.System2
	Speed int	
}

type SSEHandler struct {
	channel chan string
	game Game 
}

//go:embed index.html
var indexFile embed.FS 

func (h *SSEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	log.Printf("Client has connected.")
	h.channel = make(chan string) 

	defer func() {
		if h.channel != nil {
			close(h.channel)
			h.channel = nil
		}
		log.Printf("Client connection is lost.")
	}()

	go func() {
		for data := range h.channel {
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		}
	}()

	for {
		output, newSystem := conway.RunSystem2(h.game.System) 
		h.game.System = newSystem
		h.channel <- output 
		time.Sleep(time.Duration(h.game.Speed) * time.Millisecond)
	}
}

func main() {
	game := new(Game)
	game.System = *new(conway.System2)
	flag.IntVar(&(game.Speed), "speed", 500, "Speed of the game.")
	flag.IntVar(&(game.System.Size), "length", 32, "Length of the system.")
	flag.Parse()

	connection := new(SSEHandler)
	connection.game = *game

	http.Handle("/game", connection)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := indexFile.ReadFile("index.html")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, string(content))
	})

	log.Fatal("HTTP server error: ", http.ListenAndServe(":8080", nil))
}

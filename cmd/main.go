package main

import (
	"fmt"
	"net/http"

	"github.com/andersoncouto/sse/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() { // Thread 1
	out := make(chan amqp.Delivery)
	rabbitmqChannel, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	go rabbitmq.Consume("msgs", rabbitmqChannel, out) // Thread 2

	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		for m := range out {
			fmt.Fprintf(w, "event: message\n")
			fmt.Fprintf(w, "data: %s\n\n", m.Body)
			w.(http.Flusher).Flush()
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/index.html")
	})
	http.ListenAndServe(":8000", nil)
}

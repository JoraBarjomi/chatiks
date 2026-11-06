package main

import (
	"flag"
	"log"
	"net/http"

	"main/internal"
)

var addr = flag.String("addr", ":8090", "http service address")

func main() {
	flag.Parse()
	hub := internal.NewHub()

	go hub.Run()

	log.Println("Running on 8090...")
	http.HandleFunc("/", internal.ServeChat)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { internal.ServeWs(hub, w, r) })

	http.Handle("/ui/static/", http.StripPrefix("/ui/static/", http.FileServer(http.Dir("./ui/static/"))))

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

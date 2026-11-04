package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8090", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "home.html")
}

// как переключать между hub где хранить их в какой структуре данных нужно ли записывать данные в бд?
// добавлять удалять редачить как их хранить в бд же надо данные о них и подгружать данные о бд во время запуска но я могу в тестовом виде пока создать

func main() {
	flag.Parse()
	hub := newHub()

	go hub.run()

	log.Println("Running on 8090...")
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { serveWs(hub, w, r) })

	// http.Handle("/files", http.FileServer(http.Dir("files")))

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

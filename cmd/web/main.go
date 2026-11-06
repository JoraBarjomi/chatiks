package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"main/internal"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env: ", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	var addr = flag.String("addr", ":"+port, "http service address")
	flag.Parse()

	var store *internal.Store
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		if store, err = internal.InitStore(ctx, dsn); err != nil {
			log.Fatal("postgres: ", err)
		}
		log.Println("postgres: ok")
	} else {
		log.Println(".env: DATABASE_URL not found")
	}

	hub := internal.NewHub(store)

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

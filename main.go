package main

import (
	"log"
	"net/http"
)

func main() {
	// Serve static assets.
	http.Handle("/assets/", http.StripPrefix(
		"/assets/",
		http.FileServer(http.Dir("./web/assets")),
	))

	// Serve the frontend.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, "./web/index.html")
	})

	log.Println("CapIt is running at http://localhost:3000")

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

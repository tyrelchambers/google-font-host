package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	var PORT string = ":8080"

	if os.Getenv("PORT") != "" {
		PORT = ":" + os.Getenv("PORT")
	}

	if _, err := os.Stat("./fonts"); os.IsNotExist(err) {
		os.Mkdir("./fonts", 0755)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		fontName := query.Get("fontName")
		weight := query.Get("weight")
		variant := query.Get("variant")

		fontService := &FontServiceImpl{}

		fmt.Println("Looking for font ", fontName, weight, variant)

		if fontName == "" {
			http.Error(w, "Missing font name", 400)
			return
		}

		font, err := fontService.GetFont(fontName, weight, variant)

		if err != nil {
			http.Error(w, "Something went wrong", 500)
			return
		}

		w.Write(font)
	})

	handler := cors.Default().Handler(mux)

	fmt.Println("Listening on " + PORT)
	log.Fatal(http.ListenAndServe(PORT, handler))
}

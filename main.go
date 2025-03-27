package main

import (
	"embed"
	"flag"
	"html/template"
	"log"
	"net/http"
)

//go:embed web/templates/*
var htmlTemplates embed.FS

var port = flag.String("addr", ":1718", "http service address")

var t = template.Must(template.ParseFS(htmlTemplates, "web/templates/*"))

var tablero = NuevoTablero()

func main() {

	tablero[1] = X
	tablero[5] = O

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "index.html", tablero)
	})

	http.HandleFunc("/{indice}", func(w http.ResponseWriter, r *http.Request) {
		// Los valores de rune del 48 al 57 corresponden a los numeros 0 al 9
		indice := int([]rune(r.PathValue("indice"))[0]) - 48
		if indice < 0 || indice > 8 {
			t.ExecuteTemplate(w, "index.html", tablero)
			return
		}
		tablero[indice] = X
		t.ExecuteTemplate(w, "index.html", tablero)
	})

	log.Println("listening on", *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}

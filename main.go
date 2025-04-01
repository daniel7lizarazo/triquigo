package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
)

//go:embed web/templates/*
var htmlTemplates embed.FS

var port = flag.String("addr", ":1718", "http service address")

var t = template.Must(template.ParseFS(htmlTemplates, "web/templates/*"))

type EstadoJuego int

const (
	EnJuego = iota
	Ganado
	Perdido
	Empate
)

type EstadoTriqui struct {
	TableroActual Tablero
	Estado        EstadoJuego
	TrioGanador   Trio
}

var estadoTriqui = EstadoTriqui{
	TableroActual: NuevoTablero(),
	Estado:        EnJuego,
	TrioGanador:   Trio{},
}

func (tablero Tablero) ObtenerDisponibles() []int {
	disponibles := make([]int, 0, 9)
	k := 0
	for i := range tablero {
		if tablero[i] != Vacio {
			continue
		}
		disponibles = append(disponibles, i)
		k++
	}
	return disponibles
}

func (tablero *Tablero) JugarCelda() {
	disponibles := tablero.ObtenerDisponibles()
	rangoDisponible := len(disponibles)
	indiceCeldaDisponible := rand.Intn(rangoDisponible)
	indiceCeldaAJugar := disponibles[indiceCeldaDisponible]
	tablero[indiceCeldaAJugar] = O
}

// Los valores de rune del 48 al 57 corresponden a los numeros 0 al 9
// El string se convierte en un array de runes, y se toma el primer valore
// que es el unico que nos interesa, luego a ese valor se le resta 48
// que es el valor de "0"
func DigitoRune(s string) (int, error) {
	if len(s) <= 0 {
		return 0, fmt.Errorf("El string está vacio")
	}
	numero := []rune(s)[0] - 48
	if numero < 0 || numero > 9 {
		return 0, fmt.Errorf("El string %s tiene el valor %c por lo que no es un número", s)
	}

	return int(numero), nil
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if err := t.ExecuteTemplate(w, "index.html", estadoTriqui); err != nil {
			log.Print(err.Error())
		}

	})

	http.HandleFunc("/{indice}", func(w http.ResponseWriter, r *http.Request) {
		indice, err := DigitoRune(r.PathValue("indice"))

		if err != nil {
			return
		}

		estadoTriqui.TableroActual[indice] = X

		resultadoTrio, ganador := estadoTriqui.TableroActual.EstablecerGanador()

		if ganador == X {
			estadoTriqui.Estado = Ganado
			estadoTriqui.TrioGanador = resultadoTrio

			if err := t.ExecuteTemplate(w, "index.html", estadoTriqui); err != nil {
				log.Print(err.Error())
			}

			return
		}

		disponibles := estadoTriqui.TableroActual.ObtenerDisponibles()

		if len(disponibles) == 0 {
			estadoTriqui.Estado = Empate

			if err := t.ExecuteTemplate(w, "index.html", estadoTriqui); err != nil {
				log.Print(err.Error())
			}

			return
		}

		estadoTriqui.TableroActual.JugarCelda()
		resultadoTrio, ganador = estadoTriqui.TableroActual.EstablecerGanador()

		if ganador == O {
			estadoTriqui.Estado = Perdido
			estadoTriqui.TrioGanador = resultadoTrio

			if err := t.ExecuteTemplate(w, "index.html", estadoTriqui); err != nil {
				log.Print(err.Error())
			}

			return
		}

		if err := t.ExecuteTemplate(w, "index.html", estadoTriqui); err != nil {
			log.Print(err.Error())
		}

	})

	http.HandleFunc("/resetear", func(w http.ResponseWriter, r *http.Request) {

		estadoTriqui.TableroActual = NuevoTablero()
		estadoTriqui.Estado = EnJuego

		http.Redirect(w, r, "/", http.StatusFound)

	})

	log.Println("listening on", *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}

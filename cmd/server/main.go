package main

import (
	"daniel7lizarazo/triquigo/tablero"
	"embed"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
)

//go:embed web/templates/*
var htmlTemplates embed.FS

var t = template.Must(template.ParseFS(htmlTemplates, "web/templates/*"))

type EstadoJuego int

// Estado del juego
const (
	EnJuego = iota
	Ganado
	Perdido
	Empate
)

// Modos de juego
const (
	Tradicional  = "tradicional"
	Sincronizado = "sincronizado"
)

type EstadoTriqui struct {
	TableroActual  Tablero
	Estado         EstadoJuego
	ModoDeJuego    string
	TableroGanador Tablero
}

var estadoTriqui = EstadoTriqui{
	TableroActual:  NuevoTablero(),
	Estado:         EnJuego,
	ModoDeJuego:    Tradicional,
	TableroGanador: NuevoTablero(),
}

func (tablero *Tablero) ObtenerDisponibles() []int {
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

func (tablero *Tablero) ObtenerIndiceAleatorio() int {
	disponibles := tablero.ObtenerDisponibles()
	rangoDisponible := len(disponibles)
	indiceCeldaDisponible := rand.Intn(rangoDisponible)
	return disponibles[indiceCeldaDisponible]
}

func (tablero *Tablero) ObtenerOrdenado() int {
	for i := range tablero {
		if tablero[i] == Vacio {
			return i
		}
	}
	return 0
}

func (tablero *Tablero) EliminarBloqueada() {
	for i := range tablero {
		if tablero[i] == Bloqueada {
			tablero[i] = Vacio
			break
		}
	}
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

func tableroTradicionalHandler(w http.ResponseWriter, r *http.Request) {

	if err := t.ExecuteTemplate(w, "tableroTradicional.html", estadoTriqui); err != nil {
		log.Print(err.Error())
	}
}

func jugarTradicionalHandler(w http.ResponseWriter, r *http.Request) {

	indice, err := DigitoRune(r.PathValue("indice"))

	if err != nil {
		return
	}

	estadoTriqui.TableroActual[indice] = X

	_, ganador := estadoTriqui.TableroActual.EstablecerGanador()

	if ganador == X {
		estadoTriqui.Estado = Ganado

		if err := t.ExecuteTemplate(w, "tableroTradicional.html", estadoTriqui); err != nil {
			log.Print(err.Error())
		}

		return
	}

	disponibles := estadoTriqui.TableroActual.ObtenerDisponibles()

	if len(disponibles) == 0 {
		estadoTriqui.Estado = Empate

		if err := t.ExecuteTemplate(w, "tableroTradicional.html", estadoTriqui); err != nil {
			log.Print(err.Error())
		}

		return
	}

	celdaAJugar := estadoTriqui.TableroActual.ObtenerIndiceAleatorio()
	estadoTriqui.TableroActual[celdaAJugar] = O
	_, ganador = estadoTriqui.TableroActual.EstablecerGanador()

	if ganador == O {
		estadoTriqui.Estado = Perdido

		if err := t.ExecuteTemplate(w, "tableroTradicional.html", estadoTriqui); err != nil {
			log.Print(err.Error())
		}

		return
	}

	if err := t.ExecuteTemplate(w, "tableroTradicional.html", estadoTriqui); err != nil {
		log.Print(err.Error())
	}
}

func resetearTableroTradicionalHandler(w http.ResponseWriter, r *http.Request) {

	estadoTriqui.TableroActual = NuevoTablero()
	estadoTriqui.Estado = EnJuego
	estadoTriqui.TableroGanador = NuevoTablero()

	http.Redirect(w, r, "/tradicional", http.StatusFound)
}

func tableroSincronizadoHandler(w http.ResponseWriter, r *http.Request) {
	if err := t.ExecuteTemplate(w, "tableroSincronizado.html", estadoTriqui); err != nil {
		log.Print(err.Error())
	}
}

func jugarSincronizadoHandler(w http.ResponseWriter, r *http.Request) {

	indice, err := DigitoRune(r.PathValue("indice"))

	if err != nil {
		return
	}

	estadoTriqui.TableroGanador = NuevoTablero()

	// celdaAJugar := estadoTriqui.TableroActual.ObtenerIndiceAleatorio()
	celdaAJugar := estadoTriqui.TableroActual.ObtenerOrdenado()

	estadoTriqui.TableroActual.EliminarBloqueada()

	var disponibles []int

	if indice == celdaAJugar {

		estadoTriqui.TableroActual[indice] = Bloqueada

		disponibles = estadoTriqui.TableroActual.ObtenerDisponibles()

		if len(disponibles) == 0 {
			estadoTriqui.Estado = Empate

			if err := t.ExecuteTemplate(w, "tableroSincronizado.html", estadoTriqui); err != nil {
				log.Print(err.Error())
			}

			return
		}

		if err := t.ExecuteTemplate(w, "tableroSincronizado.html", estadoTriqui); err != nil {
			log.Print(err.Error())
		}

		return
	}

	estadoTriqui.TableroActual[indice] = X
	estadoTriqui.TableroActual[celdaAJugar] = O

	trioX, errX := estadoTriqui.TableroActual.EstablecerGanadorEsp(X)
	trioO, errO := estadoTriqui.TableroActual.EstablecerGanadorEsp(O)

	// Empate porque ganaron al tiempo
	if errX == nil && errO == nil {
		estadoTriqui.TableroGanador.AgregarTrioSignos(trioX, X)
		estadoTriqui.TableroGanador.AgregarTrioSignos(trioO, O)
		estadoTriqui.TableroActual.VaciarTrioTablero(trioX)
		estadoTriqui.TableroActual.VaciarTrioTablero(trioO)

		if err := t.ExecuteTemplate(w, "tableroSincronizado.html", estadoTriqui); err != nil {
			log.Print(err.Error())
		}

		return
	}

	// Si ambas no son nil entonces alguna debe ser diferente de nil, probaremos cual si es nil para declarar el ganador
	if errX == nil {
		estadoTriqui.Estado = Ganado
		estadoTriqui.TableroGanador.AgregarTrioSignos(trioX, X)

		if err := t.ExecuteTemplate(w, "tableroSincronizado.html", estadoTriqui); err != nil {
			log.Print(err.Error())
		}

		return
	}
	if errO == nil {
		estadoTriqui.Estado = Perdido
		estadoTriqui.TableroGanador.AgregarTrioSignos(trioO, O)

		if err := t.ExecuteTemplate(w, "tableroSincronizado.html", estadoTriqui); err != nil {
			log.Print(err.Error())
		}

		return
	}

	disponibles = estadoTriqui.TableroActual.ObtenerDisponibles()

	if len(disponibles) == 0 {
		estadoTriqui.Estado = Empate

		if err := t.ExecuteTemplate(w, "tableroSincronizado.html", estadoTriqui); err != nil {
			log.Print(err.Error())
		}

		return
	}

	if err := t.ExecuteTemplate(w, "tableroSincronizado.html", estadoTriqui); err != nil {
		log.Print(err.Error())
	}
}

func resetearTableroSincronizadoHandler(w http.ResponseWriter, r *http.Request) {

	estadoTriqui.TableroActual = NuevoTablero()
	estadoTriqui.Estado = EnJuego
	estadoTriqui.TableroGanador = NuevoTablero()

	http.Redirect(w, r, "/sincronizado", http.StatusFound)
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	if err := t.Execute(w, "menu.html"); err != nil {
		log.Print(err.Error())
	}
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	http.HandleFunc("/", menuHandler)

	http.HandleFunc("/tradicional", tableroTradicionalHandler)

	http.HandleFunc("/tradicional/{indice}", jugarTradicionalHandler)

	http.HandleFunc("/resetearTradicional", resetearTableroTradicionalHandler)

	http.HandleFunc("/sincronizado", tableroSincronizadoHandler)

	http.HandleFunc("/sincronizado/{indice}", jugarSincronizadoHandler)

	http.HandleFunc("/resetearSincronizado", resetearTableroSincronizadoHandler)

	port := ":1718"
	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

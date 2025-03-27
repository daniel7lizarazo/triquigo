package main

import (
	"fmt"
	"testing"
)

func TestTableroVacio(t *testing.T) {
	tablero := Tablero{}
	for i := range tablero {
		if tablero[i] != Vacio {
			t.Error("Hay valores diferentes a Vacio en el tablero vacio")
		}
	}
}

func TestTableroVacioEstablecerGanador(t *testing.T) {
	tablero := Tablero{}
	ganador := tablero.EstablecerGanador()
	if ganador != Vacio {
		t.Error(fmt.Sprintf("El ganador debería ser Vacio pero es %v", ganador))
	}
}

func TestTableroEstablecerGanadorXFila(t *testing.T) {
	tablero := Tablero{X, X, X}
	ganador := tablero.EstablecerGanador()
	fmt.Printf("El ganador es %v\n", ganador)
	if ganador != X {
		t.Error(fmt.Sprintf("El ganador debería ser X(1) pero es %v", ganador))
	}
}

func TestTableroEstablecerGanadorXFilaIncompleta(t *testing.T) {
	tablero := Tablero{X, Vacio, X}
	ganador := tablero.EstablecerGanador()
	fmt.Printf("El ganador es %v\n", ganador)
	if ganador != Vacio {
		t.Error(fmt.Sprintf("El ganador debería ser Vacio pero es %v", ganador))
	}
}

func TestTableroEstablecerGanadorOFila(t *testing.T) {
	tablero := Tablero{O, O, O}
	ganador := tablero.EstablecerGanador()
	fmt.Printf("El ganador es %v\n", ganador)
	if ganador != O {
		t.Error(fmt.Sprintf("El ganador debería ser O(2) pero es %v", ganador))
	}
}

func TestTableroEstablecerGanadorXColumna(t *testing.T) {
	tablero := Tablero{0: X, 3: X, 6: X}
	ganador := tablero.EstablecerGanador()
	fmt.Printf("El ganador es %v\n", ganador)
	if ganador != X {
		t.Error(fmt.Sprintf("El ganador debería ser X(1) pero es %v", ganador))
	}
}

func TestTableroEstablecerGanadorXColumnaIncompleta(t *testing.T) {
	tablero := Tablero{0: X, 3: Vacio, 6: X}
	ganador := tablero.EstablecerGanador()
	fmt.Printf("El ganador es %v\n", ganador)
	if ganador != Vacio {
		t.Error(fmt.Sprintf("El ganador debería ser Vacio pero es %v", ganador))
	}
}

func TestTableroEstablecerGanadorXDiagonal(t *testing.T) {
	tablero := Tablero{0: X, 4: X, 8: X}
	ganador := tablero.EstablecerGanador()
	fmt.Printf("El ganador es %v\n", ganador)
	if ganador != X {
		t.Error(fmt.Sprintf("El ganador debería ser X(1) pero es %v", ganador))
	}
}

func TestTableroEstablecerGanadorXDiagonalIncompleta(t *testing.T) {
	tablero := Tablero{0: X, 4: Vacio, 8: X}
	ganador := tablero.EstablecerGanador()
	fmt.Printf("El ganador es %v\n", ganador)
	if ganador != Vacio {
		t.Error(fmt.Sprintf("El ganador debería ser Vacio pero es %v", ganador))
	}
}

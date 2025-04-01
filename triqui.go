package main

const (
	Vacio = ' '
	X     = 'X'
	O     = 'O'
)

type Tablero [9]rune

type Trio [3]int

func (tablero *Tablero) compararTres(trio Trio) bool {
	return tablero[trio[0]] == tablero[trio[1]] && tablero[trio[1]] == tablero[trio[2]]
}

func NuevoTablero() Tablero {
	tablero := Tablero{}
	for i := 0; i < len(tablero); i++ {
		tablero[i] = Vacio
	}
	return tablero
}

func NuevoCeldasDisponibles() []int {
	celdasDisponibles := make([]int, 9)
	for i := 0; i < len(celdasDisponibles); i++ {
		celdasDisponibles[i] = i
	}
	return celdasDisponibles
}

func (trio *Trio) AsignarTrio(i, j, k int) {
	trio[0], trio[1], trio[2] = i, j, k
}

func (tablero *Tablero) EstablecerGanador() (Trio, rune) {
	// Verificaremos solo las primera fila y columna
	// En 0 y en 2 verificaremos las diagonales
	//
	// 0|1|2
	// 3|4|5
	// 6|7|8
	//
	// Si 0 es vacio no se verifican ni filas ni columnas pero si tiene algo verificamos
	// En 0 se puede ganar con una fila, una columna y una diagonal
	var trio Trio
	if tablero[0] != Vacio {
		// Verificamos la fila
		// 0|1|2      X|X|X
		// 3|4|5  ->  3|4|5
		// 6|7|8      6|7|8
		trio.AsignarTrio(0, 1, 2)
		if tablero.compararTres(trio) {
			return trio, tablero[trio[0]]
		}
		// Verificamos la columna
		// 0|1|2      X|1|2
		// 3|4|5  ->  X|4|5
		// 6|7|8      X|7|8
		trio.AsignarTrio(0, 3, 6)
		if tablero.compararTres(trio) {
			return trio, tablero[trio[0]]
		}
		// Verificamos la diagonal
		// 0|1|2      X|1|2
		// 3|4|5  ->  3|X|5
		// 6|7|8      6|7|X
		trio.AsignarTrio(0, 4, 8)
		if tablero.compararTres(trio) {
			return trio, tablero[trio[0]]
		}
	}
	// Seguimos por la misma fila comparando 1,
	// En esta celda la unica forma de ganar es completando una columna
	if tablero[1] != Vacio {
		// Verificamos la columna
		// 0|1|2      0|X|2
		// 3|4|5  ->  3|X|5
		// 6|7|8      6|X|8
		trio.AsignarTrio(1, 4, 7)
		if tablero.compararTres(trio) {
			return trio, tablero[trio[0]]
		}
	}
	// Seguimos por la misma fila comparando 2,
	// En esta celda la forma de ganar es completando una columna o la diagonal
	if tablero[2] != Vacio {
		// Verificamos la columna
		// 0|1|2      0|1|X
		// 3|4|5  ->  3|4|X
		// 6|7|8      6|7|X
		trio.AsignarTrio(2, 5, 8)
		if tablero.compararTres(trio) {
			return trio, tablero[trio[0]]
		}
		// Verificamos la diagonal
		// 0|1|2      0|1|X
		// 3|4|5  ->  3|X|5
		// 6|7|8      X|7|8
		trio.AsignarTrio(2, 4, 6)
		if tablero.compararTres(trio) {
			return trio, tablero[trio[0]]
		}
	}
	// Seguimos por la primera columna comparando 3,
	// En esta celda la unica forma de ganar es completando una fila
	if tablero[3] != Vacio {
		// Verificamos la columna
		// 0|1|2      0|1|2
		// 3|4|5  ->  X|X|X
		// 6|7|8      6|7|8
		trio.AsignarTrio(3, 4, 5)
		if tablero.compararTres(trio) {
			return trio, tablero[trio[0]]
		}
	}
	// Seguimos por la misma columna comparando 6,
	// En esta celda la forma de ganar es completando una fila,
	// ya que la diagonal ya fue verificada
	if tablero[6] != Vacio {
		// Verificamos la columna
		// 0|1|2      0|1|2
		// 3|4|5  ->  3|4|5
		// 6|7|8      X|X|X
		trio.AsignarTrio(6, 7, 8)
		if tablero.compararTres(trio) {
			return trio, tablero[trio[0]]
		}
	}
	return trio, Vacio
}

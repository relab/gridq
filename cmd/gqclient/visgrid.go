package main

import "fmt"

type visualGrid [][]rune

func newVisualGrid(rows, cols int) visualGrid {
	vgrid := make([][]rune, rows)
	cells := make([]rune, rows*cols)
	for i := range vgrid {
		vgrid[i], cells = cells[:cols], cells[cols:]
	}
	return vgrid
}

func (vs visualGrid) setRowQuorum(row uint32) {
	val := '-'
	for i := range vs {
		if i == int(row) {
			val = 'Q'
		}
		for j := range vs[i] {
			vs[i][j] = val
		}
		val = '-'
	}
}

func (vs visualGrid) setColQuorum(col uint32) {
	for i := range vs {
		for j := range vs[i] {
			if j == int(col) {
				vs[i][j] = 'Q'
			} else {
				vs[i][j] = '-'
			}
		}
	}
}

func (vs visualGrid) print() {
	fmt.Println("quorum:")
	for i := range vs {
		fmt.Printf("%c\n", vs[i])
	}
}

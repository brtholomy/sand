package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Coord struct {
	x uint32
	y uint32
}

type Pile struct {
	grid          [][]uint32
	height        uint32
	center_weight float32
}

type Record struct {
	seq      []Pile
	cascades map[int]int
}

func MakeRecord(size int, iters int) Record {
	seq := make([]Pile, iters, iters*size)
	casc := make(map[int]int, size)
	return Record{seq, casc}
}

func MakeGrid(size int) [][]uint32 {
	g := make([][]uint32, size)
	for i := range size {
		g[i] = make([]uint32, size)
	}
	return g
}

func MakePile(size int, height uint32, weight float32) Pile {
	g := MakeGrid(size)
	return Pile{g, height, weight}
}

func RandomCoord(p *Pile) Coord {
	size := len(p.grid) - 1
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Coord{uint32(r.Intn(size)), uint32(r.Intn(size))}
}

func PlaceGrain(p *Pile, c *Coord) {
	p.grid[c.x][c.y] += 1
}

func WillFall(p *Pile, c *Coord) bool {
	return p.grid[c.x][c.y] > p.height
}

func WithinGrid(p *Pile, c *Coord) bool {
	size := uint32(len(p.grid))
	return c.x < size && c.y < size
}

func GetNeighbors(c *Coord) [4]Coord {
	return [4]Coord{
		Coord{c.x + 1, c.y},
		Coord{c.x + 1, c.y + 1},
		Coord{c.x - 1, c.y},
		Coord{c.x - 1, c.y - 1},
	}
}

func Cascade(rec *Record, p *Pile, c *Coord, step int) {
	if WillFall(p, c) {
		p.grid[c.x][c.y] -= 4
		rec.cascades[step] += 1
		for _, v := range GetNeighbors(c) {
			if WithinGrid(p, &v) {
				PlaceGrain(p, &v)
				rec.seq = append(rec.seq, *p)
				Cascade(rec, p, &v, step)
			}
		}
	}
}

func Run(rec *Record, p *Pile, iters int) {
	for step := range iters {
		c := RandomCoord(p)
		PlaceGrain(p, &c)
		Cascade(rec, p, &c, step)
		rec.seq = append(rec.seq, *p)
	}
}

func GetTotals(rec *Record) map[int]int {
	totals := make(map[int]int, len(rec.cascades))
	for _, c := range rec.cascades {
		totals[c] += 1
	}
	return totals
}

func main() {
	size := 20
	iters := 1000
	var height uint32 = 4
	var weight float32 = 1.0
	rec := MakeRecord(size, iters)
	pile := MakePile(size, height, weight)

	Run(&rec, &pile, iters)
	totals := GetTotals(&rec)

	fmt.Println(pile)
	fmt.Println(totals)
}

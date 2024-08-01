package main

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/jessevdk/go-flags"
	// "golang.org/x/exp/maps"
	"math/rand"
	"os"
	"sort"
	"time"
)

var Opts struct {
	Size    int     `short:"s" long:"size" default:"10" description:"size of the grid"`
	Iters   int     `short:"i" long:"iters" default:"1000" description:"number of iterations"`
	Height  uint32  `short:"h" long:"height" default:"4" description:"maximum height of a pile"`
	Weight  float32 `short:"w" long:"weight" default:"1" description:"weight of center placement"`
	Verbose bool    `short:"v" long:"verbose" description:"verbose output"`
	Chart   bool    `short:"c" long:"chart" description:"make a chart of the totals"`
}

type Coord struct {
	x uint32
	y uint32
}

type Pile struct {
	grid          [][]uint32
	height        uint32
	center_weight float32
	r             rand.Rand
}

type Record struct {
	seq      []Pile
	cascades map[int]int
}

func MakeRecord(size int, iters int) Record {
	// NOTE: make enough room for the whole sequence:
	seq := make([]Pile, iters)
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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Pile{g, height, weight, *r}
}

func RandomCoord(p *Pile) Coord {
	size := len(p.grid) - 1
	return Coord{uint32(p.r.Intn(size)), uint32(p.r.Intn(size))}
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
		// TODO: only record a pile once per step for now:
		rec.seq[step] = *p
		// progress at some prime
		if 0 == step%255 {
			fmt.Println("step:", step)
		}
	}
}

func GetTotals(rec *Record) map[int]int {
	totals := make(map[int]int, len(rec.cascades))
	for _, c := range rec.cascades {
		totals[c] += 1
	}
	return totals
}

func MakeChart(totals map[int]int) {
	bar := charts.NewBar()

	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "sandpile simulation",
		Subtitle: "number of cascades per size of cascade",
	}))

	// TODO: this is stupid. should be a sorted map.
	// https://github.com/elliotchance/orderedmap
	values := make([]opts.BarData, 0)
	keys := make([]int, 0)
	for k, _ := range totals {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for k := range keys {
		values = append(values, opts.BarData{Value: totals[k]})
	}

	bar.SetXAxis(keys).
		AddSeries("number of cascades", values)
	f, _ := os.Create("bar.html")
	bar.Render(f)
}

func main() {
	_, err := flags.Parse(&Opts)
	if err != nil {
		panic(err)
	}
	rec := MakeRecord(Opts.Size, Opts.Iters)
	pile := MakePile(Opts.Size, Opts.Height, Opts.Weight)

	Run(&rec, &pile, Opts.Iters)
	totals := GetTotals(&rec)

	if Opts.Verbose {
		fmt.Println(pile)
	}
	if Opts.Chart {
		MakeChart(totals)
	}
	fmt.Println(totals)
}

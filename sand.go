package main

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	chartopts "github.com/go-echarts/go-echarts/v2/opts"
	"github.com/jessevdk/go-flags"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

type Opts struct {
	Size    int     `short:"s" long:"size" default:"10" description:"size of the grid"`
	Iters   int     `short:"i" long:"iters" default:"1000" description:"number of iterations"`
	Height  uint32  `short:"x" long:"max" default:"4" description:"maximum height of a pile"`
	Weight  float32 `short:"w" long:"weight" default:"1" description:"weight of center placement"`
	Verbose bool    `short:"v" long:"verbose" description:"verbose output"`
	Chart   bool    `short:"c" long:"chart" description:"make a chart of the totals"`
	Pixel   bool    `short:"p" long:"pixel" description:"make a pixel image of the pile"`
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
	return p.grid[c.x][c.y] >= p.height
}

func WithinGrid(p *Pile, c *Coord) bool {
	size := uint32(len(p.grid))
	return c.x < size && c.y < size
}

func GetNeighbors(c *Coord) [4]Coord {
	return [4]Coord{
		Coord{c.x + 1, c.y},
		Coord{c.x - 1, c.y},
		Coord{c.x, c.y + 1},
		Coord{c.x, c.y - 1},
	}
}

func Cascade(rec *Record, p *Pile, c *Coord, step int) {
	if WillFall(p, c) {
		p.grid[c.x][c.y] -= p.height
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
		// print progress every 100th
		if 0 == step%int(iters/100) {
			// NOTE: the \r carriage return
			fmt.Printf("\rprogress: %2d%%", int(100*step/iters))
		}
	}
	fmt.Println()
}

func GetTotals(rec *Record) map[int]float64 {
	totals := make(map[int]float64, len(rec.cascades))
	for _, c := range rec.cascades {
		totals[c] += 1
	}
	return totals
}

func GetLogTotals(totals *map[int]float64) map[int]float64 {
	logs := make(map[int]float64, len(*totals))
	for k, v := range *totals {
		logs[k] = math.Log(float64(v))
	}
	return logs
}

func MakeChart(totals map[int]float64, opts Opts) {
	bar := charts.NewBar()

	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(chartopts.Title{
		Title: "sandpile simulation",
		Subtitle: fmt.Sprintf(
			"log(cascades) per size of cascade, size:%d iters:%d height:%d",
			opts.Size, opts.Iters, opts.Height),
	}))

	// TODO: this is stupid. should be a sorted map.
	// https://github.com/elliotchance/orderedmap
	values := make([]chartopts.BarData, 0)
	keys := make([]int, 0)
	for k, _ := range totals {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for k := range keys {
		values = append(values, chartopts.BarData{Value: totals[k]})
	}

	bar.SetXAxis(keys).
		AddSeries("number of cascades", values)
	f, _ := os.Create("bar.html")
	bar.Render(f)
}

func MakeImage(p *Pile, opts Opts) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{opts.Size, opts.Size}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for i, x := range p.grid {
		for j, y := range x {
			factor := 255 / p.height
			c := color.RGBA{uint8(factor * y), 0, 0, 0xff}
			img.Set(i, j, c)
		}
	}
	f, _ := os.Create(fmt.Sprintf(
		"pile_%d-size_%d-iters_%d-height.png", opts.Size, opts.Iters, opts.Height))
	png.Encode(f, img)
}

func main() {
	opts := Opts{}
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}
	rec := MakeRecord(opts.Size, opts.Iters)
	pile := MakePile(opts.Size, opts.Height, opts.Weight)

	Run(&rec, &pile, opts.Iters)
	totals := GetTotals(&rec)
	logs := GetLogTotals(&totals)

	fmt.Println("log(cascade) at size 1:", logs[1])
	if opts.Verbose {
		fmt.Println("last pile:", pile)
		fmt.Println("totals:", totals)
		fmt.Println("logs:", logs)
	}
	if opts.Chart {
		MakeChart(logs, opts)
	}
	if opts.Pixel {
		MakeImage(&pile, opts)
	}
}

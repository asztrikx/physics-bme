package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

//plotWrite creates predefined plot style based on XYs values
func plotWrite(xys plotter.XYs, filename string) {
	//plot
	p, err := plot.New()
	if err != nil {
		log.Fatalln(err)
	}

	//texts
	p.X.Label.Text = "Fok"
	p.Y.Label.Text = "Eltérés %"
	p.Title.Text = "Eltérés % kis kitérítésű fonálingának periódusától az egész fokokra"
	size := 2 * vg.Centimeter
	p.X.Label.Font.Size = size
	p.Y.Label.Font.Size = size
	p.Y.Tick.Label.Font.Size = size
	p.X.Tick.Label.Font.Size = 1.5 * vg.Centimeter
	p.Title.Font.Size = size

	//axis tick
	p.X.Tick.Marker = tick{}
	p.Y.Tick.Marker = tick{}

	//grid
	g := plotter.NewGrid()
	g.Horizontal.Width = 1
	p.Add(g)

	//line
	line, err := plotter.NewLine(xys)
	if err != nil {
		panic(err)
	}
	line.Color = color.RGBA{R: 255, A: 255}
	line.Width = vg.Points(4)
	p.Add(line)

	//save
	if err := p.Save(5000, 2000, filename); err != nil {
		panic(err)
	}
}

type tick struct{}

func (tick) Ticks(min, max float64) []plot.Tick {
	tickS := make([]plot.Tick, int(math.Ceil(max))-int(math.Floor(min))+1)
	for i := 0; i < len(tickS); i++ {
		tickS[i].Value = float64(int(math.Floor(min)) + i)
		tickS[i].Label = fmt.Sprintf("%2.0f", tickS[i].Value)
	}
	return tickS
}

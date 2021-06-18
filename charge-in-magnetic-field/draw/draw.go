package draw

import (
	"charge-in-magnetic-field/physics"
	"charge-in-magnetic-field/vector"
	"fmt"
	"log"

	"github.com/fogleman/gg"
)

var config physics.Config
var ctx *gg.Context
var buffer []vector.Vector
var bufferIndex int

func SetRGB(r, g, b float64) {
	//unflushed items would also get this colour so flush
	if bufferIndex != 0 {
		flush()
	}

	ctx.SetRGB(r, g, b)
}
func Draw(v vector.Vector) {
	buffer[bufferIndex] = v
	bufferIndex++

	//buffer may get full
	if bufferIndex == len(buffer) {
		flush()
	}
}
func flush() {
	for i := 0; i < bufferIndex; i++ {
		ctx.DrawPoint(buffer[i].X, float64(config.MF.Height)-buffer[i].Y, 1)
	}
	ctx.Fill()
	bufferIndex = 0
}
func Start(c physics.Config) {
	config = c

	buffer = make([]vector.Vector, config.RAM)
	ctx = gg.NewContext(int(config.MF.Width), int(config.MF.Height))
}
func Save() {
	//buffer may have unflushed points
	if bufferIndex != 0 {
		flush()
	}

	if err := ctx.SavePNG(fmt.Sprintf("result-%s.png", config.TimeDelta)); err != nil {
		log.Fatalln(err)
	}
}

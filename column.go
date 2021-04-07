package main

import (
	"math/rand"
)

type Columns struct {
	Gate []gate
	ColumnWidth
	GateHeight
	ColumnCount float64
	Speed float64
	MaxColumnCount float64
	lastID float64
	window size
}

//type window struct {
//	width int
//	height int
//}

type gate struct {
	X float64
	Y float64
	Width float64
	Height float64
	ID float64
}

type ColumnWidth struct {
	min float64
	max float64
}

type GateHeight struct {
	Start float64
	Finish float64
}

func (c Columns) Frame () Frame {

	data := make(Frame, 0, int(c.MaxColumnCount * 2))

	for _, gate := range c.Gate {
		data = append(data, obj{
			gate.X,
			0,
			gate.Width,
			gate.Y,
			 2.0,//"column",
		})

		data = append(data, obj{
			gate.X,
			gate.Y + gate.Height,
			gate.Width,
			c.window.height,
			2.0,//"column",
		})
	}

	return data
}

func (c *Columns) NewColumn () {
	width := RandMinMax(c.ColumnWidth.min, c.ColumnWidth.max)
	GateSize := RandMinMax(c.GateHeight.Start, c.GateHeight.Finish)
	startPointByY := RandMinMax(100, c.window.height - 100)
	id := c.GetID()

	c.Gate = append(c.Gate, gate{
		c.window.width,
		startPointByY,
		width,
		GateSize,
		id,
	})

	c.ColumnCount++
}

func (c *Columns) Process () {
	removeStatus := false
	if len(c.Gate) == 0 {
		c.NewColumn()
	} else {
		for i := range c.Gate {
			// gate move
			c.Gate[i].X = c.Gate[i].X - c.Speed

			// gate remove
			if c.Gate[i].X + c.Gate[i].Width < 0 {
				removeStatus = true
			}

			// gate add
			if i == len(c.Gate) - 1 {
				if c.Gate[i].X + c.Gate[i].Width < c.window.width - (c.window.width / c.MaxColumnCount) {
					c.NewColumn()
				}
			}
		}
	}

	if removeStatus {
		c.RemoveFirstColumn()
	}
}

func (c *Columns) RemoveFirstColumn () {
	if len(c.Gate) > 0 {
		if len(c.Gate) == 1 {
			c.Gate = make([]gate, 0, int(c.MaxColumnCount))
		} else {
			c.Gate = c.Gate[1:]
		}
		c.ColumnCount--
	}
}

func (c *Columns) GetID () float64 {
	c.lastID++
	return c.lastID
}

func (c Columns) GetGate () []gate {
	return c.Gate
}

func RandMinMax (min float64, max float64) float64 {
	// rand.Seed(time.Now().UnixNano())
	return float64(rand.Intn(int(max) - int(min)) + int(min))
}
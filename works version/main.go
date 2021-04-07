package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"syscall/js"
	"time"
)

func main() {

	//defer Render.Release()

	//ch := make(chan interface{}, 500)
	frame := make(Frame, 0, 20)

	bird := Bird{}

	Col := Columns{}

	tick := time.Tick((1000 / 60) * time.Millisecond)

	go Render(&frame)

	// works
	js.Global().Set("GetRays", GetRays(&bird))
	js.Global().Set("GetScore", GetScore(&bird))
	js.Global().Set("GetStatus", GetStatus(&bird))

	for {


		// ask restart/start status
		document := js.Global().Call("Restart")
		test := fmt.Sprintf("%v", document)
		if test != "<boolean: true>" {
			time.Sleep(1 * time.Second)
			continue
		}

		restart(&bird, &Col)

		for range tick {
			Col.Process()

			// build frames
			frame = make(Frame, 0, 40)

			// ask jump status
			document := js.Global().Call("isStatus")
			test := fmt.Sprintf("%v", document)
			//fmt.Println(test)
			//fmt.Println(test == "<boolean: true>")
			if test == "<boolean: true>" && bird.GetStatus() {
				bird.Jump()
			} else {
				bird.Gravity(Col.window.height)
			}

			// coliseum
			bd := bird.StampOfData()

			// rays
			rays := make([]Rays, 0, 7)

			Gates := Col.GetGate()

			rays = append(rays,
				raysData(
					bird.Position.x,
					bird.Position.y,
					-90,
					ray(
						bird.Position.x,
						bird.Position.y,
						Gates,
						Col.window,
						-90,
						1,
						0)),
				raysData(bird.Position.x+bird.Size.width,
					bird.Position.y, -90,
					ray(bird.Position.x+bird.Size.width,
						bird.Position.y,
						Gates,
						Col.window,
						-90,
						1,
						0)),
				raysData(bird.Position.x+bird.Size.width,
					bird.Position.y, -45,
					ray(bird.Position.x+bird.Size.width,
						bird.Position.y,
						Gates,
						Col.window,
						-45,
						1,
						0)),

						// тут новый луч
				raysData(bird.Position.x+bird.Size.width,
					bird.Position.y + (bird.Size.height / 2), -5,
					ray(bird.Position.x+bird.Size.width,
						bird.Position.y + (bird.Size.height / 2),
						Gates,
						Col.window,
						-5,
						1,
						0)),
				raysData(bird.Position.x+bird.Size.width,
					bird.Position.y + (bird.Size.height / 2), 5,
					ray(bird.Position.x+bird.Size.width,
						bird.Position.y + (bird.Size.height / 2),
						Gates,
						Col.window,
						5,
						1,
						0)),

				raysData(bird.Position.x+bird.Size.width,
					bird.Position.y+bird.Size.height, 45, ray(bird.Position.x+bird.Size.width,
						bird.Position.y+bird.Size.height,
						Gates,
						Col.window,
						45,
						1,
						0)),
				raysData(bird.Position.x+bird.Size.width,
					bird.Position.y+bird.Size.height, 90, ray(bird.Position.x+bird.Size.width,
						bird.Position.y+bird.Size.height,
						Gates,
						Col.window,
						90,
						1,
						0)),
				raysData(bird.Position.x,
					bird.Position.y+bird.Size.height, 90, ray(bird.Position.x,
						bird.Position.y+bird.Size.height,
						Gates,
						Col.window,
						90,
						1,
						0)),
				// первый луч вверх),
			)

			for _, gate := range Gates { // Col.GetGate()

				if bd.X+bd.Size.width >= gate.X && bd.X <= gate.X+gate.Width {
					if bd.Y <= gate.Y || bd.Y+bd.Size.height >= gate.Y+gate.Height {
						bird.IsAlive = false
					}
				}

				if bd.X >= gate.X+gate.Width && bd.X <= gate.X+gate.Width+Col.Speed {
					bird.Score++
				}
			}

			bird.Rays = rays
			//fmt.Println(bird.Rays)

			frame = append(frame, Col.Frame()...)
			frame = append(frame, bird.Frame(), bird.GetScore())
			frame = append(frame, bird.FrameForRays()...)

			if !bird.GetStatus() {
				break
			}
			//ch <- frame
			//fmt.Println("works")
		}
	}

}

// ray -
// x - current Point position by x
// y - current Point position by y
// g - gate struct with position and size
// angle - float64
// speed - float64
// distances - float64

//func makeCh () chan []float64 {
//	return make(chan []float64, 50)
//}

func raysData(x, y float64, angle float64, speed float64) Rays {
	NewX := 0.0
	NewY := 0.0
	u := 1.0
	if angle != 0 {
		u = angle * (math.Pi / 180)
		NewX = x + math.Cos(u)*speed // новая точка
		NewY = y + math.Sin(u)*speed
	} else {
		NewX = x + 1.0 *speed // новая точка
		NewY = y // новая точка
	}
	  // что-то ? :)


	length := math.Abs(math.Sqrt(math.Pow(NewX-x, 2) + math.Pow(NewY-y, 2)))

	// дистация от точки до точки
	return Rays{
		length,
		x,
		y,
		NewX,
		NewY,
	}
}

func ray(x, y float64, gates []gate, wSize size, angle float64, speed float64, distances float64) float64 {

	NewX := 0.0
	NewY := 0.0
	u := 1.0
	if angle != 0 {
		u = angle * (math.Pi / 180)
		NewX = x + math.Cos(u)*speed // новая точка
		NewY = y + math.Sin(u)*speed
	} else {
		NewX = x + (1.0 * speed) // новая точка
		NewY = y+1 // новая точка
	}

	if NewY >= wSize.height || NewY <= 0  || NewX >= wSize.width {
		return math.Abs(math.Sqrt(math.Pow(NewX-x, 2) + math.Pow(NewY-y, 2)))
	} else {
		for _, gt := range gates {
			if (NewX >= gt.X && NewX <= gt.X+gt.Width) && NewY <= gt.Y {
				return math.Abs(math.Sqrt(math.Pow(NewX-x, 2) + math.Pow(NewY-y, 2)))
			}

			if (NewX >= gt.X && NewX <= gt.X+gt.Width) && NewY >= gt.Y+gt.Height {
				return math.Abs(math.Sqrt(math.Pow(NewX-x, 2) + math.Pow(NewY-y, 2)))
			}
		}
		//if (NewX >= g.X && NewX <= g.X + g.Width) && NewY <= g.Y {
		//	return math.Abs(math.Sqrt(math.Pow(NewX - x, 2)  + math.Pow(NewY - y, 2)))
		//}
		//
		//if (NewX >= g.X && NewX <= g.X + g.Width) && NewY >= g.Y + g.Height {
		//	return math.Abs(math.Sqrt(math.Pow(NewX - x, 2)  + math.Pow(NewY - y, 2)))
		//}
	}
	//if (NewX >= g.X && NewY <= g.Y) ||
	//	(NewX >= g.X && NewY >= g.Y + g.Height) ||
	//	(NewY >= wSize.height || NewY <= 0) {
	//	return math.Abs(math.Sqrt(math.Pow(NewX - x, 2)  + math.Pow(NewY - y, 2)))
	//}
	// дистация от точки до точки
	return distances + ray(NewX, NewY, gates, wSize, angle, speed, math.Abs(math.Sqrt(math.Pow(NewX-x, 2)+math.Pow(NewY-y, 2))))
	//distances = distances + math.Abs(math.Sqrt(math.Pow(NewX - x, 2)  + math.Pow(NewY - y, 2)))
}

func GetRays(data *Bird) js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		v, err := json.Marshal((*data).getRaysLengths())
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println((*data).Rays)
		dst := js.Global().Get("Uint8Array").New(len(v))
		js.CopyBytesToJS(dst, v)
		// fmt.Println(dst)
		return dst
	})
}

func GetScore(data *Bird) js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		v, err := json.Marshal((*data).Score)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println((*data).Score)
		dst := js.Global().Get("Uint8Array").New(len(v))
		js.CopyBytesToJS(dst, v)
		// fmt.Println(dst)
		return dst
	})
}

func GetStatus(data *Bird) js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		v, err := json.Marshal((*data).GetStatus())
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println((*data).Score)
		dst := js.Global().Get("Uint8Array").New(len(v))
		js.CopyBytesToJS(dst, v)
		// fmt.Println(dst)
		return dst
	})
}

func restart(bird *Bird, Col *Columns) {
	*bird = Bird{
		Position: position{
			640,
			360,
		},
		IsAlive: true,
		Score:   0,
		Size: size{
			40,
			25,
		},
		ID:             0,
		Speed:          7,
		ForceOfGravity: 7,
	}

	*Col = Columns{
		Gate: make([]gate, 0, 5),
		ColumnWidth: ColumnWidth{
			150,
			250,
		},
		GateHeight: GateHeight{
			100,
			250,
		},
		ColumnCount:    0,
		Speed:          2,
		MaxColumnCount: 5,
		lastID:         0,
		window: size{
			1280,
			720,
		},
	}
}

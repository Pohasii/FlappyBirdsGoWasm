package main

import (
	"encoding/json"
	"log"
	"math"
	"syscall/js"
	"time"
)

func main() {

	//defer Render.Release()

	//ch := make(chan interface{}, 500)
	frame := make(Frame, 0, 20)

	countOfBs := js.Global().Call("getCountOfBirds").Int()
	// fmt.Println(countOfBs)

	birds := make([]Bird, countOfBs, 20)

	Col := Columns{}

	tick := time.Tick((1000 / 60) * time.Millisecond)

	go Render(&frame)

	// works
	js.Global().Set("GetRays", GetRays(&birds))
	js.Global().Set("GetScore", GetScore(&birds))
	js.Global().Set("GetStatus", GetStatus(&birds))

	for {



		// ask restart/start status
		ReStatus := js.Global().Call("Restart").Bool()
		// test := fmt.Sprintf("%v", document)
		if ReStatus != true { // if test != "<boolean: true>" {
			time.Sleep(1 * time.Second)
			continue
		}

		//countOfBs := js.Global().Call("getCountOfBirds").Int()
		////data := js.Global().Call("getCountOfBirds")
		////countOfBs := 0
		////err := json.Unmarshal([]byte(data),&countOfBs)
		////if err != nil {
		////	log.Println(err)
		////}
		//fmt.Println(countOfBs)

		restart(&birds, countOfBs, &Col)

		for range tick {
			Col.Process()

			// build frames
			frame = make(Frame, 0, 30*countOfBs)

			// ask jump status
			document := js.Global().Call("getBirdsJumpStatus").String()
			statuses := make([]bool, 0, 30)
			err := json.Unmarshal([]byte(document), &statuses)
			if err != nil {
				log.Println(err)
			}

			// fmt.Println(statuses)

			Gates := Col.GetGate()

			for i, _ := range birds {
				if statuses[i] == true && birds[i].GetStatus() {
					birds[i].Jump()
				} else {
					birds[i].Gravity(Col.window.height)
				}

				// coliseum
				bd := birds[i].StampOfData()

				// rays
				rays := make([]Rays, 0, 7)

				rays = append(rays,
					raysData(
						birds[i].Position.x,
						birds[i].Position.y,
						-90,
						ray(
							birds[i].Position.x,
							birds[i].Position.y,
							Gates,
							Col.window,
							-90,
							1,
							0)),
					raysData(birds[i].Position.x+birds[i].Size.width,
						birds[i].Position.y, -90,
						ray(birds[i].Position.x+birds[i].Size.width,
							birds[i].Position.y,
							Gates,
							Col.window,
							-90,
							1,
							0)),
					raysData(birds[i].Position.x+birds[i].Size.width,
						birds[i].Position.y, -45,
						ray(birds[i].Position.x+birds[i].Size.width,
							birds[i].Position.y,
							Gates,
							Col.window,
							-45,
							1,
							0)),

					// тут новый луч
					raysData(birds[i].Position.x+birds[i].Size.width,
						birds[i].Position.y + (birds[i].Size.height / 2), -5,
						ray(birds[i].Position.x+birds[i].Size.width,
							birds[i].Position.y + (birds[i].Size.height / 2),
							Gates,
							Col.window,
							-5,
							1,
							0)),
					raysData(birds[i].Position.x+birds[i].Size.width,
						birds[i].Position.y + (birds[i].Size.height / 2), 5,
						ray(birds[i].Position.x+birds[i].Size.width,
							birds[i].Position.y + (birds[i].Size.height / 2),
							Gates,
							Col.window,
							5,
							1,
							0)),

					raysData(birds[i].Position.x+birds[i].Size.width,
						birds[i].Position.y+birds[i].Size.height, 45, ray(birds[i].Position.x+birds[i].Size.width,
							birds[i].Position.y+birds[i].Size.height,
							Gates,
							Col.window,
							45,
							1,
							0)),
					raysData(birds[i].Position.x+birds[i].Size.width,
						birds[i].Position.y+birds[i].Size.height, 90, ray(birds[i].Position.x+birds[i].Size.width,
							birds[i].Position.y+birds[i].Size.height,
							Gates,
							Col.window,
							90,
							1,
							0)),
					raysData(birds[i].Position.x,
						birds[i].Position.y+birds[i].Size.height, 90, ray(birds[i].Position.x,
							birds[i].Position.y+birds[i].Size.height,
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
							birds[i].IsAlive = false
						}
					}

					if bd.X >= gate.X+gate.Width && bd.X <= gate.X+gate.Width+Col.Speed {
						birds[i].Score++
					}
				}

				birds[i].Rays = rays
				//fmt.Println(bird.Rays)

				if i == 0 {
					frame = append(frame, birds[i].Frame(), birds[i].GetScore())
					frame = append(frame, birds[i].FrameForRays()...)
				} else {
					if birds[i].GetStatus() {
						frame = append(frame, birds[i].Frame())
					}
				}
			}

			frame = append(frame, Col.Frame()...)

			stats := 0
			for i := range birds {
				if birds[i].GetStatus() {
					stats++
				}
			}

			if stats == 0 {
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
	}

	// дистация от точки до точки
	return distances + ray(NewX, NewY, gates, wSize, angle, speed, math.Abs(math.Sqrt(math.Pow(NewX-x, 2)+math.Pow(NewY-y, 2))))
	//distances = distances + math.Abs(math.Sqrt(math.Pow(NewX - x, 2)  + math.Pow(NewY - y, 2)))
}

func GetRays(data *[]Bird) js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		id := args[0].Int()
		v, err := json.Marshal((*data)[id].getRaysLengths())
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

func GetScore(data *[]Bird) js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		id := args[0].Int()
		v, err := json.Marshal((*data)[id].Score)
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

func GetStatus(data *[]Bird) js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		id := args[0].Int()
		v, err := json.Marshal((*data)[id].GetStatus())
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

func restart(bird *[]Bird, count int, Col *Columns) {

	for i := 0; i < count; i++ { // i := range *bird
		(*bird)[i] = Bird{
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

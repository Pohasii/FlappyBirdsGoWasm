package main

import (
	"encoding/json"
	"fmt"
	"log"
	"syscall/js"
)

type Rays struct {
	length float64
	StartX float64
	StartY float64
	FinishX float64
	FinishY float64
}

type size struct {
	width  float64
	height float64
}

type obj struct {
	X float64
	Y float64
	Width float64
	Height float64
	Type float64
}

type Frame []obj

func (d Frame) Get(ch chan interface{}) js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		v, err := json.Marshal(<-ch)
		if err != nil {
			log.Fatal(err)
		}

		dst := js.Global().Get("Uint8Array").New(len(v))
		js.CopyBytesToJS(dst, v)

		return dst
	})
}

//func (d Frame) Get(ch chan interface{}) js.Func {
//
//	//data := make(chan []interface{})
//	//
//	//fmt.Println(f)
//	//value := make([]interface{}, 0,21)
//	//for _, v := range f {
//	//	fmt.Println(v)
//	//	value = append(value, v)
//	//}
//	//fmt.Println(value)
//	//data <- value
//	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
//		// Return a JS dictionary with two keys (of heterogeneous type)
//		//return map[string]interface{}{
//		//	"hello":  "world12",
//		//	"answer": 43,
//		//}
//		//var d interface{}
//		//select {
//		//case da := <-ch:
//		//	d = da
//		//default:
//		//	d = "[]"
//		//}
//
//
//		v, err := json.Marshal(<-ch)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		dst := js.Global().Get("Uint8Array").New(len(v))
//		js.CopyBytesToJS(dst, v)
//
//		//value := make([]interface{}, 0,21)
//		//for _, v := range d {
//		//	// fmt.Println(v)
//		//	value = append(value, v)
//		//}
//		//asd["asd"] = value
//		return dst
//	})
//}

func Render (f *Frame) {
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "canvas")
	ctx := canvasEl.Call("getContext", "2d")

	var renderFrame js.Func
	var tmark float64
	var markCount = 0
	var tdiffSum float64


	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		now := args[0].Float()
		tdiffSum += now - tmark
		markCount++
		if markCount > 10 {

			ctx.Set("fillStyle", "black")
			ctx.Call("fillText", fmt.Sprintf("FPS: %.01f", 1000/(tdiffSum/float64(markCount))), 10, 10)
			//doc.Call("getElementById", "fps").Set("innerHTML", fmt.Sprintf("FPS: %.01f", 1000/(tdiffSum/float64(markCount))))
			tdiffSum, markCount = 0, 0
		}
		tmark = now

		ctx.Call("clearRect", 0, 0, 1280, 720)

		for _, obj := range *f {
			switch obj.Type {
			case 1.0: //bird
				ctx.Set("fillStyle", "#b22222")
				ctx.Call("fillRect", obj.X, obj.Y, obj.Width, obj.Height)
			case 2.0: //column
				ctx.Set("fillStyle", "green")
				ctx.Call("fillRect", obj.X, obj.Y, obj.Width, obj.Height)
			case 3.0: //Score
				ctx.Set("fillStyle", "black")
				ctx.Call("fillText", obj.Width, obj.X, obj.Y)
			case 4.0: //Score
				ctx.Set("fillStyle", "red")
				// ctx.Call("fillRect", obj.X, obj.Y, obj.Width, obj.Height)
				ctx.Call("beginPath")
				ctx.Call("moveTo", obj.X, obj.Y)
				ctx.Call("lineTo", obj.Width, obj.Height)
				ctx.Call("stroke")
			default:
				// freebsd, openbsd,
				// plan9, windows...
				fmt.Println("err draw")
			}
		}

		//ctx.Set("fillStyle", "yellow")
		//ctx.Call("fillRect", 50, 50, 50, 50)

		// can.fillStyle = 'yellow'
		// can.fillRect(obj.X, obj.Y, obj.Width, obj.Height)
		//ctx.Set("globalAlpha", 0.5)
		//ctx.Call("beginPath")
		//ctx.Set("fillStyle", fmt.Sprintf("#%06x", dot.color))
		//ctx.Set("strokeStyle", fmt.Sprintf("#%06x", dot.color))
		//ctx.Set("lineWidth", dot.size)
		//ctx.Call("arc", dot.pos[0], dot.pos[1], dot.size, 0, 2*math.Pi)
		//ctx.Call("fill")

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	//defer renderFrame.Release()

	// Start running
	js.Global().Call("requestAnimationFrame", renderFrame)

	}
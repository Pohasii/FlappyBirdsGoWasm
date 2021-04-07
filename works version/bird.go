package main

type Bird struct {
	Position position
	IsAlive  bool
	Score    float64
	Size     size
	ID       float64
	Speed float64
	ForceOfGravity float64
	Rays []Rays
}

type position struct {
	x float64
	y float64
}

//type size struct {
//	width  int
//	height int
//}

func (b Bird) getRaysLengths () []float64{
	rays := make([]float64, 0, 10)
	for _, Ray := range b.Rays {
		rays = append(rays, Ray.length)
	}
	return rays
}

func (b *Bird) Gravity (wSize float64) {
	b.Position.y = b.Position.y + b.ForceOfGravity
	if b.Position.y + b.Size.height >= wSize {
		b.IsAlive = false
		b.Position.y = wSize
	}
}

func (b *Bird) Jump () {
	if b.GetStatus() {
		b.Position.y = b.Position.y - b.Speed
		if b.Position.y <= 0 {
			b.IsAlive = false
			b.Position.y = 0
		}
	}
}

func (b Bird) GetScore() obj {
	return obj {
		b.Position.x + b.Size.width /2,
		b.Position.y + b.Size.height /2,
		b.Score,
		0,
		3.0,
	}
}

func (b Bird) Frame() obj {
	return obj {
		b.Position.x , // - b.Size.width
		b.Position.y , // - b.Size.height
		b.Size.width,
		b.Size.height,
		1.0,
	}
}

func (b Bird) FrameForRays() []obj {

	r := make([]obj, 0, 12)

	for _, ra := range b.Rays {
		r = append(r, obj {
			ra.StartX,
			ra.StartY,
			ra.FinishX,
			ra.FinishY,
			4.0,
		})
	}
	return r
}

func (b Bird) GetStatus() bool {
	return b.IsAlive
}

//func (b Bird) GetScore() int {
//	return b.Score
//}

func (b Bird) GetID() float64 {
	return b.ID
}

func (b Bird) StampOfData() struct {
	X      float64
	Y      float64
	Status bool
	Score  float64
	ID     float64
	Type   float64
	Size     size
} {
	return struct {
		X      float64
		Y      float64
		Status bool
		Score  float64
		ID     float64
		Type   float64
		Size     size
	}{
		b.Position.x,
		b.Position.y,
		b.GetStatus(),
		b.Score,
		b.GetID(),
		1.0,
		b.Size,
	}
}

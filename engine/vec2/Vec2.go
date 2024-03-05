package vec2

type Vec2 struct {
	X, Y float32
}

func New(x, y float32) Vec2 {
	return Vec2{X: x, Y: y}
}

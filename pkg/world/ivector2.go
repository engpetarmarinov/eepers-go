package world

// IVector2 is an integer-based 2D vector.
type IVector2 struct {
	X, Y int
}

// Add returns the sum of two vectors.
func (v IVector2) Add(other IVector2) IVector2 {
	return IVector2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub returns the difference of two vectors.
func (v IVector2) Sub(other IVector2) IVector2 {
	return IVector2{X: v.X - other.X, Y: v.Y - other.Y}
}

// Mul returns the vector scaled by an integer.
func (v IVector2) Mul(s int) IVector2 {
	return IVector2{X: v.X * s, Y: v.Y * s}
}

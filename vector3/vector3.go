package vector3

type Vector struct {
	X, Y, Z int
}

func (v Vector) Add(vo Vector) Vector {
	return Vector{
		X: v.X + vo.X,
		Y: v.Y + vo.Y,
		Z: v.Z + vo.Z,
	}
}
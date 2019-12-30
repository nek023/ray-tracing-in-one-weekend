package main

import (
	"fmt"

	"github.com/golang/geo/r3"
)

type Ray struct {
	Org r3.Vector
	Dir r3.Vector
}

func NewRay(org, dir r3.Vector) Ray {
	return Ray{org, dir}
}

func (r Ray) pointAtParameter(t float64) r3.Vector {
	return r3.Vector{
		r.Org.X + t*r.Dir.X,
		r.Org.Y + t*r.Dir.Y,
		r.Org.Z + t*r.Dir.Z,
	}
}

func color(r Ray) r3.Vector {
	unit := r.Dir.Normalize()
	t := 0.5 * (unit.Y + 1.0)
	return r3.Vector{X: 1.0, Y: 1.0, Z: 1.0}.Mul(1.0 - t).Add(r3.Vector{X: 0.5, Y: 0.7, Z: 1.0}.Mul(t))
}

func main() {
	nx := 200
	ny := 100
	fmt.Printf("P3\n%d %d\n255\n", nx, ny)
	lowerLeftCorner := r3.Vector{X: -2.0, Y: -1.0, Z: -1.0}
	horizontal := r3.Vector{X: 4.0, Y: 0.0, Z: 0.0}
	vertical := r3.Vector{X: 0.0, Y: 2.0, Z: 0.0}
	origin := r3.Vector{X: 0.0, Y: 0.0, Z: 0.0}
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			u := float64(i) / float64(nx)
			v := float64(j) / float64(ny)
			r := NewRay(origin, lowerLeftCorner.Add(horizontal.Mul(u)).Add(vertical.Mul(v)))
			col := color(r)
			ir := int(255.990 * col.X)
			ig := int(255.990 * col.Y)
			ib := int(255.990 * col.Z)
			fmt.Printf("%d %d %d\n", ir, ig, ib)
		}
	}
}

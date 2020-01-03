package main

import (
	"fmt"
	"math"

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
	t := hitSphere(r3.Vector{X: 0, Y: 0, Z: -1}, 0.5, r)
	if t > 0 {
		n := r.pointAtParameter(t).Sub(r3.Vector{X: 0, Y: 0, Z: -1}).Normalize()
		return r3.Vector{X: n.X + 1, Y: n.Y + 1, Z: n.Z + 1}.Mul(0.5)
	}
	unitDir := r.Dir.Normalize()
	t = 0.5 * (unitDir.Y + 1.0)
	return r3.Vector{X: 1.0, Y: 1.0, Z: 1.0}.Mul(1.0 - t).Add(r3.Vector{X: 0.5, Y: 0.7, Z: 1.0}.Mul(t))
}

func hitSphere(center r3.Vector, radius float64, r Ray) float64 {
	oc := r.Org.Sub(center)
	a := r.Dir.Dot(r.Dir)
	b := 2.0 * oc.Dot(r.Dir)
	c := oc.Dot(oc) - radius*radius
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return -1.0
	}
	return (-b - math.Sqrt(discriminant)) / (2.0 * a)
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

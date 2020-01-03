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

type HitRecord struct {
	T      float64
	P      r3.Vector
	Normal r3.Vector
}

type Hitable interface {
	hit(r Ray, tMin float64, tMax float64, rec *HitRecord) bool
}

type Sphere struct {
	Center r3.Vector
	Radius float64
}

func NewSphere(center r3.Vector, radius float64) Sphere {
	return Sphere{Center: center, Radius: radius}
}

func (s Sphere) hit(r Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	oc := r.Org.Sub(s.Center)
	a := r.Dir.Dot(r.Dir)
	b := oc.Dot(r.Dir)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - a*c
	if discriminant > 0 {
		temp := (-b - math.Sqrt(discriminant)) / a
		if temp > tMin && temp < tMax {
			rec.T = temp
			rec.P = r.pointAtParameter(rec.T)
			rec.Normal = rec.P.Sub(s.Center).Mul(1.0 / s.Radius)
			return true
		}
		temp = (-b + math.Sqrt(discriminant)) / a
		if temp > tMin && temp < tMax {
			rec.T = temp
			rec.P = r.pointAtParameter(rec.T)
			rec.Normal = rec.P.Sub(s.Center).Mul(1.0 / s.Radius)
			return true
		}
	}
	return false
}

type HitableList struct {
	List     []Hitable
	ListSize int
}

func NewHitableList(list []Hitable, listSize int) HitableList {
	return HitableList{List: list, ListSize: listSize}
}

func (hl HitableList) hit(r Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	var tempRec HitRecord
	hitAnything := false
	closestSoFar := tMax
	for i := 0; i < hl.ListSize; i++ {
		if hl.List[i].hit(r, tMin, closestSoFar, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T
			*rec = tempRec
		}
	}
	return hitAnything
}

func color(r Ray, world Hitable) r3.Vector {
	var rec HitRecord
	if world.hit(r, 0, math.MaxFloat64, &rec) {
		return r3.Vector{X: rec.Normal.X + 1, Y: rec.Normal.Y + 1, Z: rec.Normal.Z + 1}.Mul(0.5)
	}
	unitDir := r.Dir.Normalize()
	t := 0.5 * (unitDir.Y + 1.0)
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
	var list []Hitable
	list = append(list, NewSphere(r3.Vector{X: 0, Y: 0, Z: -1}, 0.5))
	list = append(list, NewSphere(r3.Vector{X: 0, Y: -100.5, Z: -1}, 100))
	world := NewHitableList(list, len(list))
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			u := float64(i) / float64(nx)
			v := float64(j) / float64(ny)
			r := NewRay(origin, lowerLeftCorner.Add(horizontal.Mul(u)).Add(vertical.Mul(v)))
			col := color(r, world)
			ir := int(255.990 * col.X)
			ig := int(255.990 * col.Y)
			ib := int(255.990 * col.Z)
			fmt.Printf("%d %d %d\n", ir, ig, ib)
		}
	}
}

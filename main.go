package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/golang/geo/r3"
)

func NewVector(x, y, z float64) r3.Vector {
	return r3.Vector{X: x, Y: y, Z: z}
}

type Ray struct {
	Org r3.Vector
	Dir r3.Vector
}

func NewRay(org, dir r3.Vector) Ray {
	return Ray{org, dir}
}

func (r Ray) pointAtParameter(t float64) r3.Vector {
	return NewVector(
		r.Org.X+t*r.Dir.X,
		r.Org.Y+t*r.Dir.Y,
		r.Org.Z+t*r.Dir.Z,
	)
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

type Camera struct {
	Origin          r3.Vector
	LowerLeftCorner r3.Vector
	Horizontal      r3.Vector
	Vertical        r3.Vector
}

func NewCamera(origin, lowerLeftCorner, horizontal, vertical r3.Vector) Camera {
	return Camera{
		Origin:          origin,
		LowerLeftCorner: lowerLeftCorner,
		Horizontal:      horizontal,
		Vertical:        vertical,
	}
}

func color(r Ray, world Hitable) r3.Vector {
	var rec HitRecord
	if world.hit(r, 0, math.MaxFloat64, &rec) {
		return NewVector(rec.Normal.X+1, rec.Normal.Y+1, rec.Normal.Z+1).Mul(0.5)
	}
	unitDir := r.Dir.Normalize()
	t := 0.5 * (unitDir.Y + 1.0)
	return NewVector(1.0, 1.0, 1.0).Mul(1.0 - t).Add(NewVector(0.5, 0.7, 1.0).Mul(t))
}

func (c Camera) GetRay(u, v float64) Ray {
	return NewRay(c.Origin, c.LowerLeftCorner.Add(c.Horizontal.Mul(u)).Add(c.Vertical.Mul(v)).Sub(c.Origin))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	nx := 200
	ny := 100
	ns := 100
	fmt.Printf("P3\n%d %d\n255\n", nx, ny)
	var list []Hitable
	list = append(list, NewSphere(NewVector(0, 0, -1), 0.5))
	list = append(list, NewSphere(NewVector(0, -100.5, -1), 100))
	world := NewHitableList(list, len(list))
	cam := NewCamera(
		NewVector(0.0, 0.0, 0.0),
		NewVector(-2.0, -1.0, -1.0),
		NewVector(4.0, 0.0, 0.0),
		NewVector(0.0, 2.0, 0.0),
	)
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			col := NewVector(0, 0, 0)
			for s := 0; s < ns; s++ {
				u := (float64(i) + rand.Float64()) / float64(nx)
				v := (float64(j) + rand.Float64()) / float64(ny)
				r := cam.GetRay(u, v)
				col = col.Add(color(r, world))
			}
			col = col.Mul(1.0 / float64(ns))
			ir := int(255.990 * col.X)
			ig := int(255.990 * col.Y)
			ib := int(255.990 * col.Z)
			fmt.Printf("%d %d %d\n", ir, ig, ib)
		}
	}
}

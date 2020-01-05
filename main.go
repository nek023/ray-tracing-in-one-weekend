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
	Mat    *Material
}

type Hitable interface {
	Hit(r Ray, tMin float64, tMax float64, rec *HitRecord) bool
}

type Sphere struct {
	Center r3.Vector
	Radius float64
	Mat    Material
}

func NewSphere(center r3.Vector, radius float64, mat Material) Sphere {
	return Sphere{Center: center, Radius: radius, Mat: mat}
}

func (s Sphere) Hit(r Ray, tMin float64, tMax float64, rec *HitRecord) bool {
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
			rec.Mat = &s.Mat
			return true
		}
		temp = (-b + math.Sqrt(discriminant)) / a
		if temp > tMin && temp < tMax {
			rec.T = temp
			rec.P = r.pointAtParameter(rec.T)
			rec.Normal = rec.P.Sub(s.Center).Mul(1.0 / s.Radius)
			rec.Mat = &s.Mat
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

func (hl HitableList) Hit(r Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	var tempRec HitRecord
	hitAnything := false
	closestSoFar := tMax
	for i := 0; i < hl.ListSize; i++ {
		if hl.List[i].Hit(r, tMin, closestSoFar, &tempRec) {
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

type Material interface {
	Scatter(rIn Ray, rec *HitRecord, attenuation *r3.Vector, scattered *Ray) bool
}

type Lambertian struct {
	Albedo r3.Vector
}

func NewLambertian(albedo r3.Vector) Lambertian {
	return Lambertian{Albedo: albedo}
}

func (l Lambertian) Scatter(rIn Ray, rec *HitRecord, attenuation *r3.Vector, scattered *Ray) bool {
	target := rec.P.Add(rec.Normal).Add(randomInUnitSphere())
	*scattered = NewRay(rec.P, target.Sub(rec.P))
	*attenuation = l.Albedo
	return true
}

type Metal struct {
	Albedo r3.Vector
}

func NewMetal(albedo r3.Vector) Metal {
	return Metal{Albedo: albedo}
}

func (m Metal) Scatter(rIn Ray, rec *HitRecord, attenuation *r3.Vector, scattered *Ray) bool {
	reflected := reflect(rIn.Dir.Normalize(), rec.Normal)
	*scattered = NewRay(rec.P, reflected)
	*attenuation = m.Albedo
	return scattered.Dir.Dot(rec.Normal) > 0
}

func reflect(v, n r3.Vector) r3.Vector {
	return v.Sub(n.Mul(2 * v.Dot(n)))
}

func randomInUnitSphere() r3.Vector {
	var p r3.Vector
	for {
		p = NewVector(rand.Float64(), rand.Float64(), rand.Float64()).Mul(2).Sub(NewVector(1, 1, 1))
		if p.Norm2() < 1.0 {
			break
		}
	}
	return p
}

func VectorMul(a, b r3.Vector) r3.Vector {
	return NewVector(a.X*b.X, a.Y*b.Y, a.Z*b.Z)
}

func color(r Ray, world Hitable, depth int) r3.Vector {
	var rec HitRecord
	if world.Hit(r, 0.001, math.MaxFloat64, &rec) {
		var attenuation r3.Vector
		var scattered Ray
		if depth < 50 && (*rec.Mat).Scatter(r, &rec, &attenuation, &scattered) {
			return VectorMul(attenuation, color(scattered, world, depth+1))
		}
		return NewVector(0, 0, 0)
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
	list = append(list, NewSphere(NewVector(0, 0, -1), 0.5, NewLambertian(NewVector(0.8, 0.3, 0.3))))
	list = append(list, NewSphere(NewVector(0, -100.5, -1), 100, NewLambertian(NewVector(0.8, 0.8, 0))))
	list = append(list, NewSphere(NewVector(1, 0, -1), 0.5, NewMetal(NewVector(0.8, 0.6, 0.2))))
	list = append(list, NewSphere(NewVector(-1, 0, -1), 0.5, NewMetal(NewVector(0.8, 0.8, 0.8))))
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
				col = col.Add(color(r, world, 0))
			}
			col = col.Mul(1.0 / float64(ns))
			col = NewVector(math.Sqrt(col.X), math.Sqrt(col.Y), math.Sqrt(col.Z))
			ir := int(255.990 * col.X)
			ig := int(255.990 * col.Y)
			ib := int(255.990 * col.Z)
			fmt.Printf("%d %d %d\n", ir, ig, ib)
		}
	}
}

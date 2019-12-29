package main

import (
	"fmt"

	"github.com/golang/geo/r3"
)

func main() {
	nx := 200
	ny := 100
	fmt.Printf("P3\n%d %d\n255\n", nx, ny)
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			col := r3.Vector{
				float64(i) / float64(nx),
				float64(j) / float64(ny),
				0.2,
			}
			ir := int(255.990 * col.X)
			ig := int(255.990 * col.Y)
			ib := int(255.990 * col.Z)
			fmt.Printf("%d %d %d\n", ir, ig, ib)
		}
	}
}

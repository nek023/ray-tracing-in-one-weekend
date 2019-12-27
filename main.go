package main

import "fmt"

func main() {
	nx := 200
	ny := 100
	fmt.Printf("P3\n%d %d\n255\n", nx, ny)
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			r := float64(i) / float64(nx)
			g := float64(j) / float64(ny)
			b := 0.2
			ir := int(255.990 * r)
			ig := int(255.990 * g)
			ib := int(255.990 * b)
			fmt.Printf("%d %d %d\n", ir, ig, ib)
		}
	}
}

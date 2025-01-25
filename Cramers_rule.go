package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Размерность матрицы
const N = 100

func main() {
	start := time.Now()

	//создание матрицы
	data := [N][N]float64{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			//матрица, состоящая из случайных элементов от -10 до 10
			data[i][j] = -10 + r.Float64()*20
		}
	}

	b := [N]float64{}

	for i := 0; i < N; i++ {
		//столбец, состоящий из случайных элементов от -10 до 10
		b[i] = -10 + r.Float64()*20
	}

	//printMatrix(data)

	//fmt.Println(b)

	D := det(data)

	//fmt.Println(det)

	if D == 0 {
		fmt.Println("Матрица вырождена")
		return
	}

	x := make([]float64, N)

	for i := 0; i < N; i++ {
		temp := data

		for k := 0; k < N; k++ {
			temp[k][i] = b[k]
		}

		//printMatrix(temp)

		det_i := det(temp)
		x[i] = det_i / D
	}

	fmt.Println("Решение СЛАУ")
	fmt.Println(x)

	duration := time.Since(start)
	fmt.Printf("Время выполнения кода: %v\n", duration)
}

func printMatrix(X [N][N]float64) {
	for i := range len(X) {
		for j := range len(X) {
			fmt.Printf("%10f", X[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func det(a [N][N]float64) float64 {
	det := 1.0

	for k := 0; k < N; k++ {
		iMax := k
		max := math.Abs(a[k][k])
		for i := k + 1; i < N; i++ {
			if abs := math.Abs(a[i][k]); abs > max {
				iMax = i
				max = abs
			}
		}

		if a[iMax][k] == 0 {
			return 0
		}

		if iMax != k {
			a[k], a[iMax] = a[iMax], a[k]
			det *= -1
		}

		det *= a[k][k]

		for i := k + 1; i < N; i++ {
			for j := k + 1; j < N; j++ {
				a[i][j] -= a[k][j] * (a[i][k] / a[k][k])
			}
			a[i][k] = 0
		}
	}

	return det
}

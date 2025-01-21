package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Matrix [][]float32

// Размерность матрицы
const N = 500
const M = N + 1

func main() {
	start := time.Now()

	//создание расширенной матрицы
	data := [N][M]float64{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			//матрица, состоящая из случайных элементов от -10 до 10
			data[i][j] = -10 + r.Float64()*20
		}
	}

	x, err := Gauss(data)
	if !err {
		fmt.Println("Матрица вырождена")
		return
	}
	fmt.Println("Решение СЛАУ")
	fmt.Println(x)

	duration := time.Since(start)
	fmt.Printf("Время выполнения кода: %v\n", duration)
}

func printMatrix(X [N][M]float64) {
	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			fmt.Printf("%10f", X[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func Gauss(a [N][M]float64) ([]float64, bool) {
	x := make([]float64, N)
	flag := true

	//Прямой ход
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
			flag = false
			return x, flag
		}
		a[k], a[iMax] = a[iMax], a[k]

		for i := k + 1; i < N; i++ {
			for j := k + 1; j <= N; j++ {
				a[i][j] -= a[k][j] * (a[i][k] / a[k][k])
			}
			a[i][k] = 0
		}
	}

	//Обратный ход
	for i := N - 1; i >= 0; i-- {
		x[i] = a[i][N]
		for j := i + 1; j < N; j++ {
			x[i] -= a[i][j] * x[j]
		}
		x[i] /= a[i][i]
	}
	return x, flag
}

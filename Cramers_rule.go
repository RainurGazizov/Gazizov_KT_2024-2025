package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Matrix [][]float64

func (m Matrix) Rows() int {
	return len(m)
}

func (m Matrix) Columns() int {
	return len(m[0])
}

func (m Matrix) IsSquare() bool {
	return m.Columns() == m.Rows()
}

func (m Matrix) IsMatrix() bool {
	if m.Rows() == 0 {
		return false
	}

	for _, row := range m {
		if len(m[0]) != len(row) {
			return false
		}
	}
	return true
}

func InBetween(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	} else {
		return false
	}
}

func (m Matrix) ExcludeColumn(col_index int) (Matrix, error) {

	if !InBetween(col_index, 1, m.Columns()) {
		return Matrix{}, fmt.Errorf("input not in range")
	}

	result := make(Matrix, m.Rows())
	for i, row := range m {
		for j, el := range row {
			if j == col_index-1 {
				continue
			}
			result[i] = append(result[i], el)
		}
	}
	return result, nil
}

func (m Matrix) ExcludeRow(row_index int) (Matrix, error) {
	if !InBetween(row_index, 1, m.Rows()) {
		return Matrix{}, fmt.Errorf("input not in range")
	}

	var result Matrix
	for i, r := range m {
		if i == row_index-1 {
			continue
		}
		result = append(result, r)
	}
	return result, nil
}

func (m Matrix) Det() (float64, error) {

	if !m.IsMatrix() || !m.IsSquare() {
		return -1, fmt.Errorf("determinant is not defined for the input [Matrix: %t][Square: %t]",
			m.IsMatrix(), m.IsSquare())
	}

	if m.Rows() == 2 {
		return m[0][0]*m[1][1] - m[0][1]*m[1][0], nil
	}

	// исключаем первый столбец
	partial_matrix, err := m.ExcludeRow(1)
	if err != nil {
		return -1, err
	}

	var temp float64 = 0

	// раскладываем по элементам первого столбца
	for i, el := range m[0] {

		reduced_matrix, err := partial_matrix.ExcludeColumn(i + 1)
		if err != nil {
			return -1, err
		}

		partial_det, err := reduced_matrix.Det()
		if err != nil {
			return -1, err
		}

		temp = temp + partial_det*el*math.Pow(-1, float64(i))
	}

	return temp, nil
}

// Размерность матрицы
const N = 10

func main() {
	start := time.Now()

	//создание матрицы
	data := make([][]float64, N)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < N; i++ {
		data[i] = make([]float64, N)
		for j := 0; j < N; j++ {
			//матрица, состоящая из случайных элементов от -10 до 10
			data[i][j] = -10 + r.Float64()*20
		}
	}

	b := make([]float64, N)

	for i := 0; i < N; i++ {
		//столбец, состоящий из случайных элементов от -10 до 10
		b[i] = -10 + r.Float64()*20
	}

	//printMatrix(data)

	//fmt.Println(b)

	det, _ := Matrix(data).Det()

	//fmt.Println(det)

	if det == 0 {
		fmt.Println("Матрица вырождена")
		return
	}

	x := make([]float64, N)

	for i := 0; i < N; i++ {
		temp := make([][]float64, N)

		for s := 0; s < N; s++ {
			temp[s] = make([]float64, N)
			for j := 0; j < N; j++ {
				temp[s][j] = data[s][j]
			}
		}

		for k := 0; k < N; k++ {
			temp[k][i] = b[k]
		}

		//printMatrix(temp)

		det_i, _ := Matrix(temp).Det()
		x[i] = det_i / det
	}

	fmt.Println("Решение СЛАУ")
	fmt.Println(x)

	duration := time.Since(start)
	fmt.Printf("Время выполнения кода: %v\n", duration)
}

func printMatrix(X [][]float64) {
	for i := range len(X) {
		for j := range len(X) {
			fmt.Printf("%10f", X[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

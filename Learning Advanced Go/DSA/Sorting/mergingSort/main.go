package main

import "fmt"

func merge(a, b []int) []int {
	i, j := 0, 0
	result := []int{}

	for i < len(a) && j < len(b) {
		if a[i] > b[j] {
			result = append(result, b[j])
			j++
		} else {
			result = append(result, a[i])
			i++
		}
	}

	for i < len(a) {
		result = append(result, a[i])
		i++
	}

	for j < len(b) {
		result = append(result, b[j])
		j++
	}

	return result
}

func main() {
	a := []int{1, 4, 6, 8, 20}
	b := []int{2, 5, 10, 12, 30}

	fmt.Println(merge(a, b))
}

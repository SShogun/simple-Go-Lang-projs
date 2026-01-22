package main

import "fmt"

func twoSum(nums []int, target int) []int {
	left, right := 0, len(nums)-1

	for left < right {
		sum := nums[left] + nums[right]
		if sum == target {
			return []int{left, right}
		} else if sum < target {
			left++
		} else {
			right--
		}
	}
	return nil
}

func main() {
	a := []int{1, 4, 6, 8, 20}
	results := twoSum(a, 28)
	fmt.Printf("%s %s", fmt.Sprint(results[0]), fmt.Sprint(results[1]))
}

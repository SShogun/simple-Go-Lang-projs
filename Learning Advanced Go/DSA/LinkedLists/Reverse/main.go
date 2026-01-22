package main

type Node struct {
	Val  int
	Next *Node
}

func reverseList(head *Node) *Node {
	var prev *Node
	curr := head

	for curr != nil {
		nextTemp := curr.Next
		curr.Next = prev
		prev = curr
		curr = nextTemp
	}
	return prev
}

func singleNumber(nums []int) int {
	result := 0
	for _, n := range nums {
		result ^= n
	}
	return result
}

func subsets(nums []int) [][]int {
	var result [][]int
	var backtrack func(start int, current []int)

	backtrack = func(start int, current []int) {
		// 1. Add a copy of the current path to results
		temp := make([]int, len(current))
		copy(temp, current)
		result = append(result, temp)

		// 2. Explore further choices
		for i := start; i < len(nums); i++ {
			current = append(current, nums[i]) // Choose
			backtrack(i+1, current)            // Explore
			current = current[:len(current)-1] // Undo (Backtrack)
		}
	}

	backtrack(0, []int{})
	return result
}

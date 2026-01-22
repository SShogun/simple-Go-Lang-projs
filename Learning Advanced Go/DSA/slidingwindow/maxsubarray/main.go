package main

func maxSubArray(nums []int, k int) int {
	if k <= 0 || k > len(nums) {
		return 0 // invalid window size; avoid panics
	}

	winSum, maxSum := 0, 0

	// initial window
	for i := 0; i < k; i++ {
		winSum += nums[i]
	}
	maxSum = winSum

	// sliding window: i marks the element leaving; i+3 is the entering element (k=3)
	for i := 0; i < len(nums)-k; i++ {
		// add the new element entering the window and remove the one leaving
		winSum = winSum + nums[i+3] - nums[i]
		if winSum > maxSum {
			maxSum = winSum
		}
	}
	return maxSum
}

func main() {
	nums := []int{2, 1, 5, 1, 3, 2}
	k := 3
	result := maxSubArray(nums, k)
	println(result)
}

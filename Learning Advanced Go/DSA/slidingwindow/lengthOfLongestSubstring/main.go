package main

func leng(s string) int {
	last := make(map[rune]int)
	start, best := 0, 0

	for i, ch := range s {
		if prev, ok := last[ch]; ok && prev >= start {
			start = prev + 1
		}

		if length := i - start + 1; length > best {
			best = length
		}

		last[ch] = i
	}

	return best
}

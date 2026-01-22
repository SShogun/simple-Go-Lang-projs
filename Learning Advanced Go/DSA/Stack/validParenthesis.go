package main

func validPar(s string, corr map[rune]rune) bool {
	stack := []rune{}
	for _, ch := range s {
		if ch == '(' || ch == '{' || ch == '[' {
			stack = append(stack, ch)
		} else {
			if len(stack) == 0 || stack[len(stack)-1] != corr[ch] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}
func main() {
	corr := make(map[rune]rune)
	corr[')'] = '('
	corr[']'] = '['
	corr['}'] = '{'

	s := "{[()]]"
	println(validPar(s, corr))
}

package query

func reversePattern(pattern []string) {
	for i, j := 0, len(pattern)-1; i < j; i, j = i+1, j-1 {
		pattern[i], pattern[j] = pattern[j], pattern[i]
	}
}
func MatchPattern(path []string, pattern string) bool {
	p := Parse(pattern)
	return matchPattern(path, p)
}

func matchPattern(pathParts []string, pattern []string) bool {
	reversePattern(pattern)

	var match func(int, int) bool
	match = func(idx, patIdx int) bool {
		if patIdx == len(pattern) {
			return idx == -1
		}
		if idx < 0 {
			return false
		}

		switch pattern[patIdx] {
		case "/":
			return match(idx, patIdx+1)
		case "//":
			for subIdx := idx; subIdx >= -1; subIdx-- {
				if match(subIdx, patIdx+1) {
					return true
				}
			}
			return false
		default:
			if pattern[patIdx] == pathParts[idx] || pattern[patIdx] == "*" {
				return match(idx-1, patIdx+1)
			}
			return false
		}
	}

	return match(len(pathParts)-1, 0)
}

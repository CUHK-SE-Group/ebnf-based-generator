package grammar

type Set map[string]bool

func NewSet(elements []string) Set {
	s := make(Set)
	for _, e := range elements {
		s[e] = true
	}
	return s
}

func SetDifference(a, b Set) Set {
	diff := make(Set)
	for k := range a {
		if _, ok := b[k]; !ok {
			diff[k] = true
		}
	}
	return diff
}

func SRange(input string) []string {
	var output []string
	for _, v := range input {
		output = append(output, string(v))
	}
	return output
}
func CRange(start string, end string) []string {
	var output []string
	for i := start[0]; i <= end[0]; i++ {
		output = append(output, string(i))
	}
	return output
}

// Find all non-terminal symbols in an expansion, returns as a slice.
func NonTerminals(expansion string) []string {
	return ReNonterminal.FindAllString(expansion, -1)
}

func IsNonTerminals(expansion string) bool {
	return ReNonterminal.MatchString(expansion)
}

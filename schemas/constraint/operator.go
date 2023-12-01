package constraint

import "regexp"

type Operator[Tv1 any, Tv2 any] func(v1 Var[Tv1], v2 Var[Tv2]) bool
type OpCons[Tv1 any, Tv2 any] interface {
	func(v1 Var[Tv1], v2 Var[Tv2]) bool
}

type SingleNodeOp[T any] OpCons[NodeInfo, T]

// DefinedBefore Read Operation, compare whether v1 is defined before v2
// When generating type V2, the function will check whether the V1 is defined. If not, V2 will be discarded.
func DefinedBefore(v1 Var[NodeInfo], v2 Var[NodeInfo]) bool {
	return true
}

// ArriveCntBound Restrict the max number of occurrence of a symbol
// When generating type V1, the function will check whether the occurrence of V1 exceeds V2. If exceeds, V1 will be discarded.
func ArriveCntBound(v1 Var[NodeInfo], v2 Var[int]) bool {
	return v1.Get().OccurredTimes >= v2.Get()
}

// RegexSubstitute : Write Operation, substitute a grammar node with a regexp
func RegexSubstitute(v1 Var[NodeInfo], v2 Var[regexp.Regexp]) bool {
	return true
}

// SampleSubstitute Write Operation, substitute a grammar node with a probability map
func SampleSubstitute(v1 Var[NodeInfo], v2 Var[map[string]int]) bool {
	return true
}

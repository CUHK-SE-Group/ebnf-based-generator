package constraint

import (
	"fmt"
	"testing"
)

func TestBuildConstraint(t *testing.T) {
	builder := BuildFromNode{}
	v1 := builder.Build()
	v2 := builder.Build()
	c1 := BuildConstraint(DefinedBefore, v1(), v2())
	c2 := BuildConstraint(ArriveCntBound, v1(), NewVar[int]("node", 2))

	RegisterConstraint("test", c1)
	RegisterConstraint("test", c2)
	fmt.Println(NodeHooks)
}

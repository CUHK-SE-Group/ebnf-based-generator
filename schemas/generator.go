package schemas

type Generator interface {
	Generate(constraint Constraint) (Output, error)
}

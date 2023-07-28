package schemas

type Constraint map[string]interface{}
type Generator interface {
	Generate(constraint Constraint) (Output, error)
}

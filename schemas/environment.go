package schemas

type Input map[string]interface{}
type Output map[string]interface{}
type Environment interface {
	Interact(input Input) (Output, error)
}

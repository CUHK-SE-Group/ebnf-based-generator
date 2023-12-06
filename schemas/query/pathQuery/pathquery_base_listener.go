// Code generated from ./PathQuery.g4 by ANTLR 4.13.0. DO NOT EDIT.

package pathQuery // PathQuery
import "github.com/antlr4-go/antlr/v4"

// BasePathQueryListener is a complete listener for a parse tree produced by PathQueryParser.
type BasePathQueryListener struct{}

var _ PathQueryListener = &BasePathQueryListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasePathQueryListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasePathQueryListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasePathQueryListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasePathQueryListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterQuery is called when production query is entered.
func (s *BasePathQueryListener) EnterQuery(ctx *QueryContext) {}

// ExitQuery is called when production query is exited.
func (s *BasePathQueryListener) ExitQuery(ctx *QueryContext) {}

// EnterNode is called when production Node is entered.
func (s *BasePathQueryListener) EnterNode(ctx *NodeContext) {}

// ExitNode is called when production Node is exited.
func (s *BasePathQueryListener) ExitNode(ctx *NodeContext) {}

// EnterAny is called when production Any is entered.
func (s *BasePathQueryListener) EnterAny(ctx *AnyContext) {}

// ExitAny is called when production Any is exited.
func (s *BasePathQueryListener) ExitAny(ctx *AnyContext) {}

// EnterRootNode is called when production rootNode is entered.
func (s *BasePathQueryListener) EnterRootNode(ctx *RootNodeContext) {}

// ExitRootNode is called when production rootNode is exited.
func (s *BasePathQueryListener) ExitRootNode(ctx *RootNodeContext) {}

// EnterChild is called when production Child is entered.
func (s *BasePathQueryListener) EnterChild(ctx *ChildContext) {}

// ExitChild is called when production Child is exited.
func (s *BasePathQueryListener) ExitChild(ctx *ChildContext) {}

// EnterAll is called when production All is entered.
func (s *BasePathQueryListener) EnterAll(ctx *AllContext) {}

// ExitAll is called when production All is exited.
func (s *BasePathQueryListener) ExitAll(ctx *AllContext) {}

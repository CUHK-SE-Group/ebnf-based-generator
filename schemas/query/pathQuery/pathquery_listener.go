// Code generated from ./PathQuery.g4 by ANTLR 4.13.0. DO NOT EDIT.

package pathQuery // PathQuery
import "github.com/antlr4-go/antlr/v4"

// PathQueryListener is a complete listener for a parse tree produced by PathQueryParser.
type PathQueryListener interface {
	antlr.ParseTreeListener

	// EnterQuery is called when entering the query production.
	EnterQuery(c *QueryContext)

	// EnterNode is called when entering the Node production.
	EnterNode(c *NodeContext)

	// EnterAny is called when entering the Any production.
	EnterAny(c *AnyContext)

	// EnterRootNode is called when entering the rootNode production.
	EnterRootNode(c *RootNodeContext)

	// EnterChild is called when entering the Child production.
	EnterChild(c *ChildContext)

	// EnterAll is called when entering the All production.
	EnterAll(c *AllContext)

	// ExitQuery is called when exiting the query production.
	ExitQuery(c *QueryContext)

	// ExitNode is called when exiting the Node production.
	ExitNode(c *NodeContext)

	// ExitAny is called when exiting the Any production.
	ExitAny(c *AnyContext)

	// ExitRootNode is called when exiting the rootNode production.
	ExitRootNode(c *RootNodeContext)

	// ExitChild is called when exiting the Child production.
	ExitChild(c *ChildContext)

	// ExitAll is called when exiting the All production.
	ExitAll(c *AllContext)
}

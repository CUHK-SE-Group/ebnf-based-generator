// Code generated from ./EBNFParser.g4 by ANTLR 4.13.0. DO NOT EDIT.

package ebnf // EBNFParser
import "github.com/antlr4-go/antlr/v4"

// EBNFParserListener is a complete listener for a parse tree produced by EBNFParser.
type EBNFParserListener interface {
	antlr.ParseTreeListener

	// EnterEbnf is called when entering the ebnf production.
	EnterEbnf(c *EbnfContext)

	// EnterProduction is called when entering the production production.
	EnterProduction(c *ProductionContext)

	// EnterExpr is called when entering the expr production.
	EnterExpr(c *ExprContext)

	// EnterTerm is called when entering the term production.
	EnterTerm(c *TermContext)

	// EnterQUOTE is called when entering the QUOTE production.
	EnterQUOTE(c *QUOTEContext)

	// EnterCHOICE is called when entering the CHOICE production.
	EnterCHOICE(c *CHOICEContext)

	// EnterBRACKET is called when entering the BRACKET production.
	EnterBRACKET(c *BRACKETContext)

	// EnterID is called when entering the ID production.
	EnterID(c *IDContext)

	// EnterBRACE is called when entering the BRACE production.
	EnterBRACE(c *BRACEContext)

	// EnterNone is called when entering the None production.
	EnterNone(c *NoneContext)

	// EnterREP is called when entering the REP production.
	EnterREP(c *REPContext)

	// EnterPLUS is called when entering the PLUS production.
	EnterPLUS(c *PLUSContext)

	// EnterEXT is called when entering the EXT production.
	EnterEXT(c *EXTContext)

	// EnterSUB is called when entering the SUB production.
	EnterSUB(c *SUBContext)

	// EnterIdentifier is called when entering the identifier production.
	EnterIdentifier(c *IdentifierContext)

	// ExitEbnf is called when exiting the ebnf production.
	ExitEbnf(c *EbnfContext)

	// ExitProduction is called when exiting the production production.
	ExitProduction(c *ProductionContext)

	// ExitExpr is called when exiting the expr production.
	ExitExpr(c *ExprContext)

	// ExitTerm is called when exiting the term production.
	ExitTerm(c *TermContext)

	// ExitQUOTE is called when exiting the QUOTE production.
	ExitQUOTE(c *QUOTEContext)

	// ExitCHOICE is called when exiting the CHOICE production.
	ExitCHOICE(c *CHOICEContext)

	// ExitBRACKET is called when exiting the BRACKET production.
	ExitBRACKET(c *BRACKETContext)

	// ExitID is called when exiting the ID production.
	ExitID(c *IDContext)

	// ExitBRACE is called when exiting the BRACE production.
	ExitBRACE(c *BRACEContext)

	// ExitNone is called when exiting the None production.
	ExitNone(c *NoneContext)

	// ExitREP is called when exiting the REP production.
	ExitREP(c *REPContext)

	// ExitPLUS is called when exiting the PLUS production.
	ExitPLUS(c *PLUSContext)

	// ExitEXT is called when exiting the EXT production.
	ExitEXT(c *EXTContext)

	// ExitSUB is called when exiting the SUB production.
	ExitSUB(c *SUBContext)

	// ExitIdentifier is called when exiting the identifier production.
	ExitIdentifier(c *IdentifierContext)
}

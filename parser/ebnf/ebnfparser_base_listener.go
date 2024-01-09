// Code generated from ./EBNFParser.g4 by ANTLR 4.13.0. DO NOT EDIT.

package ebnf // EBNFParser
import "github.com/antlr4-go/antlr/v4"

// BaseEBNFParserListener is a complete listener for a parse tree produced by EBNFParser.
type BaseEBNFParserListener struct{}

var _ EBNFParserListener = &BaseEBNFParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseEBNFParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseEBNFParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseEBNFParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseEBNFParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterEbnf is called when production ebnf is entered.
func (s *BaseEBNFParserListener) EnterEbnf(ctx *EbnfContext) {}

// ExitEbnf is called when production ebnf is exited.
func (s *BaseEBNFParserListener) ExitEbnf(ctx *EbnfContext) {}

// EnterProduction is called when production production is entered.
func (s *BaseEBNFParserListener) EnterProduction(ctx *ProductionContext) {}

// ExitProduction is called when production production is exited.
func (s *BaseEBNFParserListener) ExitProduction(ctx *ProductionContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseEBNFParserListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseEBNFParserListener) ExitExpr(ctx *ExprContext) {}

// EnterTerm is called when production term is entered.
func (s *BaseEBNFParserListener) EnterTerm(ctx *TermContext) {}

// ExitTerm is called when production term is exited.
func (s *BaseEBNFParserListener) ExitTerm(ctx *TermContext) {}

// EnterQUOTE is called when production QUOTE is entered.
func (s *BaseEBNFParserListener) EnterQUOTE(ctx *QUOTEContext) {}

// ExitQUOTE is called when production QUOTE is exited.
func (s *BaseEBNFParserListener) ExitQUOTE(ctx *QUOTEContext) {}

// EnterCHOICE is called when production CHOICE is entered.
func (s *BaseEBNFParserListener) EnterCHOICE(ctx *CHOICEContext) {}

// ExitCHOICE is called when production CHOICE is exited.
func (s *BaseEBNFParserListener) ExitCHOICE(ctx *CHOICEContext) {}

// EnterBRACKET is called when production BRACKET is entered.
func (s *BaseEBNFParserListener) EnterBRACKET(ctx *BRACKETContext) {}

// ExitBRACKET is called when production BRACKET is exited.
func (s *BaseEBNFParserListener) ExitBRACKET(ctx *BRACKETContext) {}

// EnterID is called when production ID is entered.
func (s *BaseEBNFParserListener) EnterID(ctx *IDContext) {}

// ExitID is called when production ID is exited.
func (s *BaseEBNFParserListener) ExitID(ctx *IDContext) {}

// EnterBRACE is called when production BRACE is entered.
func (s *BaseEBNFParserListener) EnterBRACE(ctx *BRACEContext) {}

// ExitBRACE is called when production BRACE is exited.
func (s *BaseEBNFParserListener) ExitBRACE(ctx *BRACEContext) {}

// EnterNone is called when production None is entered.
func (s *BaseEBNFParserListener) EnterNone(ctx *NoneContext) {}

// ExitNone is called when production None is exited.
func (s *BaseEBNFParserListener) ExitNone(ctx *NoneContext) {}

// EnterREP is called when production REP is entered.
func (s *BaseEBNFParserListener) EnterREP(ctx *REPContext) {}

// ExitREP is called when production REP is exited.
func (s *BaseEBNFParserListener) ExitREP(ctx *REPContext) {}

// EnterPLUS is called when production PLUS is entered.
func (s *BaseEBNFParserListener) EnterPLUS(ctx *PLUSContext) {}

// ExitPLUS is called when production PLUS is exited.
func (s *BaseEBNFParserListener) ExitPLUS(ctx *PLUSContext) {}

// EnterEXT is called when production EXT is entered.
func (s *BaseEBNFParserListener) EnterEXT(ctx *EXTContext) {}

// ExitEXT is called when production EXT is exited.
func (s *BaseEBNFParserListener) ExitEXT(ctx *EXTContext) {}

// EnterSUB is called when production SUB is entered.
func (s *BaseEBNFParserListener) EnterSUB(ctx *SUBContext) {}

// ExitSUB is called when production SUB is exited.
func (s *BaseEBNFParserListener) ExitSUB(ctx *SUBContext) {}

// EnterIdentifier is called when production identifier is entered.
func (s *BaseEBNFParserListener) EnterIdentifier(ctx *IdentifierContext) {}

// ExitIdentifier is called when production identifier is exited.
func (s *BaseEBNFParserListener) ExitIdentifier(ctx *IdentifierContext) {}

parser grammar EBNFParser;

options {
    tokenVocab = 'EBNFLexer';
}

regex : QUOTE regexContents* QUOTE;

regexContents : TEXT
               | ESCAPE
               ;

unaryOp: REP | EXT | PLUS;
binaryOp: OR;

symbol:
     ID #SubProduction
    | regex #Terminal
    ;

tmp: symbol #None
    | LPAREN expr RPAREN #SubSymbol
    | tmp unaryOp #SymbolWithUOp
    ;

expr: expr binaryOp expr+ #SymbolWithBOp
    | tmp+ #SymbolWithCat
    ;
production: ID COLON expr SEMICOLON;

ebnf : production*;

parser grammar EBNFParser;

options {
    tokenVocab = 'EBNFLexer';
}

regex : QUOTE regexContents* QUOTE;

regexContents : TEXT
               | ESCAPE
               ;

unaryOp: REP | EXT;
binaryOp: OR;

expr: expr binaryOp expr+ #SymbolWithBOp
    | tmp+ #SymbolWithCat
    ;

tmp: symbol #None
    | LPAREN expr RPAREN #SubSymbol
    | tmp unaryOp #SymbolWithUOp
    ;

symbol: 
     ID #SubProduction
    | regex #Terminal
    ;

production: ID COLON expr SEMICOLON;

ebnf : production*;

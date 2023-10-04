parser grammar EBNFParser;

options {
    tokenVocab = 'EBNFLexer';
}
ebnf: production*;

production: ID EQUAL expr SEMICOLON;

expr: term (OR term)*;

term: factor (COMMA factor)*;

factor: identifier #ID
      | STRING #STR
      | LPAREN expr RPAREN #PAREN
      | LBRACKET expr RBRACKET #BRACKET
      | LBRACE expr RBRACE  #BRACE
      | QUOTE TEXT QUOTE #QUOTE
      | factor (REP | PLUS | EXT | SUB) #RECU
      ;

identifier: ID;








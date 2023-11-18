parser grammar EBNFParser;

options {
    tokenVocab = 'EBNFLexer';
}
ebnf: production*;

production: ID EQUAL expr SEMICOLON;

expr: term (COMMA term)*;

term: factor (OR factor)*;

factor: identifier #ID
      | factor choice factor? #CHOICE
      | LBRACKET expr RBRACKET #BRACKET
      | LBRACE expr RBRACE  #BRACE
      | QUOTE TEXT QUOTE #QUOTE
      | DOUBLEQUOTE REGTEXT DOUBLEQUOTE #QUOTE
      | LPAREN expr RPAREN #None
      ;


choice: REP #REP
        | PLUS #PLUS
        | EXT #EXT
        | SUB #SUB
        ;

identifier: ID;

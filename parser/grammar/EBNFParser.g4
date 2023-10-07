parser grammar EBNFParser;

options {
    tokenVocab = 'EBNFLexer';
}
ebnf: production*;

production: ID EQUAL expr SEMICOLON;

expr: term (OR term)*;

term: factor (COMMA factor)*;

factor: identifier #ID
      | LPAREN expr RPAREN #PAREN
      | LBRACKET expr RBRACKET #BRACKET
      | LBRACE expr RBRACE  #BRACE
      | factor choice #None
      | QUOTE TEXT QUOTE #QUOTE
      | DOUBLEQUOTE REGTEXT DOUBLEQUOTE #QUOTE
      ;


choice: REP #REP
        | PLUS #PLUS
        | EXT #EXT
        | SUB #SUB
        ;

identifier: ID;









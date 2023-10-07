lexer grammar EBNFLexer;

// Tokens
LINE_COMMENT: '//' ~[\r\n]* -> skip;

COLON: ':';
LPAREN: '(';
RPAREN: ')';
LBRACKET: '[';
RBRACKET: ']';
LBRACE: '{';
RBRACE: '}';
SEMICOLON: ';';
EQUAL: '=' | '::=';

OR: '|';
SUB: '-';
REP: '*';
PLUS: '+';
EXT: '?';
COMMA: ',';
ID: [_\p{Alpha}][_\p{Alnum}]*;
WHITESPACE: [ \r\n\t]+ -> skip;
QUOTE: '\'' -> pushMode(IN_STRING);


mode IN_STRING;
DEQUOTE: '\'' -> type(QUOTE), popMode;
TEXT: (~[\\']|ESC)+;
fragment ESC: '\\' .;
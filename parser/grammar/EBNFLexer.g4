lexer grammar EBNFLexer;

// Tokens
LINE_COMMENT: '//' ~[\r\n]* -> skip;

COLON: ':';
LPAREN: '(';
RPAREN: ')';
SEMICOLON: ';';

OR : '|';
REP: '*';
EXT: '?';

QUOTE: '\'' -> pushMode(IN_STRING);

ID: [_\p{Alpha}][_\p{Alnum}]*;

WHITESPACE: [ \r\n\t]+ -> skip;

mode IN_STRING;

TEXT: ~[\\']+;
ESCAPE: '\\' .;
DEQUOTE: '\'' -> type(QUOTE), popMode;




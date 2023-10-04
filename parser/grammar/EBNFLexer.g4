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

QUOTE: '\'' -> pushMode(IN_STRING);

ID: [_\p{Alpha}][_\p{Alnum}]*;

WHITESPACE: [ \r\n\t]+ -> skip;

mode IN_STRING;

TEXT: ~[\\']+;
ESCAPE: '\\' .;
DEQUOTE: '\'' -> type(QUOTE), popMode;

DUMQUOTE: '\'';

STRING: '"' (~["\n\r])* '"';
WS: [ \t\r\n]+ -> skip;
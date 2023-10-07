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
DOUBLEQUOTE: '"' -> pushMode(IN_REGEX);


mode IN_STRING;
DEQUOTE: '\'' -> type(QUOTE), popMode;
TEXT: (~[\\']|ESC)+;
fragment ESC: '\\' .;

mode IN_REGEX;
DEDOUBLEQUOTE: '"' -> type(DOUBLEQUOTE), popMode;
REGTEXT: (~[\\"]|ESC)+;
fragment REGESC: '\\' .;

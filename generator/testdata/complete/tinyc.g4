grammar tinyc;

program     : statement ;

statement   : 'if' paren_expr statement
            | 'if' paren_expr statement 'else' statement
            | 'while' paren_expr statement
            | 'do' statement 'while' paren_expr ';'
            | '{' statement* '}'
            | expr ';'
            | ';' ;

paren_expr  : '(' expr ')' ;

expr        : test
            | ID '=' expr ;

test        : sum
            | sum '<' sum ;

sum         : term
            | sum '+' term
            | sum '-' term ;

term        : ID
            | INT
            | paren_expr ;

ID          : [a-z] ;

INT         : [0-9]+ ;
grammar PathQuery;

query       : (rootNode)? segment (pathSeparator segment)* EOF ;

segment     : NODE_ID #Node
        | '*' #Any
        ;

rootNode    : '/' ;
pathSeparator : '/' #Child
    | '//' #All
    ;

NODE_ID     : [a-zA-Z_0-9#]+ ;
WS          : [ \t\r\n]+ -> skip ;

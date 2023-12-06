#!/bin/sh
cd "$(dirname "$0")"
alias antlr4='java -Xmx500M -cp "../../../third_party/antlr-4.13.0/antlr-4.13.0-complete.jar:$CLASSPATH" org.antlr.v4.Tool'
# antlr4 -Dlanguage=Go -no-visitor -package parsing Partial.g4 -o ../parsing
antlr4 -Dlanguage=Go -package pathQuery ./PathQuery.g4 -o ../pathQuery

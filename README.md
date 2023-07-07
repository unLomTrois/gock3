# What is this?
Golang implementation of linting PDXScript files for Crusader Kings 3

## Structure
It is composed of three parts:
- lexer (tokenizer)
- parser (LL(1) parser)
- ~~linter~~ not updated yet

Lexer creates stream of tokens that is consumed by parser, catches lexical errors, like unknown tokens (e.g. you can't write `!=`)

Parser makes AST (Abstract Syntax Tree), catches syntax errors (e.g. not closed curly brace)

## Example of work:

### Raw file

![image](https://github.com/unLomTrois/gock3/assets/51882489/bc502829-7a9e-40d1-9b82-7343fb69cf01)

### AST

![image](https://github.com/unLomTrois/gock3/assets/51882489/3836d10e-6411-4b28-92aa-89120350a667)


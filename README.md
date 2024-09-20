[![Go](https://github.com/unLomTrois/gock3/actions/workflows/go.yml/badge.svg)](https://github.com/unLomTrois/gock3/actions/workflows/go.yml)

# What is this?
Golang implementation of linting PDXScript files for Crusader Kings 3

## Structure
It is composed of three parts:
- lexer (tokenizer)
- parser (LL(1) parser)
- ~~linter~~ not updated yet

Lexer creates stream of tokens that is consumed by parser, catches lexical errors, like unknown tokens (e.g. you can't write `!=`)

Parser makes AST (Abstract Syntax Tree), catches syntax errors (e.g. not closed curly brace)

## Todo
- [x] Implement Lexer
- [x] Add tests on Lexer
- [x] Implement Parser
- [x] Update Linter
- [X] Add tests on Parser
- [ ] Make a VSCode extension

## Example of work:

### Raw file

![image](https://github.com/unLomTrois/gock3/assets/51882489/1aee3cad-f633-41a9-979d-50b4280541ea)

## Linted file:

![image](https://github.com/unLomTrois/gock3/assets/51882489/9818b66d-2c2b-483e-bc7b-eb5c64cd7ab3)


## P.S.

If you have zero knowledge in parsing, I recommend watch this playlist first: [Building a Parser from scratch](https://www.youtube.com/playlist?list=PLGNbPb3dQJ_5FTPfFIg28UxuMpu7k0eT4)

Then you can read first 2-3 chapters of [Compilers: Principles, Techniques, and Tools](https://en.wikipedia.org/wiki/Compilers:_Principles,_Techniques,_and_Tools)



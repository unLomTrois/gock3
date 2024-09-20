[![Go](https://github.com/unLomTrois/gock3/actions/workflows/go.yml/badge.svg)](https://github.com/unLomTrois/gock3/actions/workflows/go.yml)

# GOCK3 - PDXScript Tools for Crusader Kings 3

**GOCK3** is a collection of tools written in Go for tokenizing, parsing, and validating PDXScript files used in [Crusader Kings 3](https://www.crusaderkings.com/). This project aims to assist mod developers by providing utilities that can analyze PDXScript code, catch errors, and improve code quality.

> This project was inspired by [ck3-tiger](https://github.com/amtep/ck3-tiger). Some concepts and code structures have been adapted with gratitude.

## Features

The project consists of three main components:

- **Lexer**: Tokenizes PDXScript code and catches lexical errors, such as unknown tokens (e.g., invalid operators like `!=`).
- **Parser**: Constructs an Abstract Syntax Tree (AST) from the token stream and catches syntax errors (e.g., unclosed curly braces).
- ~~**Linter**: Analyzes the AST for potential issues such as improper structure, deprecated syntax, or inefficient patterns, and suggests improvements to enhance the quality of PDXScript code.~~
- ~~**Validator**: Validates the AST against predefined rules to ensure code correctness and consistency.~~

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Contributing](#contributing)
- [Todo](#todo)
- [Acknowledgments](#acknowledgments)
- [Resources](#resources)
- [License](#license)

## Installation

To install **GOCK3**, ensure you have Go installed (version 1.16 or higher), then run:

```bash
go get github.com/unLomTrois/gock3
```

## Usage

You can use GOCK3 as a command-line tool or as a library in your Go projects.

### Command-Line Tool

To lint a PDXScript file:
```bash
gock3 lint your_file.txt
```

### Library
Import GOCK3 into your Go project:

```go
import "github.com/unLomTrois/gock3"
```

Use the lexer, parser, and validator in your code:
```go
tokens, err := gock3.Lexer(fileContent)
if err != nil {
    // Handle lexical errors
}

ast, err := gock3.Parser(tokens)
if err != nil {
    // Handle syntax errors
}

validationErrors := gock3.Validator(ast)
if len(validationErrors) > 0 {
    // Handle validation errors
}
```

## Examples

### Raw File

![image](https://github.com/unLomTrois/gock3/assets/51882489/1aee3cad-f633-41a9-979d-50b4280541ea)

### Linted file:

![image](https://github.com/unLomTrois/gock3/assets/51882489/9818b66d-2c2b-483e-bc7b-eb5c64cd7ab3)

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## Todo
- [x] Implement Lexer
- [x] Add tests for Lexer
- [x] Implement Parser
- [x] Add tests for Parser
- [ ] Implement Linter
- [ ] Add tests for Linter
- [ ] Implement Validator
- [ ] Add tests for Validator
- [ ] Create a VSCode extension


# Acknowledgments
- [ck3-tiger](https://github.com/amtep/ck3-tiger) - Inspiration for code concepts and structure.


## Resources
If you are new to parsing and compilers, here are some resources to get you started:

- **Video Playlist**: [Building a Parser from scratch](https://www.youtube.com/playlist?list=PLGNbPb3dQJ_5FTPfFIg28UxuMpu7k0eT4)
- **Book**: [Compilers: Principles, Techniques, and Tools](https://en.wikipedia.org/wiki/Compilers:_Principles,_Techniques,_and_Tools) (First 2-3 chapters recommended)

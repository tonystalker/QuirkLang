package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

type Token int

const (
	EOF = iota
	ILLEGAL
	IDENTIFIER
	INT

	// Operators
	ASSIGN    // =
	PLUS      // +
	MINUS     // -
	MULTIPLY  // *
	DIVIDE    // /
	MODULUS   // %
	AMPERSAND // &
	GREATER   // >
	LESSER    // <
	NOT       // !

	// Keywords
	FN
	VAR
	IF
	ELSE
	RETURN
	LOOP

	// ONE OR TWO CHARACTER TOKENS
	EQUAL             // =
	EQUAL_EQUAL       // ==
	GREATER_EQUAL     // >=
	LESS_EQUAL        // <=
	LEFT_PARENTHESIS  // (
	RIGHT_PARENTHESIS // )
	LEFT_BRACE        // {
	RIGHT_BRACE       // }
	LEFT_BRACKET      // [
	RIGHT_BRACKET     // ]
	COMMA             // ,
	COLON             // :
	DOUBLE_QUOTE      // "
	SEMICOLON         //;

)

var tokens = [...]string{
	EOF:        "EOF",
	ILLEGAL:    "ILLEGAL",
	IDENTIFIER: "IDENTIFIER",
	INT:        "INT",
	//OPERATORS
	SEMICOLON: ";",
	PLUS:      "+",
	MINUS:     "-",
	MULTIPLY:  "*",
	DIVIDE:    "/",
	MODULUS:   "%",
	AMPERSAND: "&",
	GREATER:   ">",
	LESSER:    "<",
	NOT:       "!",
	//KEYWORDS
	FN:     "FN",
	VAR:    "VAR",
	IF:     "IF",
	ELSE:   "ELSE",
	RETURN: "RETURN",
	LOOP:   "LOOP",
	//ONE OR TWO CHARACTER TOKENS
	EQUAL:             "=",
	EQUAL_EQUAL:       "==",
	GREATER_EQUAL:     ">=",
	LESS_EQUAL:        "<=",
	LEFT_PARENTHESIS:  "(",
	RIGHT_PARENTHESIS: ")",
	LEFT_BRACE:        "{",
	RIGHT_BRACE:       "}",
	LEFT_BRACKET:      "[",
	RIGHT_BRACKET:     "]",
	COMMA:             ",",
	COLON:             ":",
	DOUBLE_QUOTE:      "\"",
}

func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	line int
	col  int
}
type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, col: 0},
		reader: bufio.NewReader(reader),
	}
}
func (l *Lexer) backup() {
	if l.pos.col > 0 {
		l.pos.col--
		l.reader.UnreadRune()
	}
}
func (l *Lexer) resetPosition() {
	l.pos.col = 0
	l.pos.line++
}
func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// Reached the end of the integer literal
				return lit
			}
			// Handle other errors if necessary
			panic(err)
		}

		l.pos.col++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			// Scanned something not in the integer literal
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}
			panic(err)
		}

		l.pos.col++

		switch r {
		case '\n':
			l.resetPosition()
			// Handle the newline character, if needed.
		case ';':
			return l.pos, SEMICOLON, ";"
		case '=':
			return l.pos, EQUAL, "="
		case '+':
			return l.pos, PLUS, "+"
		case '-':
			return l.pos, MINUS, "-"
		case '*':
			return l.pos, MULTIPLY, "*"
		case '/':
			return l.pos, DIVIDE, "/"
		case '%':
			return l.pos, MODULUS, "%"
		case '&':
			return l.pos, AMPERSAND, "&"
		case '>':
			return l.pos, GREATER, ">"
		case '<':
			return l.pos, LESSER, "<"
		case '!':
			return l.pos, NOT, "!"
		case '(':
			return l.pos, LEFT_PARENTHESIS, "("
		case ')':
			return l.pos, RIGHT_PARENTHESIS, ")"
		case '{':
			return l.pos, LEFT_BRACE, "{"
		case '}':
			return l.pos, RIGHT_BRACE, "}"
		case '[':
			return l.pos, LEFT_BRACKET, "["
		case ']':
			return l.pos, RIGHT_BRACKET, "]"
		case ',':
			return l.pos, COMMA, ","
		case ':':
			return l.pos, COLON, ":"
		case '"':
			return l.pos, DOUBLE_QUOTE, "\""
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return startPos, INT, lit
			} else if unicode.IsLetter(r) {
				// backup and let lexIdent rescan the beginning of the ident
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()
				return startPos, IDENTIFIER, lit
			} else {
				return l.pos, ILLEGAL, string(r)
			}
		}
	}
}
func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// Reached the end of the identifier
				return lit
			}
			// Handle other errors if necessary
			panic(err)
		}

		l.pos.col++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			// Scanned something not in the identifier
			l.backup()
			return lit
		}
	}
}
func main() {
	file, err := os.Open("Quirk.test")
	if err != nil {
		panic(err)
	}

	lexer := NewLexer(file)
	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}

		fmt.Printf("%d:%d\t%s\t%s\n", pos.line, pos.col, tok, lit)
	}
}

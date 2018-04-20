package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var separator = " "
var eos = ";"

// var enclosers = map[string]map[byte]int {
// 	"left": map[byte]int {
// 	'[': 1,
// 	'{': 1,
// 	'(': 1,
// 	'`': 1,
// 	'"': 1,
// 	'\'': 1,
// 	'|': 1,
// 	'_': 1,
// },
// "right": map[byte]int {
// 	'[': 1,
// 	'{': 1,
// 	'(': 1,
// 	'`': 1,
// 	'"': 1,
// 	'\'': 1,
// 	'|': 1,
// 	'_': 1,
// },
// }

// 	var leftEnclosers = map[byte]int {
// 	'[': 1,
// 	'{': 1,
// 	'(': 1,
// 	'`': 1,
// 	'"': 1,
// 	'\'': 1,
// 	'|': 1,
// 	'_': 1,
// }

// var rightEnclosers = map[byte]int {
// 	'[': 1,
// 	'{': 1,
// 	'(': 1,
// 	'`': 1,
// 	'"': 1,
// 	'\'': 1,
// 	'|': 1,
// 	'_': 1,
// }

// FIXME: all of the Expect attributes are incorrect
var tokenMap = map[string]Token{
	eos: Token{
		ID:     4,
		Type:   "eos",
		String: eos,
		Expect: "eos",
	},
	separator: Token{
		ID:     1,
		Type:   "separator",
		String: separator,
		// Set the expected value to the next expected token value based on the token before this; the "last" token
	},
	"var": Token{
		ID:     2,
		Type:   "type",
		String: "var",
		// Add the expect function and then call the expect function
		Expect: "ident",
	},
	"=": Token{
		ID:     3,
		Type:   "equals",
		String: "=",
		Expect: "expr",
	},
	"+": Token{
		ID:     4,
		Type:   "sec_op",
		String: "+",
		Expect: "expr",
	},
	"-": Token{
		ID:     4,
		Type:   "sec_op",
		String: "-",
		Expect: "expr",
	},
	".+": Token{
		ID:     4,
		Type:   "sec_op",
		String: ".+",
		Expect: "expr",
	},
	".-": Token{
		ID:     4,
		Type:   "sec_op",
		String: ".-",
		Expect: "expr",
	},
	"*": Token{
		ID:     4,
		Type:   "pri_op",
		String: "*",
		Expect: "expr",
	},
	".*": Token{
		ID:     4,
		Type:   "pri_op",
		String: ".*",
		Expect: "expr",
	},
	"/": Token{
		ID:     4,
		Type:   "pri_op",
		String: "/",
		Expect: "expr",
	},
	"./": Token{
		ID:     4,
		Type:   "pri_op",
		String: "./",
		Expect: "expr",
	},
}

// type

type parseMetaLiteral struct {
	Enclosed struct {
		Value   string
		Matched bool
	}
	Period bool
	Tick   bool
}

// Char ...
type Char struct {
	Last    byte
	Current byte
	Next    byte
}

// Token ...
type Token struct {
	ID       int
	Type     string
	True     interface{}
	String   string
	Location [2]string
	Expect   string
}

// Program ...
type Program struct {
	Value       string
	Index       int
	Length      int
	Accumulator string
	EOF         bool
	EOS         bool
	Char        Char
	Tokens      []Token
}

// AddToken ...
func (p *Program) AddToken(token Token) {
	p.Tokens = append(p.Tokens, token)
}

// GetLastChar ...
func (p *Program) GetLastChar() byte {
	if p.Index-1 < 0 {
		return 0
	}
	return p.Value[p.Index-1]
}

// GetCurrentChar ...
func (p *Program) GetCurrentChar() byte {
	if p.Index > p.Length {
		return 0
	}
	return p.Value[p.Index]
}

// GetNextChar ...
func (p *Program) GetNextChar() byte {
	if p.Index+1 > p.Length {
		return 0
	}
	return p.Value[p.Index+1]
}

// Accumulate ...
func (p *Program) Accumulate() {
	p.Accumulator += string(p.GetCurrentChar())
}

// ClearAccumulator ...
func (p *Program) ClearAccumulator() {
	p.Accumulator = ""
}

// GetAccumulator ...
func (p *Program) GetAccumulator() string {
	return p.Accumulator
}

// ShiftChar ...
func (p *Program) ShiftChar() {
	p.Accumulate()

	next := p.GetNextChar()
	if next == 0 {
		p.EOF = true
		return
	}

	p.Char.Last = p.GetCurrentChar()
	p.Char.Current = next
	p.Index++
	p.Char.Next = p.GetNextChar()
}

// GetToken ...
func (p *Program) GetToken() Token {
	return tokenMap[p.GetAccumulator()]
}

// GetSeparator ...
func (p *Program) GetSeparator() Token {
	// TODO: this should actually just swallow all separators until it hits a non-separator
	if string(p.GetCurrentChar()) == separator {
		p.ShiftChar()
		p.ClearAccumulator()
		// TODO: check this for ok
		return tokenMap[separator]
	}
	return Token{}
}

// GetLiteral ...
// TODO: we will have to consider quotes here at a later date
func (p *Program) GetLiteral() Token {
	meta := parseMetaLiteral{}

	for {
		switch {
		case string(p.GetNextChar()) == ";" || string(p.GetNextChar()) == " ":
			p.ShiftChar()

			t := Token{
				ID:     -6,
				Type:   "literal",
				String: p.GetAccumulator(),
				// Location: ,
				// TODO: this has to be based on the last token
				//Expect
			}

			switch {
			// TODO: check enclosers before this
			case meta.Period:
				// FIXME: default to 64 for now
				trueValue, err := strconv.ParseFloat(p.GetAccumulator(), 64)
				if err != nil {
					// TODO: idk handle this who cares
					fmt.Println("got a fuckin err brah", err)
				}

				t.True = trueValue
				return t

				//TODO: we would look up the value before the default in current then global scope
			default:
				trueValue, err := strconv.Atoi(p.GetAccumulator())
				if err != nil {
					// TODO: this is where we would look up the value in the current variable scope and change the token
				}
				t.True = trueValue
				return t
			}

		case string(p.GetCurrentChar()) != " ":
			switch p.GetCurrentChar() {
			case '.':
				if meta.Period {
					os.Exit(666)
					//TODO: just exit for now if there is two periods
				}
				meta.Period = true
				// TODO: forget about these for now
				// case '\'':
				// 	meta.Tick = true
				// case "[":
				// 	meta.Enclosed.Value = "bracket"
				// 	meta.
				// }

			}
			p.ShiftChar()

			// TODO: need to implement someway to tell if we need to include the space in our literal (string) or not
			// TODO: we also need to do this for plus, minus, basically anything else that should separate a literal (i.e, 6 + 5, 6+5)
		}
	}
}

// GetIdent ...
func (p *Program) GetIdent() Token {
	// TODO: need to add excluded characters and logic stuff in here for
	for {
		if string(p.GetNextChar()) != separator {
			a := p.GetAccumulator()
			p.ClearAccumulator()
			return Token{
				ID:     -1,
				Type:   "ident",
				String: a,
			}
		}

		p.ShiftChar()
	}
}

// GetStatement ...
func (p *Program) GetStatement() {
	p.GetExpr()
}

// GetExpr ...
func (p *Program) GetExpr() {
	p.GetTerm()
}

// GetTerm ...
func (p *Program) GetTerm() Token {
	// factor := p.GetFactor()
	// fmt.Println(factor)
	// return Token{}
	return p.GetFactor()
}

// GetFactor ...
func (p *Program) GetFactor() Token {
	// check if there is a paren
	// check for literal
	// parse entire body of the statement
	// 	check if it is a literal
	// 	if it doesn't match then check if it is a var
	// 	check if it is a func
	//	return error

	lit := p.GetLiteral()
	if lit.ID == -6 {
		p.ClearAccumulator()
	}
	fmt.Printf("GetLiteral %#v\n", lit)

	return lit
}

func main() {
	programName := "program.expr"

	input, err := ioutil.ReadFile(programName)
	if err != nil {
		fmt.Println("error reading input program", programName)
		os.Exit(1)
	}

	fmt.Println("input:", input)

	program := Program{
		Value:  string(input),
		Length: len(input) - 1,
	}

	fmt.Println("Program:", program)
	fmt.Println("GetLastChar", string(program.GetLastChar()))
	fmt.Println("GetCurrentChar", string(program.GetCurrentChar()))
	fmt.Println("GetNextChar", string(program.GetNextChar()))
	fmt.Println()

	for {
		program.ShiftChar()
		// fmt.Println("ShiftChar", program)

		token := program.GetToken()
		if token.ID != 0 {
			program.ClearAccumulator()
		}

		// fmt.Println("Accumulator", program.GetAccumulator())
		fmt.Printf("GetToken %#v\n", token)
		// TODO: Right here we need to have a switch based on the token that comes back and call a function to "expect" a certain token
		// if token == "" {
		// 	program.ClearAccumulator()
		// } else
		// if string(program.GetNextChar()) == ";" {
		// 	// program.GetToken()
		// 	// something
		// } else
		switch {
		case token.String == "var":
			program.ClearAccumulator()

			sep := program.GetSeparator()
			if sep.String == "sep" {
				program.ClearAccumulator()
			}
			fmt.Printf("GetSeparator %#v\n", sep)
			program.ShiftChar()

			fmt.Printf("GetIdent %#v\n", program.GetIdent())
		case token.String == "=":
			sep := program.GetSeparator()
			if sep.String == "sep" {
				program.ClearAccumulator()
			}
			fmt.Printf("GetSeparator %#v\n", sep)

			program.GetExpr()
		case token.Type == "sec_op":
			sep := program.GetSeparator()
			if sep.String == "sep" {
				program.ClearAccumulator()
			}
			fmt.Printf("GetSeparator %#v\n", sep)

			program.GetExpr()
		case token.Type == "pri_op":
			sep := program.GetSeparator()
			if sep.String == "sep" {
				program.ClearAccumulator()
			}
			fmt.Printf("GetSeparator %#v\n", sep)

			program.GetTerm()
		}

		if program.EOF == true {
			// program.Accumulator += string(program.GetCurrentChar())
			fmt.Println()
			fmt.Println("Reached end-of-file")
			fmt.Println("ShiftChar", program)
			fmt.Printf("%#v\n", program)
			return
		}
	}
}

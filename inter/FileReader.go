package inter

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var builtinTokens = map[string]struct{}{
	"add": {}, "sub": {}, "mul": {}, "div": {},
	"print": {}, "println": {}, "readln": {}, "return": {},
	"equalTo": {}, "compareTo": {}, "not": {}, "and": {}, "or": {},
	"toString": {}, "typeOf": {}, "getTime": {},
	"if": {}, "/if": {}, "for": {}, "/for": {},
}

var keywordTokens = map[string]struct{}{
	"butter": {}, "run": {}, "func": {}, "/func": {},
}

func OpenFile(filename string) (BytecodeFile, []int) {
	if filepath.Ext(filename) != ".but" {
		NewError("Error: can not read file", filename, 0)
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Can't open file %s: %v\n", filename, err.Error())
		return nil, nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Can't close file %s: %v\n", filename, err.Error())
			return
		}
	}(file)

	scanner := bufio.NewScanner(file)
	var words BytecodeFile
	var lines []int
	lineCnt := 1
	for scanner.Scan() {
		text := scanner.Text()
		tokenText := tokenizeLine(text, filename, lineCnt)
		if len(tokenText) > 0 {
			compiled := compileLine(tokenText, lineCnt)
			words = append(words, compiled)
			lines = append(lines, lineCnt)
		}
		lineCnt++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err.Error())
	}

	return words, lines
}

func compileLine(tokens []string, lineCnt int) Bytecode {
	var bytecode Bytecode
	for _, token := range tokens {
		switch token {
		case "(":
			bytecode = append(bytecode, BytecodeToken{Op: OpParenOpen, Text: token, Variant: VariantToButter(token, lineCnt)})
		case ")":
			bytecode = append(bytecode, BytecodeToken{Op: OpParenClose, Text: token, Variant: VariantToButter(token, lineCnt)})
		case "=":
			bytecode = append(bytecode, BytecodeToken{Op: OpAssign, Text: token})
		default:
			if _, ok := builtinTokens[token]; ok {
				bytecode = append(bytecode, BytecodeToken{Op: OpBuiltin, Text: token})
				continue
			}
			if _, ok := keywordTokens[token]; ok {
				bytecode = append(bytecode, BytecodeToken{Op: OpKeyword, Text: token})
				continue
			}
			if inMap(token, types) {
				bytecode = append(bytecode, BytecodeToken{Op: OpType, Text: token, Variant: &ButterVariants{Type: Object, value: TypeToBType[token]}})
				continue
			}
			variant := VariantToButter(token, lineCnt)
			if variant.Type == Variable {
				bytecode = append(bytecode, BytecodeToken{Op: OpName, Text: token})
			} else {
				bytecode = append(bytecode, BytecodeToken{Op: OpLiteral, Text: token, Variant: variant})
			}
		}
	}
	return bytecode
}

func tokenizeLine(text, filename string, lineCnt int) []string {
	var tokens []string
	var current strings.Builder
	inString := false

	for i := 0; i < len(text); i++ {
		ch := text[i]
		if inString {
			current.WriteByte(ch)
			if ch == '"' {
				inString = false
				tokens = append(tokens, current.String())
				current.Reset()
			}
			continue
		}

		if i+1 < len(text) && ch == '/' && text[i+1] == '/' {
			break
		}

		switch ch {
		case ' ', '\t', ',':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		case '(', ')', '=':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(ch))
		case '"':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			inString = true
			current.WriteByte(ch)
		default:
			current.WriteByte(ch)
		}
	}

	if inString {
		NewError("Error: unmatched quotation marks", filename, lineCnt)
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}

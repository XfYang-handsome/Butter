package inter

import (
	"fmt"
	"os"
)

func EachButterFunction() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file to interpret.")
		return
	}
	filename := os.Args[1]

	var ButterFunc = new(Functions)
	var line BytecodeFile
	var eachLine []int
	var foundFunc bool

	words, lines := OpenFile(filename)
	var ifCnt int
	var ifLines []int
	var forCnt int
	var forLines []int
	var funcCnt int

	for ind, word := range words {
		if len(word) == 0 {
			continue
		}

		if word[0].Text == "run" || word[0].Text == "func" {
			if funcCnt == 0 {
				line = BytecodeFile{}
				eachLine = []int{}
			} else {
				NewError("Error: nested function definition is not allowed", filename, lines[ind])
			}
			funcCnt++
			ButterFunc = new(Functions)
			ButterFunc.name = word[1].Text
			ButterFunc.DoRun = word[0].Text == "run"
			ButterFunc.args = make(map[string]*ButterVariants)
			ButterFunc.argNames = []string{}
			for i := 3; i < len(word); i++ {
				if word[i].Text == ")" {
					break
				}
				if word[i].Text == "," || word[i].Op == OpParenOpen {
					continue
				}
				if i+1 < len(word) {
					typeToken := word[i+1].Text
					if typeToken == "int" || typeToken == "float" || typeToken == "string" || typeToken == "bool" || typeToken == "char" || typeToken == "object" {
						ButterFunc.args[word[i].Text] = &ButterVariants{Type: TypeToBType[typeToken], value: nil}
						ButterFunc.argNames = append(ButterFunc.argNames, word[i].Text)
						i++
						continue
					}
				}
				ButterFunc.args[word[i].Text] = &ButterVariants{Type: Object, value: nil}
				ButterFunc.argNames = append(ButterFunc.argNames, word[i].Text)
			}
			NameToFunctions[ButterFunc.name] = ButterFunc
			foundFunc = true
			continue
		} else if word[0].Text == "if" {
			ifCnt++
			ifLines = append(ifLines, lines[ind])
		} else if word[0].Text == "/if" {
			ifs[ifLines[ifCnt-1]] = lines[ind]
			ifLines = ifLines[:len(ifLines)-1]
			ifCnt--
		} else if word[0].Text == "for" {
			forCnt++
			forLines = append(forLines, lines[ind])
		} else if word[0].Text == "/for" {
			fors[forLines[forCnt-1]] = lines[ind]
			forLines = forLines[:len(forLines)-1]
			forCnt--
		} else if word[0].Text == "/func" {
			funcCnt--
			if funcCnt < 0 {
				NewError("Error: unmatched function", filename, lines[ind])
			}
			if len(line) > 0 {
				ButterFunctions[ButterFunc.name] = line
				ButterLines[ButterFunc.name] = eachLine
			} else {
				NewError("Error: Nothing in function "+ButterFunc.name+".", filename, eachLine[len(eachLine)-1])
			}

		}

		line = append(line, word)
		eachLine = append(eachLine, lines[ind])
	}

	if ifCnt != 0 {
		NewError("Error: unmatched \"if\"", filename, eachLine[len(eachLine)-1])
	}
	if forCnt != 0 {
		NewError("Error: unmatched \"for\"", filename, eachLine[len(eachLine)-1])
	}
	if funcCnt != 0 {
		NewError("Error: unclosed function", filename, eachLine[len(eachLine)-1])
	}

	if len(line) > 0 {
		if !foundFunc {
			NewError("Error: no function defined", filename, eachLine[len(eachLine)-1])
		}
		ButterFunctions[ButterFunc.name] = line
		ButterLines[ButterFunc.name] = eachLine
	}
}

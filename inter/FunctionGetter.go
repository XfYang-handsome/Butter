package inter

import (
	"fmt"
	"os"
)

func EachButterFunction() { //以run或者func关键字分割每个Butter函数，run在开始时运行，func不运行
	ButterFunc := new(Function)
	var line [][]string
	var eachLine []int

	if len(os.Args) < 2 {
		fmt.Println("Please provide a file to interpret.")
		return
	}
	filename := os.Args[1]

	var words, lines = OpenFile(filename)
	var ifCnt int
	var ifLines []int
	var forCnt int
	var forLines []int

	for ind, word := range words {
		if word[0] == "run" || word[0] == "func" { //如果是一个新的函数
			if len(line) > 0 { //将文件的某一行存入ButterFunctions表中
				ButterFunctions[ButterFunc.name] = line
				ButterLines[ButterFunc.name] = eachLine
				line = [][]string{}
				eachLine = []int{}
			}
			ButterFunc = new(Function)
			ButterFunc.name = word[1]
			ButterFunc.DoRun = word[0] == "run"
			ButterFunc.args = make(map[string]*ButterVariants)
			for ind, arg := range word[3:] { //Butter函数的初始参数
				if arg != "(" && arg != ")" && (!inMap(arg, types)) {
					bVar := &ButterVariants{Type: TypeToBType[word[ind+4]], value: types[word[ind+4]]}
					ButterFunc.args[arg] = bVar
				}
			}
			NameToFunctions[ButterFunc.name] = ButterFunc
			continue
		} else if word[0] == "if" { //如果是if或/if
			ifCnt++
			ifLines = append(ifLines, lines[ind])
		} else if word[0] == "/if" {
			ifs[ifLines[ifCnt-1]] = lines[ind]
			ifLines = ifLines[:len(ifLines)-1]
			ifCnt--
		} else if word[0] == "for" { //如果是for或/for
			forCnt++
			forLines = append(forLines, lines[ind])
		} else if word[0] == "/for" {
			fors[forLines[forCnt-1]] = lines[ind]
			forLines = forLines[:len(forLines)-1]
			forCnt--
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

	ButterFunctions[ButterFunc.name] = line
	ButterLines[ButterFunc.name] = eachLine
}

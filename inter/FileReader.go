package inter

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func OpenFile(filename string) ([][]string, []int) {
	var words [][]string
	var lines []int
	var lineCnt = 1

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
	for scanner.Scan() {

		text := scanner.Text()
		var str = ""
		var isPacked = false
		var eachWord []string
		var cntL int
		var cntR int

		re := regexp.MustCompile(`^\s+`)
		text = re.ReplaceAllString(text, "")

		for ind, line := range text {
			if !isPacked {

				if line == '/' && text[ind+1] == '/' {
					break
				}

				if line != ' ' && line != ',' {
					if line == '(' || line == ')' {
						if str != "" {
							eachWord = append(eachWord, str)
							str = ""
						}
						eachWord = append(eachWord, string(line))
						if line == '(' {
							cntL++
						} else {
							cntR++
						}
					} else {
						str += string(line)
					}
				} else {
					if str != "" {
						eachWord = append(eachWord, str)
						str = ""
					}
				}
			} else {
				str += string(line)
			}
			if line == '"' {
				isPacked = !isPacked
			}

		}
		if str != "" {
			eachWord = append(eachWord, str)
		}
		if len(eachWord) > 0 {
			words = append(words, eachWord)
			lines = append(lines, lineCnt)
		}
		if cntL != cntR {
			NewError("Error: incompleted or illegal expression", filename, lineCnt)
		}

		lineCnt++
	}
	return words, lines
}

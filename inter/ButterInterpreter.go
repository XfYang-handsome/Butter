package inter

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func inMap(key string, mapping map[string]any) bool {
	_, ok := mapping[key]
	return ok
}
func inMapT(key string, mapping map[string]Types) bool {
	_, ok := mapping[key]
	return ok
}
func getKey(val int, m map[int]int) (int, error) {
	for k, v := range m {
		if v == val {
			return k, nil
		}
	}
	return 0, errors.New("can not find key")
}
func getKeyT(val Types, m map[string]Types) (string, error) {
	for k, v := range m {
		if v == val {
			return k, nil
		}
	}
	return "", errors.New("can not find key")
}

func NewError(s string, f string, l int) { //创建一个新的Error
	println()
	fmt.Println(errors.New(s))
	fmt.Print("\tat file: ", f, ", line: ", l)
	println()
	syscall.Exit(0)
}

func VariantToButter(v *string, l int) *ButterVariants { //将某个字符串转换为Butter变量
	var filename = os.Args[1]

	str := *v
	if str[0] == '"' && str[len(str)-1] == '"' {
		bVar := &ButterVariants{Type: String, value: str[1 : len(str)-1]}
		return bVar
	}

	numI, errI := strconv.ParseInt(str, 10, 64)
	if errI == nil {
		bVar := &ButterVariants{Type: Int, value: numI}
		return bVar
	}

	numF, errF := strconv.ParseFloat(str, 64)
	if errF == nil {
		bVar := &ButterVariants{Type: Float, value: numF}
		return bVar
	}

	numB, errB := strconv.ParseBool(str)
	if errB == nil {
		bVar := &ButterVariants{Type: Bool, value: numB}
		return bVar
	}

	if str == ")" || str == "(" { //括号需作为单独的变量读取
		bVar := &ButterVariants{Type: String, value: str}
		return bVar
	}

	if inMapT(str, TypeToBType) {
		bVar := &ButterVariants{Type: Object, value: TypeToBType[str]}
		return bVar
	}

	NewError("Error: can not parse the element "+str, filename, l)
	return nil
}

func ButterInterpreter(words [][]string, lines []int, bFunc Function, fatherFunc *Function) { //编译某个Butter代码块

	var filename = os.Args[1]
	variants := map[string]*ButterVariants{}
	var currentLine int //函数变量，当前行
	var endIf int
	var notif, doEnd, toFor, doVariant bool

	functions := map[string]func(){ //所有Butter内置函数
		"add": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments of function \"add\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"add\"", filename, currentLine)
			}
			var args []*ButterVariants
			for bFunc.execStack.top().value != ")" {
				if bFunc.execStack.top().value != "(" {
					args = append(args, bFunc.execStack.top())
					if bFunc.execStack.top().Type != args[0].Type {
						NewError("Error: ununited type for \"add\" function", filename, currentLine)
					}
				}
				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}
			}
			bFunc.execStack.pop()
			if len(args) == 0 || len(args) == 1 {
				NewError("Error: wrong arguments for function \"add\"", filename, currentLine)
			}
			switch args[0].Type {
			case Int:
				var sum int64
				for _, arg := range args {
					sum += arg.value.(int64)
				}
				bFunc.execStack.push(&ButterVariants{Type: Int, value: sum})
			case Float:
				var sum float64
				for _, arg := range args {
					sum += arg.value.(float64)
				}
				bFunc.execStack.push(&ButterVariants{Type: Float, value: sum})
			case String:
				var sum string
				for _, arg := range args {
					sum += arg.value.(string)
				}
				bFunc.execStack.push(&ButterVariants{Type: String, value: sum})
			default:
				NewError("Error: ununited type for \"add\" function", filename, currentLine)
			}

		}, //数学
		"sub": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments of function \"sub\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"sub\"", filename, currentLine)
			}
			var args []*ButterVariants
			for bFunc.execStack.top().value != ")" {
				if bFunc.execStack.top().value != "(" {
					args = append(args, bFunc.execStack.top())
					if bFunc.execStack.top().Type != args[0].Type {
						NewError("Error: ununited type for \"sub\" function", filename, currentLine)
					}
				}
				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}
			}
			bFunc.execStack.pop()
			if len(args) == 0 || len(args) > 2 || len(args) == 1 {
				NewError("Error: wrong arguments for function \"sub\"", filename, currentLine)
			}
			switch args[0].Type {
			case Int:
				bFunc.execStack.push(&ButterVariants{Type: Int, value: args[0].value.(int64) - args[1].value.(int64)})
			case Float:

				bFunc.execStack.push(&ButterVariants{Type: Float, value: args[0].value.(float64) - args[1].value.(float64)})
			default:
				NewError("Error: ununited type for \"sub\" function", filename, currentLine)
			}
		},
		"mul": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments of function \"mul\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"mul\"", filename, currentLine)
			}
			var args []*ButterVariants
			for bFunc.execStack.top().value != ")" {
				if bFunc.execStack.top().value != "(" {
					args = append(args, bFunc.execStack.top())
					if bFunc.execStack.top().Type != args[0].Type {
						NewError("Error: ununited type for \"mul\" function", filename, currentLine)
					}
				}
				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}
			}
			bFunc.execStack.pop()
			if len(args) == 0 || len(args) == 1 {
				NewError("Error: wrong arguments for function \"mul\"", filename, currentLine)
			}
			switch args[0].Type {
			case Int:
				var mul int64 = 1
				for _, arg := range args {
					mul *= arg.value.(int64)
				}
				bFunc.execStack.push(&ButterVariants{Type: Int, value: mul})
			case Float:
				var mul = 1.0
				for _, arg := range args {
					mul *= arg.value.(float64)
				}
				bFunc.execStack.push(&ButterVariants{Type: Float, value: mul})
			default:
				NewError("Error: ununited type for \"mul\" function", filename, currentLine)
			}
		},
		"div": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments of function \"div\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"div\"", filename, currentLine)
			}
			var args []*ButterVariants
			for bFunc.execStack.top().value != ")" {
				if bFunc.execStack.top().value != "(" {
					args = append(args, bFunc.execStack.top())
					if bFunc.execStack.top().Type != args[0].Type {
						NewError("Error: ununited type for \"div\" function", filename, currentLine)
					}
				}
				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}
			}
			bFunc.execStack.pop()
			if len(args) == 0 || len(args) > 2 || len(args) == 1 {
				NewError("Error: wrong arguments for function \"div\"", filename, currentLine)
			}
			switch args[0].Type {
			case Int:
				if args[1].value.(int64) == 0 {
					NewError("Error: divided by zero", filename, currentLine)
				}
				bFunc.execStack.push(&ButterVariants{Type: Int, value: args[0].value.(int64) / args[1].value.(int64)})
			case Float:
				if args[1].value.(float64) == 0 {
					NewError("Error: divided by zero", filename, currentLine)
				}
				bFunc.execStack.push(&ButterVariants{Type: Float, value: args[0].value.(float64) / args[1].value.(float64)})
			default:
				NewError("Error: ununited type for \"div\" function", filename, currentLine)
			}
		},

		"print": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments for function \"print\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"print\"", filename, currentLine)
			}

			var args []*ButterVariants
			for bFunc.execStack.top().value != ")" {
				args = append(args, bFunc.execStack.top())
				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}
			}
			bFunc.execStack.pop()

			for i := 0; i < len(args)-1; i++ {
				if args[i].value != "(" {
					fmt.Print(args[i].value, " ")
				}
			}
			fmt.Print(args[len(args)-1].value)
		}, //IO
		"println": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments for function \"println\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"println\"", filename, currentLine)
			}

			var args []*ButterVariants

			for bFunc.execStack.top().value != ")" {

				args = append(args, bFunc.execStack.top())
				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}
			}
			bFunc.execStack.pop()

			for i := 0; i < len(args)-1; i++ {
				if args[i].value != "(" {
					fmt.Print(args[i].value, " ")
				}
			}
			fmt.Print(args[len(args)-1].value)
			println()
		},
		"readln": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments for function \"readln\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"readln\"", filename, currentLine)
			}

			for bFunc.execStack.top().value != ")" {
				bFunc.execStack.pop()
			}
			bFunc.execStack.pop()

			reader := bufio.NewReader(os.Stdin)
			str, _ := reader.ReadString('\n')
			str = strings.TrimRight(str, "\r\n")

			numI, errI := strconv.ParseInt(str, 10, 64)
			if errI == nil {
				bVar := &ButterVariants{Type: Int, value: numI}
				bFunc.execStack.push(bVar)

				return
			}

			numF, errF := strconv.ParseFloat(str, 64)
			if errF == nil {
				bVar := &ButterVariants{Type: Float, value: numF}
				bFunc.execStack.push(bVar)

				return
			}

			numB, errB := strconv.ParseBool(str)
			if errB == nil {
				bVar := &ButterVariants{Type: Bool, value: numB}
				bFunc.execStack.push(bVar)

				return
			}

			if inMapT(str, TypeToBType) {
				bVar := &ButterVariants{Type: Object, value: TypeToBType[str]}
				bFunc.execStack.push(bVar)

				return
			}

			bFunc.execStack.push(&ButterVariants{Type: String, value: str})

		},

		"return": func() {

			if bFunc.execStack.isEmpty() {
				NewError("Error: can not return nothing", filename, currentLine)
			} else if fatherFunc.name == "" {
				NewError("Error: return to an empty function", filename, currentLine)
			} else {

				es := &bFunc.execStack
				for !bFunc.execStack.isEmpty() {
					fatherFunc.execStack.push(es.s[len(es.s)-1])
					es.s = es.s[:len(es.s)-1]
				}

				doEnd = true
			}
		}, //控制流
		"if": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: no expressions for \"if\"", filename, currentLine)
			} else if bFunc.execStack.top().Type != Bool {
				NewError("Error: illegal expression for \"if\"", filename, currentLine)
			}
			if bFunc.execStack.top().value == false {
				endIf = ifs[currentLine]
				notif = true
			}
			bFunc.execStack.pop()
		},
		"for": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: no expressions for \"for\"", filename, currentLine)
			} else if bFunc.execStack.top().Type != Bool {
				NewError("Error: illegal expression for \"for\"", filename, currentLine)
			}
			if bFunc.execStack.top().value == false {
				endIf = fors[currentLine]
				notif = true
			}
			bFunc.execStack.pop()
		},
		"/for": func() {
			toFor = true

		},

		"equalTo": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments of function \"equalTo\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"equalTo\"", filename, currentLine)
			}
			var args []*ButterVariants
			for bFunc.execStack.top().value != ")" {
				if bFunc.execStack.top().value != "(" {
					args = append(args, bFunc.execStack.top())
					if bFunc.execStack.top().Type != args[0].Type {
						NewError("Error: ununited type for \"equalTo\" function", filename, currentLine)
					}
				}
				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}

			}
			bFunc.execStack.pop()
			if len(args) != 2 {
				NewError("Error: incorrect number of arguments for \"equalTo\" function", filename, currentLine)
			}
			bFunc.execStack.push(&ButterVariants{Type: Bool, value: args[0].value == args[1].value})
		}, //逻辑
		"compareTo": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments of function \"compareTo\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"compareTo\"", filename, currentLine)
			}
			var args []*ButterVariants
			for bFunc.execStack.top().value != ")" {
				if bFunc.execStack.top().value != "(" {
					args = append(args, bFunc.execStack.top())
					if bFunc.execStack.top().Type != args[0].Type {
						NewError("Error: ununited type for \"compareTo\" function", filename, currentLine)
					}
				}
				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}

			}
			bFunc.execStack.pop()
			if len(args) != 3 {
				NewError("Error: incorrect number of arguments for \"compareTo\" function", filename, currentLine)
			}
			switch args[0].Type {
			case Int:
				switch args[2].value {
				case int64(1):
					bFunc.execStack.push(&ButterVariants{Type: Bool, value: args[0].value.(int64) >= args[1].value.(int64)})
				case int64(0):
					bFunc.execStack.push(&ButterVariants{Type: Bool, value: args[0].value.(int64) > args[1].value.(int64)})
				default:
					NewError("Error: wrong arguments for function \"compareTo\"", filename, currentLine)
				}
			case Float:
				switch args[2].value {
				case int64(1):
					bFunc.execStack.push(&ButterVariants{Type: Bool, value: args[0].value.(float64) >= args[1].value.(float64)})
				case int64(0):
					bFunc.execStack.push(&ButterVariants{Type: Bool, value: args[0].value.(float64) > args[1].value.(float64)})
				default:
					NewError("Error: wrong arguments for function \"compareTo\"", filename, currentLine)
				}
			case String:
				switch args[2].value {
				case int64(1):
					bFunc.execStack.push(&ButterVariants{Type: Bool, value: args[0].value.(string) >= args[1].value.(string)})
				case int64(0):
					bFunc.execStack.push(&ButterVariants{Type: Bool, value: args[0].value.(string) > args[1].value.(string)})
				default:
					NewError("Error: wrong arguments for function \"compareTo\"", filename, currentLine)
				}
			default:
				NewError("Error: wrong type of arguments for function \"compareTo\"", filename, currentLine)
			}
		},
		"not": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments of function \"not\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"not\"", filename, currentLine)
			}
			var args []*ButterVariants
			for bFunc.execStack.top().value != ")" {
				if bFunc.execStack.top().value != "(" {
					args = append(args, bFunc.execStack.top())
				}
				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}

			}
			bFunc.execStack.pop()
			if len(args) != 1 {
				NewError("Error: incorrect number of arguments for \"not\" function", filename, currentLine)
			}
			switch args[0].Type {
			case Bool:
				bFunc.execStack.push(&ButterVariants{Type: Bool, value: !args[0].value.(bool)})
			default:
				NewError("Error: wrong type of argument for \"not\"", filename, currentLine)
			}

		},
		"and": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments of function \"and\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"and\"", filename, currentLine)
			}
			var args []*ButterVariants
			for bFunc.execStack.top().value != ")" {
				if bFunc.execStack.top().value != "(" {
					args = append(args, bFunc.execStack.top())
					if bFunc.execStack.top().Type != args[0].Type {
						NewError("Error: ununited type for \"and\" function", filename, currentLine)
					}
				}

				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}

			}
			bFunc.execStack.pop()
			if len(args) != 2 {
				NewError("Error: incorrect number of arguments for \"and\" function", filename, currentLine)
			}
			switch args[0].Type {
			case Bool:
				bFunc.execStack.push(&ButterVariants{Type: Bool, value: args[0].value.(bool) && args[1].value.(bool)})
			default:
				NewError("Error: wrong type of argument for \"and\"", filename, currentLine)
			}
		},
		"or": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments of function \"or\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"or\"", filename, currentLine)
			}
			var args []*ButterVariants
			for bFunc.execStack.top().value != ")" {
				if bFunc.execStack.top().value != "(" {
					args = append(args, bFunc.execStack.top())
					if bFunc.execStack.top().Type != args[0].Type {
						NewError("Error: ununited type for \"or\" function", filename, currentLine)
					}
				}

				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}

			}
			bFunc.execStack.pop()
			if len(args) != 2 {
				NewError("Error: incorrect number of arguments for \"or\" function", filename, currentLine)
			}
			switch args[0].Type {
			case Bool:
				bFunc.execStack.push(&ButterVariants{Type: Bool, value: args[0].value.(bool) || args[1].value.(bool)})
			default:
				NewError("Error: wrong type of argument for \"or\"", filename, currentLine)
			}
		},

		"toString": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments for function \"toString\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"toString\"", filename, currentLine)
			}

			var args []*ButterVariants
			for bFunc.execStack.top().value != ")" {
				if bFunc.execStack.top().value != "(" {
					args = append(args, bFunc.execStack.top())
				}
				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}
			}
			bFunc.execStack.pop()

			if len(args) != 1 {
				NewError("Error: wrong number of arguments for function \"toString\"", filename, currentLine)
			}

			switch args[0].Type {
			case Int:
				bFunc.execStack.push(&ButterVariants{Type: String, value: strconv.FormatInt(args[0].value.(int64), 10)})
			case Float:
				bFunc.execStack.push(&ButterVariants{Type: String, value: strconv.FormatFloat(args[0].value.(float64), 'f', -1, 64)})
			case String:
				bFunc.execStack.push(&ButterVariants{Type: String, value: args[0].value})
			case Bool:
				bFunc.execStack.push(&ButterVariants{Type: String, value: strconv.FormatBool(args[0].value.(bool))})
			case Char:
				bFunc.execStack.push(&ButterVariants{Type: String, value: string(args[0].value.(byte))})
			case Array:

			case Map:

			case Object:
				v, _ := getKeyT(args[0].value.(Types), TypeToBType)
				bFunc.execStack.push(&ButterVariants{Type: String, value: v})
			}

		}, //变量操作
		"typeOf": func() {
			if bFunc.execStack.isEmpty() {
				NewError("Error: arguments for function \"typeOf\" are missing", filename, currentLine)
			}
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"typeOf\"", filename, currentLine)
			}
			var args []*ButterVariants
			for bFunc.execStack.top().value != ")" {
				if bFunc.execStack.top().value != "(" {
					args = append(args, bFunc.execStack.top())
				}
				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}
			}
			bFunc.execStack.pop()

			if len(args) != 1 {
				NewError("Error: wrong number of arguments for function \"typeOf\"", filename, currentLine)
			}

			bFunc.execStack.push(&ButterVariants{Type: Object, value: args[0].Type})

		},

		"getTime": func() {
			if bFunc.execStack.top().value != "(" {
				NewError("Error: missing \"(\" of function \"compareTo\"", filename, currentLine)
			}
			for bFunc.execStack.top().value != ")" {
				bFunc.execStack.pop()
				if bFunc.execStack.isEmpty() {
					NewError("Error: incompleted or illegal expression", filename, currentLine)
				}

			}
			bFunc.execStack.pop()
			currentTime := &ButterVariants{value: time.Now().UnixMilli(), Type: Int}
			bFunc.execStack.push(currentTime)
		}, //其他函数
	}
	for k, v := range bFunc.args { //将函数初始参数存入变量表
		variants[k] = v
	}
	for h := 0; h < len(words); h++ {

		if notif {
			h = sort.Search(len(lines), func(i int) bool { return lines[i] >= endIf })
			notif = false
			continue
		}
		currentLine = lines[h]

		if doEnd {
			break
		}
		var doSet bool
		var leftStack Stack //左栈
		doVariant = false
		for i := len(words[h]) - 1; i >= 0; i-- { //读取一行Butter代码
			if variants[words[h][i]] != nil { //是否是已有变量
				if doSet { //如果在赋值符左侧
					leftStack.push(variants[words[h][i]])
				} else {
					bFunc.execStack.push(variants[words[h][i]])
				}
			} else if inMap(words[h][i], types) { //如果是类型
				doVariant = true
				if doSet { //如果在赋值符左侧
					leftStack.push(&ButterVariants{Type: Object, value: TypeToBType[words[h][i]]})
				} else {
					bFunc.execStack.push(&ButterVariants{Type: Object, value: TypeToBType[words[h][i]]})
				}
			} else if functions[words[h][i]] != nil { //是否是内置函数
				functions[words[h][i]]()
			} else if ButterFunctions[words[h][i]] != nil { //是否是自定义的Butter函数
				newFunc := NameToFunctions[words[h][i]]
				if bFunc.execStack.isEmpty() {
					NewError("Error: arguments of function \""+words[h][i]+"\" are missing", filename, currentLine)
				}
				if bFunc.execStack.top().value != "(" {
					NewError("Error: missing \"(\" of function \""+words[h][i]+"\"", filename, currentLine)
				}
				var fVar = make(map[string]ButterVariants)
				if newFunc.name == bFunc.name {
					for k, v := range bFunc.args {
						fVar[k] = *v
					}
				}

				var keys []string
				for k := range newFunc.args {
					keys = append(keys, k)
				}
				k := 0
				for bFunc.execStack.top().value != ")" {
					if bFunc.execStack.top().value != "(" {
						if k >= len(keys) {

							NewError("Error: wrong arguments for \""+words[h][i]+"\"", filename, currentLine)
						}
						if bFunc.execStack.top().Type != newFunc.args[keys[k]].Type {
							NewError("Error: wrong type for \""+words[h][i]+"\"", filename, currentLine)
						}
						newFunc.args[keys[k]].value = bFunc.execStack.top().value
						k++
					}
					bFunc.execStack.pop()

				}
				bFunc.execStack.pop()
				ButterInterpreter(ButterFunctions[newFunc.name], ButterLines[newFunc.name], *newFunc, &bFunc)
				if len(fVar) != 0 {
					for k, v := range fVar {
						variants[k] = &v
					}
				}

			} else if words[h][i] == "butter" { //如果是butter关键字
				if !doSet {
					bFunc.execStack.pop()
				} else {
					leftStack.s = leftStack.s[1:]
				}
			} else if words[h][i] == "/if" { //如果是/if
				//跳过
			} else if words[h][i] == "=" { //如果是赋值符
				doSet = true
			} else {
				if i > 0 && words[h][0] == "butter" && doVariant { //是否是新定义的变量
					if functions[words[h][i]] != nil {
						NewError("Error: can not use keyword as variant name", filename, currentLine)
					}

					if doSet { //如果在赋值符左侧
						variants[words[h][i]] = &ButterVariants{Type: leftStack.bottom().value.(Types), value: types[words[h][i+1]]}
						leftStack.push(variants[words[h][i]])
					} else {
						variants[words[h][i]] = &ButterVariants{Type: bFunc.execStack.bottom().value.(Types), value: types[words[h][i+1]]}
					}
				} else {
					butter := VariantToButter(&words[h][i], currentLine)
					bFunc.execStack.push(butter)
				}
			}

		}

		if doSet { //一行结束后，如果需要有赋值
			if int64(len(leftStack.s)) != bFunc.execStack.size() { //若长度不等
				NewError("Error: wrong number of given arguments", filename, currentLine)
			}
			for ind := len(leftStack.s) - 1; ind >= 0; ind-- {
				if bFunc.execStack.top().Type != leftStack.top().Type {
					if leftStack.top().Type != Object {
						NewError("Error: wrong type for variant", filename, currentLine)
					} else {
						leftStack.top().Type = bFunc.execStack.top().Type
					}
				}
				leftStack.top().value = bFunc.execStack.top().value
				bFunc.execStack.pop()
				leftStack.pop()
			}
		}
		if toFor { //一行结束后，如果需要跳转
			startFor, err := getKey(currentLine, fors)
			if err != nil {
				NewError("Error: can not find word \"for\"", filename, currentLine)
			}
			h = sort.Search(len(lines), func(i int) bool { return lines[i] >= startFor }) - 1
			toFor = false
		}

		if !bFunc.execStack.isEmpty() {
			NewError("Error: incompleted or illegal expression", filename, currentLine)
		}
	}
}

package inter

import (
	"errors"
	"fmt"
	"strconv"
	"syscall"
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
	if l != 0 {
		fmt.Print("\tat file: ", f, ", line: ", l)
	} else {
		fmt.Print("\tat file: ", f)
	}
	println()
	syscall.Exit(0)
}

func VariantToButter(str string, l int) *ButterVariants { //将某个字符串转换为Butter变量

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

	if str == "(" || str == ")" {
		bVar := &ButterVariants{Type: Function, value: str}
		return bVar
	}

	if inMapT(str, TypeToBType) {
		bVar := &ButterVariants{Type: Object, value: TypeToBType[str]}
		return bVar
	}

	Vars[str] = uint64(len(Vars))
	bVar := &ButterVariants{Type: Variable, value: str}
	return bVar
}

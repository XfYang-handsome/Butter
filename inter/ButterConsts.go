package inter

type Function struct { //Butter函数类型
	args      map[string]*ButterVariants
	name      string
	DoRun     bool
	execStack Stack
}
type ButterVariants struct { //Butter变量类型
	Type  Types
	value any
}

type Types int

const ( //Butter变量的所有类型
	Int Types = iota
	Float
	String
	Bool
	Char
	Array
	Map
	Object
)

var types = map[string]any{"int": int64(0), "float": 0.0, "string": "",
	"bool": false, "char": 0, "array": []any{}, "map": make(map[any]any), "object": Object} //Butter变量类型的默认值

var ButterFunctions = map[string][][]string{} //所有定义出来的Butter函数，key为函数名称，value为代码块

var NameToFunctions = map[string]*Function{} //以函数名称查询Function类型函数的map

var ButterLines = map[string][]int{} //所有函数的代码行数，key为函数名称，value的每个值和ButterFunctions中的每一行一一对应

var ifs = make(map[int]int)

var fors = make(map[int]int)

var ObjectFunc = new(Function)

var TypeToBType = map[string]Types{ //用字符串查询Butter变量Types的方法
	"int":    Int,
	"float":  Float,
	"string": String,
	"bool":   Bool,
	"char":   Char,
	"array":  Array,
	"map":    Map,
	"object": Object,
}

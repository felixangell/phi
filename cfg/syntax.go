package cfg

type LanguageSyntaxConfig struct {
	Syntax map[string]SyntaxCriteria `toml:"syntax"`
}

type DefaultSyntax map[string]string

var DefaultSyntaxSet = DefaultSyntax{}

func RegisterSyntax(name string, s string) {
	DefaultSyntaxSet[name] = s
}

func init() {
	// C LANGUAGE SYNTAX HIGHLIGHTING

	RegisterSyntax("c", `[syntax.c]
[syntax.type]
colouring = 0xf8f273
match = [
	"int", "char", "bool", "float", "double", "void",
	"uint8_t", "uint16_t", "uint32_t", "uint64_t",
	"int8_t", "int16_t", "int32_t", "int64_t", "const"
]

[syntax.keyword]
colouring = 0xf0a400
match = [
	"for", "break", "if", "else", "continue", "return",
	"goto", "static", "extern", "const", "typedef",
	"struct", "union", "register", "enum", 
	"do", "static", "sizeof", "volatile", "unsigned",
	"switch", "case", "goto", "default"
]

[syntax.string_literal]
colouring = 0x4b79fc
pattern = "\"([^\\\"]|\\.)*\""

[syntax.directive]
colouring = 0xf0a400
pattern = "^\\s*#\\s*include\\s+(?:<[^>]*>|\"[^\"]*\")\\s*"

[syntax.comment]
colouring = 0x4b79fc
pattern = '//.*'`)

	// GO LANGUAGE SYNTAX HIGHLIGHTING

	RegisterSyntax("go", `[syntax.go]
[syntax.keyword]
colouring = 0xf0a400
match = [
	"break", "default", "func", "interface", "select",
	"case", "defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range", "type",
	"continue", "for", "import", "return", "var",
]

[syntax.type]
colouring = 0xf8f273
match = [
	"int", "string", "uint", "rune",
	"int8", "int16", "int32", "int64",
	"uint8", "uint16", "uint32", "uint64",
	"byte", "float32", "float64", "complex64",
	"complex128", "uintptr", 
]

[syntax.comment]
colouring = 0x4b79fc
pattern = '//.*'

[syntax.string_literal]
colouring = 0x4b79fc
pattern = "\"([^\\\"]|\\.)*\""

[syntax.symbol]
colouring = 0xf0a400
match = [
	"=", ":", ";", "*", "&", "+", "-", "/", "%",
	"^", "#", "!", "@", "<", ">", ".", ","	
]`)

	RegisterSyntax("md", `[syntax.md]
[syntax.header]
colouring = 0xff00ff
pattern = '(?m)^#{1,6}.*'
`)

	// your syntax here!
}

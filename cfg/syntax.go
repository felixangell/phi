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
	RegisterSyntax("c", `[syntax.c]
[syntax.type]
colouring = 0xf8f273
match = [
	"int", "char", "bool", "float", "double", "void",
	"uint8_t", "uint16_t", "uint32_t", "uint64_t",
	"int8_t", "int16_t", "int32_t", "int64_t"
]

[syntax.keyword]
colouring = 0xf0a400
match = [
	"for", "break", "if", "else", "continue", "return",
	"goto", "static", "extern", "const", "typedef",
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

	RegisterSyntax("go", `[syntax.go]
[syntax.go.keyword]
colouring = 0xf0a400
match = [
	"type", "import", "package", "func", "struct",
	"append", "delete", "make", "for", "if", "while",
	"switch", "select", "chan", "else", "var", "const",
	"iota", "case"
]

[syntax.go.type]
colouring = 0xf8f273
match = [
	"int", "string", "uint",
	"int8", "int16", "int32", "int64",
	"uint8", "uint16", "uint32", "uint64",
	"rune", "byte", "float32", "float64"
]

[syntax.go.comment]
colouring = 0x4b79fc
pattern = '//.*'

[syntax.go.string_literal]
colouring = 0x4b79fc
pattern = "\"([^\\\"]|\\.)*\""

[syntax.go.symbol]
colouring = 0xf0a400
match = [
	"=", ":", ";", "*", "&", "+", "-", "/", "%",
	"^", "#", "!", "@", "<", ">", ".", ","	
]`)
}

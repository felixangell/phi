package cfg

type LanguageSyntaxConfig struct {
	Syntax map[string]*SyntaxCriteria `toml:"syntax"`
}

func MarkdownConfig() *LanguageSyntaxConfig {
	return &LanguageSyntaxConfig{Syntax: map[string]*SyntaxCriteria{
		"header": {
			Foreground: 0xff00ff,
			Pattern:    "(?m)^#{1,6}.*",
		},
	}}
}

func TOMLConfig() *LanguageSyntaxConfig {
	return &LanguageSyntaxConfig{Syntax: map[string]*SyntaxCriteria{
		"declaration": {
			Foreground: 0xf8f273,
			Pattern:    `(\[)(.*)(\])`,
		},
		"identifier": {
			Foreground: 0xf0a400,
			Pattern:    `\b([a-z]|[A-Z])+(_|([a-z]|[A-Z])+)*\b`,
		},
		"symbol": {
			Match:      []string{"="},
			Foreground: 0xf8f273,
		},
	}}
}

func CConfig() *LanguageSyntaxConfig {
	return &LanguageSyntaxConfig{Syntax: map[string]*SyntaxCriteria{
		"type": {
			Foreground: 0xf8f273,
			Match: []string{
				"int", "char", "bool", "float", "double", "void",
				"uint8_t", "uint16_t", "uint32_t", "uint64_t",
				"int8_t", "int16_t", "int32_t", "int64_t", "const",
			},
		},
		"keyword": {
			Foreground: 0xf0a400,
			Match: []string{
				"for", "break", "if", "else", "continue", "return",
				"goto", "extern", "const", "typedef",
				"struct", "union", "register", "enum",
				"do", "static", "sizeof", "volatile", "unsigned",
				"switch", "case", "default",
			},
		},
		"string_literal": {
			Foreground: 0x4b79fc,
			Pattern:    `\"([^\\\"]|\\.)*\"`,
		},
		"directive": {
			Foreground: 0xf0a400,
			Pattern:    `^\\s*#\\s*include\\s+(?:<[^>]*>|\"[^\"]*\")\\s*`,
		},
		"symbol": {
			Foreground: 0xf0a400,
			Match: []string{
				"+=", "-=", "*=", "/=", ">>", "<<", "==", "!=",
				">=", "<=", "||", "&&",
				"=", ":", ";", "*", "&", "+", "-", "/", "%",
				"^", "#", "!", "@", "<", ">", ".", ",",
			},
		},
		"comment": {
			Foreground: 0x4b79fc,
			Pattern:    `//.*`,
		},
	}}
}

func GoConfig() *LanguageSyntaxConfig {
	return &LanguageSyntaxConfig{Syntax: map[string]*SyntaxCriteria{
		"keyword": {
			Foreground: 0xf0a400,
			Match: []string{
				"break", "default", "func", "interface", "select",
				"case", "defer", "go", "map", "struct",
				"chan", "else", "goto", "package", "switch",
				"const", "fallthrough", "if", "range", "type",
				"continue", "for", "import", "return", "var",
			},
		},
		"type": {
			Foreground: 0xf8f273,
			Match: []string{
				"int", "string", "uint", "rune",
				"int8", "int16", "int32", "int64",
				"uint8", "uint16", "uint32", "uint64",
				"byte", "float32", "float64", "complex64",
				"complex128", "uintptr",
			},
		},
		"comment": {
			Foreground: 0x4b79fc,
			Pattern:    `//.*`,
		},
		"string_literal": {
			Foreground: 0x4b79fc,
			Pattern:    `\"([^\\\"]|\\.)*\"`,
		},
		"symbol": {
			Foreground: 0xf0a400,
			Match: []string{
				"+=", "-=", "*=", "/=", ">>", "<<", "==", "!=", ":=",
				">=", "<=", "||", "&&",
				"=", ":", ";", "*", "&", "+", "-", "/", "%",
				"^", "#", "!", "@", "<", ">", ".", ",",
			},
		},
	}}
}

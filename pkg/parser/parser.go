package parser

type Parser interface {
	// Parse(data []byte) (interface{}, error)
	IsCodeSyntax(code string) bool
	InterceptSyntaxErrorCode(choicesText, prefix, suffix string) string
	ExtractAccurateBlockPrefixSuffix(prefix, suffix string) (string, string)
}

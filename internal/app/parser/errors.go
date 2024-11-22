package parser

// Expected token error messages
const (
	errUnexpectedEOF   = "Unexpected end of input, expected one of: %s"
	errUnexpectedToken = "Unexpected token %q of type %q, expected one of: %s"

	// Additional error messages
	errFieldListUnexpectedToken = "[FieldList] Unexpected token %q of type %q"
	errUnexpectedFieldToken     = "Unexpected token %q of type %q when expecting a field"
	errOperatorExpectedEOF      = "Expected an operator '=', '==', or comparison, but reached end of input"
	errOperatorUnexpectedToken  = "Expected operator '=', '==', or comparison, but found %q of type %q"
	errValueExpectedEOF         = "Expected a value, but reached end of input"
	errValueUnexpectedToken     = "[Value] Unexpected token %q of type %q"
	errBlockUnexpectedToken     = "[Block] Unexpected token %q of type %q in block"
	errTokenListUnexpectedToken = "[TokenList] Unexpected token %q of type %q in token list"
	errLiteralExpectedEOF       = "Unexpected end of input when expecting a literal value"
	errLiteralUnexpectedToken   = "Unexpected token %q of type %q when expecting a literal value (word, number, boolean, or quoted string)"
	errRecoveredNonLiteralToken = "Recovered to non-literal token %q of type %q after error"
	errFailedUnquoteString      = "Failed to unquote string %q"
)

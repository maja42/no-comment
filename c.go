package nocomment

type state interface {
	Next(in rune) state
	Result() []rune
}

type cState int

const (
	normalState cState = iota
	normalEscapeState
	normalSlashState
	normalEscapedSlashState
	quoteState
	quoteEscapeState
	lineCommentState
	blockCommentState
	blockCommentStarState
)

// StripCStyleComments removes C-Style line- and block comments from the provided string
// C-Style comments have two different variations://
//	- line comments start with // and end at the next \n
//  - block comments start with /* and end with */
//
// Comments cannot start within double-quoted regions. These regions can span multiple lines if the
// new-line character \n is escaped with a backslash.
// Backslashes can also be escaped with backslashes.
func StripCStyleComments(input string) string {
	var s, out string
	state := normalState

	for _, r := range input {
		s, state = parseNextRune(state, r)
		out += s
	}
	if state == normalSlashState {
		out += "/"
	}

	return out
}

func parseNextRune(state cState, r rune) (string, cState) {
	rStr := string(r)

	switch state {

	case normalState:
		switch r {
		case '\\':
			return rStr, normalEscapeState
		case '"':
			return rStr, quoteState
		case '/':
			return "", normalSlashState
		}
		return rStr, state

	case normalEscapeState:
		if r == '/' {
			return "", normalEscapedSlashState
		}
		return rStr, normalState

	case normalEscapedSlashState:
		switch r {
		case '/':
			return "", lineCommentState
		case '*':
			return "", blockCommentState
		case '"':
			return "/" + rStr, quoteState
		}
		return "/" + rStr, normalState

	case quoteState:
		switch r {
		case '\\':
			return rStr, quoteEscapeState
		case '"':
			return rStr, normalState
		case '\n':
			return rStr, normalState
		}
		return rStr, state

	case quoteEscapeState:
		return rStr, quoteState

	case normalSlashState:
		switch r {
		case '/':
			return "", lineCommentState
		case '*':
			return "", blockCommentState
		}
		return "/" + rStr, normalState

	case lineCommentState:
		if r == '\n' {
			return rStr, normalState
		}
		return "", state

	case blockCommentState:
		if r == '*' {
			return "", blockCommentStarState
		}
		return "", state

	case blockCommentStarState:
		switch r {
		case '/':
			return "", normalState
		case '*':
			return "", state
		}
		return "", blockCommentState
	}

	return "", state // cannot reach
}

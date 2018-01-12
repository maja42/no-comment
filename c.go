package commentstripper

type state interface {
	Next(in rune) state
	Result() []rune
}

type stateMemory struct {
	emitedRunes  []rune // Runes that should not be stripped
	pendingRunes []rune // Runes where we don't know yet if they should be stripped
}

func (m *stateMemory) emitRune(r rune) *stateMemory {
	m.emitedRunes = append(m.emitedRunes, r)
	return m
}

func (m *stateMemory) emitPendingRunes() *stateMemory {
	m.emitedRunes = append(m.emitedRunes, m.pendingRunes...)
	m.pendingRunes = nil
	return m
}

func (m *stateMemory) pendingRune(r rune) *stateMemory {
	if m.pendingRunes == nil {
		m.pendingRunes = []rune{r}
	} else {
		m.pendingRunes = append(m.pendingRunes, r)
	}
	return m
}

func (m *stateMemory) clearPendingRunes() *stateMemory {
	m.pendingRunes = nil
	return m
}

func (m stateMemory) Result() []rune {
	mem := m.emitPendingRunes()
	return mem.emitedRunes
}

// initial state, nothing special. just emit characters
type normalState struct{ *stateMemory }

// we are in the normal state and just encountered an escape-character
type normalEscapeState struct{ *stateMemory }

// we are in the normal state and just encountered a slash. Comments might follow
type normalSlashState struct{ *stateMemory }

// we are within a quoted string
type quoteState struct{ *stateMemory }

//  we are within a quoted string and just encountered an escape-character
type quoteEscapeState struct{ *stateMemory }

// we are within a line comment
type lineCommentState struct{ *stateMemory }

// we are within a block-comment
type blockCommentState struct{ *stateMemory }

// we are within a block-comment and just encountered a star-character
type blockCommentStarState struct{ *stateMemory }

// newStateMachine initializes a new state-machine for stripping comments
func newStateMachine() state {
	mem := stateMemory{
		emitedRunes:  make([]rune, 0),
		pendingRunes: nil,
	}
	return normalEscapeState{&mem}
}

func (s normalState) Next(in rune) state {
	if in == '\\' {
		return normalEscapeState{s.emitRune(in)}
	}
	if in == '"' {
		return quoteState{s.emitRune(in)}
	}
	if in == '/' {
		return normalSlashState{s.pendingRune(in)}
	}
	s.emitRune(in)
	return s
}

func (s normalEscapeState) Next(in rune) state {
	return normalState{s.emitRune(in)}
}

func (s quoteState) Next(in rune) state {
	if in == '\\' {
		return quoteEscapeState{s.emitRune(in)}
	}
	if in == '"' {
		return normalState{s.emitRune(in)}
	}
	if in == '\n' {
		return normalState{s.emitRune(in)}
	}
	s.emitRune(in)
	return s
}

func (s quoteEscapeState) Next(in rune) state {
	return quoteState{s.emitRune(in)}
}

func (s normalSlashState) Next(in rune) state {
	if in == '/' {
		return lineCommentState{s.clearPendingRunes()}
	}
	if in == '*' {
		return blockCommentState{s.clearPendingRunes()}
	}
	return normalState{s.emitPendingRunes().emitRune(in)}
}

func (s lineCommentState) Next(in rune) state {
	if in == '\n' {
		return normalState{s.emitRune(in)}
	}
	return s
}

func (s blockCommentState) Next(in rune) state {
	if in == '*' {
		return blockCommentStarState{s.stateMemory}
	}
	return s
}

func (s blockCommentStarState) Next(in rune) state {
	if in == '/' {
		return normalState{s.stateMemory}
	}
	if in == '*' {
		return s
	}
	return blockCommentState{s.stateMemory}
}

// StripCStyleComments removes C-Style line- and block comments from the provided string
// The method uses a state-machine
func StripCStyleComments(input string) string {
	stateMachine := newStateMachine()

	for _, r := range input {
		stateMachine = stateMachine.Next(r)
	}

	return string(stateMachine.Result())
}

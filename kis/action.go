package kis

// Action defines the actions for KisFlow execution.
type Action struct {
	// DataReuse indicates whether to reuse data from the upper Function.
	DataReuse bool

	// ForceEntryNext overrides the default rule, where if the current Function's calculation result is 0 data entries,
	// subsequent Functions will not continue execution.
	// With ForceEntryNext set to true, the next Function will be entered regardless of the data.
	ForceEntryNext bool

	// JumpFunc specifies the Function to jump to for continued execution.
	// (Note: This can easily lead to Flow loop calls, causing an infinite loop.)
	JumpFunc string

	// Abort terminates the execution of the Flow.
	Abort bool
}

// ActionFunc is the type for KisFlow Functional Option.
type ActionFunc func(ops *Action)

// LoadActions loads Actions and sequentially executes the ActionFunc operations.
func LoadActions(acts []ActionFunc) Action {
	action := Action{}

	if acts == nil {
		return action
	}

	for _, act := range acts {
		act(&action)
	}

	return action
}

// ActionDataReuse sets the option for reusing data from the upper Function.
func ActionDataReuse(act *Action) {
	act.DataReuse = true
}

// ActionForceEntryNext sets the option to forcefully enter the next layer.
func ActionForceEntryNext(act *Action) {
	act.ForceEntryNext = true
}

// ActionJumpFunc returns an ActionFunc function and sets the funcName to Action.JumpFunc.
// (Note: This can easily lead to Flow loop calls, causing an infinite loop.)
func ActionJumpFunc(funcName string) ActionFunc {
	return func(act *Action) {
		act.JumpFunc = funcName
	}
}

// ActionAbort terminates the execution of the Flow.
func ActionAbort(action *Action) {
	action.Abort = true
}

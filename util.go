package errors

func BaseStackTrace(err error) []uintptr {
	if err == nil {
		return nil
	}
	return baseStackTrace(err)
}

func baseStackTrace(err error) []uintptr {
	x, ok := err.(interface{ StackTrace() []uintptr })
	if !ok {
		return nil
	}
	stackTrace := x.StackTrace()
	childStackTrace := baseStackTrace(Unwrap(err))
	if len(childStackTrace) == 0 {
		return stackTrace
	}
	return childStackTrace
}

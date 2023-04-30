package testing

type T struct {
	success         int
	failure         int
	failureMessages []string
}

func (t *T) Assert(condition bool, failureMessage string) {
	if condition {
		t.success++
	} else {
		t.failure++
		t.failureMessages = append(t.failureMessages, failureMessage)
	}
}

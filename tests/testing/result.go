package testing

import "fmt"

type R struct {
	Success         int
	Failure         int
	FailureMessages []string
	Errors          int
	ErrorMessages   []string
}

func (r *R) MergeResults(result *R) {
	if result == nil {
		return
	}

	r.Success += result.Success
	r.Failure += result.Failure
	r.Errors += result.Errors
	r.ErrorMessages = append(r.ErrorMessages, result.ErrorMessages...)
	r.FailureMessages = append(r.FailureMessages, result.FailureMessages...)
}

func (r *R) Add(test *T) {
	r.Success += test.success
	r.Failure += test.failure
	r.FailureMessages = append(r.FailureMessages, test.failureMessages...)
}

func (r *R) IsSucceed() bool {
	return r.Errors == 0 && r.Failure == 0
}

func (r *R) Print() {
	fmt.Println("Self testing results:")
	fmt.Printf("Successes: %d \n", r.Success)
	fmt.Println()
	fmt.Println("Failures:")
	r.formatList(r.FailureMessages)
	fmt.Println()
	fmt.Println("Errors:")
	r.formatList(r.ErrorMessages)
	fmt.Println()
}

func (*R) formatList(list []string) {
	for _, s := range list {
		fmt.Printf("  * %s \n", s)
	}
}

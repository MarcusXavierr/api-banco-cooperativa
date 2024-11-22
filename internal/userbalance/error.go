package userbalance

type financialStatusError struct {
}

func (f financialStatusError) HasExceededLimit() bool {
	return true
}

func (f financialStatusError) Error() string {
	return "Credit Limit exceeded"
}

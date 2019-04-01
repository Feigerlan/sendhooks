package models

import "fmt"

var (
	Err_responses map[string]*Response
)

type Response struct {
	code   int64
	msg    string
}

func init() {
	Err_responses = make(map[string]*Response)
	Err_responses["err1"] = &Response{401, "adfa"}
	Err_responses["err2"] = &Response{402, "adsfasdf"}
}
func ErrAll() map[string]*Response {
	fmt.Println()
	return Err_responses
}


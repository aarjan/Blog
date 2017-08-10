package main

import (
	"fmt"
)

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Err  error
	Code int
}

// statisfies error interface
func (s *StatusError) Error() string {
	return s.Err.Error()
}

// statisfies Error interface
func (s *StatusError) Status() int {
	return s.Code
}

func main() {
	fmt.Printf("%v", returnError().(Error))

}

func returnError() error {
	var err Error
	err = &StatusError{Err: fmt.Errorf("eror %s", "eh"), Code: 404}
	return err
}

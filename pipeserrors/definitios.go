package pipeserrors

import "fmt"

var ErrorIsClosed error = fmt.Errorf("pipe is closed")

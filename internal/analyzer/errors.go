package analyzer

import (
	"errors"
	"fmt"
)

var ErrFileAccess = errors.New("file access error")
var ErrParse = errors.New("parse error")

type FileAccessError struct {
	Path string
	Err  error
}

func (e *FileAccessError) Error() string {
	return fmt.Sprintf("file access error on %s: %v", e.Path, e.Err)
}
func (e *FileAccessError) Unwrap() error {
	return e.Err
}

type ParseError struct {
	Path string
	Why  string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%v: %s: %s", ErrParse, e.Path, e.Why)
}

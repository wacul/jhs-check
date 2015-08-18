package walker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type ErrorsWalker struct {
	errors map[string]error
	walker filepath.WalkFunc
}

func NewErrorsWalker(walker filepath.WalkFunc) ErrorsWalker {
	return ErrorsWalker{map[string]error{}, walker}
}

var indenter = regexp.MustCompile(`(?:\r?\n|\r)\s*`)

func (w *ErrorsWalker) Err() error {
	if len(w.errors) == 0 {
		return nil
	}
	message := ""
	for path, err := range w.errors {
		message += fmt.Sprintf("%s\n  %s\n\n", path, indenter.ReplaceAllString(err.Error(), "\n  "))
	}
	return errors.New(message)
}

func (w *ErrorsWalker) Walk(path string, info os.FileInfo, err error) error {
	if w == nil || w.walker == nil {
		return err
	}
	if err := w.walker(path, info, err); err != nil {
		if err == filepath.SkipDir {
			return err
		}
		w.errors[path] = err
	}
	return nil
}

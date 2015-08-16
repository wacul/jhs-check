package walker

import (
	"os"
	"path/filepath"
	"regexp"
)

type PatternWalker struct {
	pattern *regexp.Regexp
	walker  filepath.WalkFunc
}

func NewPatternWalker(pattern *regexp.Regexp, walker filepath.WalkFunc) PatternWalker {
	return PatternWalker{pattern, walker}
}

func (w *PatternWalker) Walk(path string, info os.FileInfo, err error) error {
	if w == nil || w.walker == nil {
		return err
	}
	if w.pattern == nil || w.pattern.MatchString(info.Name()) {
		return w.walker(path, info, err)
	}
	return err
}

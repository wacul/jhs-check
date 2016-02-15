package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/wacul/jhs-check/config"
	"github.com/wacul/jhs-check/hyperschema"
	"github.com/wacul/jhs-check/walker"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	conf := config.Config{}
	app := kingpin.New("jhs-check", "JSON Hyper-schema checker")
	app.Flag("pattern", "file name pattern").Short('p').StringVar(&conf.Pattern)
	app.Arg("source", "source pathspec").Required().ExistingFilesOrDirsVar(&conf.Sources)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if conf.Sources == nil {
		return
	}

	var set hyperschema.SchemaSet

	var re *regexp.Regexp
	if conf.Pattern == "" {
		re = nil
	} else {
		re = regexp.MustCompile(conf.Pattern)
	}

	code := 0
	pw := walker.NewPatternWalker(re, set.Collect)
	ew := walker.NewErrorsWalker(pw.Walk)
	for _, file := range conf.Sources {
		if err := filepath.Walk(file, ew.Walk); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			code |= 4
		}
	}
	if err := ew.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		code |= 2
	}
	if err := set.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		code |= 1
	}

	os.Exit(code)
}

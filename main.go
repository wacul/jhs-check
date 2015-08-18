package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/davecgh/go-spew/spew"
	"github.com/kyoh86/jhs-check/config"
	"github.com/kyoh86/jhs-check/hyperschema"
	"github.com/kyoh86/jhs-check/walker"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	conf := config.Config{}
	app := kingpin.New("jhs-check", "JSON Hyper-schema checker")
	app.Flag("pattern", "file name pattern").Short('p').StringVar(&conf.Pattern)
	app.Arg("source", "source pathspec").Required().ExistingFilesOrDirsVar(&conf.Sources)
	app.Parse(os.Args[1:])

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

	pw := walker.NewPatternWalker(re, set.Collect)
	ew := walker.NewErrorsWalker(pw.Walk)
	for _, file := range conf.Sources {
		if err := filepath.Walk(file, ew.Walk); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
	if err := ew.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	set.Walk(func(schema *hyperschema.Schema, err error) error {
		if schema != nil {
			fmt.Println(schema.Types)
		}
		return nil
	})
	spew.Dump(set)

	// if err := set.Validate(); err != nil {
	// 	fmt.Fprintln(os.Stderr, err.Error())
	// }

}

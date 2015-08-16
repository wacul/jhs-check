package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/davecgh/go-spew/spew"
	"github.com/kyoh86/jhs-check/config"
	"github.com/kyoh86/jhs-check/schema"
	"github.com/kyoh86/jhs-check/walker"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

func main() {
	conf := config.Config{}
	app := kingpin.New("jhs-check", "JSON Hyper-schema checker")
	// app.Flag("format", "Schema file format").Default("yaml").EnumVar(&conf.Format, "yaml", "json")
	app.Flag("pattern", "file name pattern").Short('p').StringVar(&conf.Pattern)
	app.Arg("source", "source pathspec").Required().ExistingFilesOrDirsVar(&conf.Sources)
	app.Parse(os.Args[1:])

	if conf.Sources == nil {
		return
	}

	re := regexp.MustCompile(conf.Pattern)
	wl := walker.NewPatternWalker(re, Walk)
	for _, file := range conf.Sources {
		if err := filepath.Walk(file, wl.Walk); err != nil {
			panic(err)
		}
	}

}

func Walk(path string, info os.FileInfo, err error) error {
	fmt.Println("hoge0", err)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}

	fmt.Println("hoge1", err)
	var sc schema.Schema
	fileBuf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(fileBuf, sc); err != nil {
		return err
	}

	spew.Dump(fileBuf)
	return nil
}

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasepe/map/builder"
	"gopkg.in/yaml.v2"
)

const (
	banner = ` _ _  _  _ 
| | |(_||_)
        | `
)

var (
	optFormat  string
	optVersion bool
	gitCommit  string
)

func main() {
	configureFlags()

	if optVersion {
		fmt.Printf("%s version: %s\n", appName(), gitCommit)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.CommandLine.Usage()
		os.Exit(2)
	}

	spec := strings.Join(flag.Args(), " ")
	res, err := builder.Parse(spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	for _, g := range res {
		if strings.EqualFold(optFormat, "json") {
			dat, err := json.Marshal(g.Gen())
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}

			var out bytes.Buffer
			if err := json.Indent(&out, dat, "", "   "); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(out.Bytes()))
		} else {
			dat, err := yaml.Marshal(g.Gen())
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(string(dat))
		}
	}
}

func configureFlags() {
	name := appName()

	flag.CommandLine.Usage = func() {
		fmt.Printf("%s\n", banner)
		fmt.Print("Generate JSON or YAML from the command-line.\n\n")

		fmt.Print("USAGE:\n\n")
		fmt.Printf("  %s [flags] <EXPRESSION SYNTAX>...\n\n", name)

		fmt.Print("EXAMPLES:\n\n")
		fmt.Printf("  %s user = { name=foo age=:30 type=C }\n\n", name)
		fmt.Printf("  %s -f yaml kind=Service metadata.name=bb-entrypoint metadata.namespace=default\n\n", name)
		fmt.Print("FLAGS:\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.PrintDefaults()
		flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
		fmt.Print("  -help\n\tprints this message\n")
		fmt.Println()

		fmt.Println("Crafted with passion by Luca Sepe - https://github.com/lucasepe/map")
	}

	flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)

	flag.CommandLine.StringVar(&optFormat, "f", "json", "output format: json, yaml")
	flag.CommandLine.BoolVar(&optVersion, "v", false, "print current version and exit")

	flag.CommandLine.Parse(os.Args[1:])
}

func appName() string {
	return filepath.Base(os.Args[0])
}

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/lucasepe/clon/parser"
	"gopkg.in/yaml.v2"
)

const (
	banner = `     _         
 ___| |___ ___    commandline
|  _| | . |   | object notation
|___|_|___|_|_|    language`
)

var (
	optYAML    bool
	optVersion bool
	optEnvFile string
	gitCommit  string
)

func main() {
	configureFlags()

	if optVersion {
		fmt.Printf("%s version: %s\n", appName(), gitCommit)
		os.Exit(0)
	}

	godotenv.Load(optEnvFile)

	res, err := parseArgsOrStdIn()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	for _, g := range res {
		if optYAML {
			if err := toYAML(os.Stdout, g); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		} else {
			if err := toJSON(os.Stdout, g); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

func configureFlags() {
	name := appName()

	flag.CommandLine.Usage = func() {
		fmt.Printf("%s\n\n", banner)
		//fmt.Print("Commandline Object Notation.\n\n")

		fmt.Print("USAGE:\n\n")
		fmt.Printf("  %s [flags] <EXPRESSION SYNTAX>...\n\n", name)

		fmt.Print("EXAMPLES:\n\n")
		fmt.Printf("  %s user = { name=foo age=:30 type=C }\n\n", name)
		fmt.Printf("  %s -yaml kind=Service metadata.name=bb-entrypoint metadata.namespace=default\n\n", name)
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

	flag.CommandLine.StringVar(&optEnvFile, "env-file", ".env", "the dot env file")
	flag.CommandLine.BoolVar(&optYAML, "yaml", false, "output format YAML (default: JSON)")
	flag.CommandLine.BoolVar(&optVersion, "v", false, "print current version and exit")

	flag.CommandLine.Parse(os.Args[1:])
}

func appName() string {
	return filepath.Base(os.Args[0])
}

func parseArgsOrStdIn() ([]parser.Generator, error) {
	if len(flag.Args()) == 0 {
		return parser.ParseReader(os.Stdin)
	}

	return parser.ParseTextLines(flag.Args())
}

func toJSON(w io.Writer, g parser.Generator) error {
	dat, err := json.Marshal(g.Do())
	if err != nil {
		return err
	}

	var out bytes.Buffer
	if err := json.Indent(&out, dat, "", "   "); err != nil {
		return err
	}
	w.Write(out.Bytes())

	return nil
}

func toYAML(w io.Writer, g parser.Generator) error {
	dat, err := yaml.Marshal(g.Do())
	if err != nil {
		return err
	}

	w.Write(dat)

	return nil
}

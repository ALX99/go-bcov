package main

import (
	"encoding/xml"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/tools/cover"
)

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(in io.Reader, out io.Writer) error {
	profiles, err := cover.ParseProfilesFromReader(in)
	if err != nil {
		return err
	}

	lines, err := getCoverage(profiles)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, xml.Header)
	encoder := xml.NewEncoder(out)
	encoder.Indent("", "\t")
	err = encoder.Encode(makeCoverage(lines))
	if err != nil {
		return err
	}

	return encoder.Flush()
}

func getCoverage(profiles []*cover.Profile) (filesWithLines, error) {
	lines := make(map[string]map[int]Line)
	for _, profile := range profiles {
		newLines, file, err := getFileCoverage(profile)
		if err != nil {
			return nil, err
		}

		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		relPath, err := filepath.Rel(wd, file)
		if err != nil {
			return nil, err
		}

		lines[relPath] = newLines
	}
	return lines, nil
}

func getFileCoverage(profile *cover.Profile) (map[int]Line, string, error) {
	fileName := profile.FileName
	absFilePath, err := findAbsFile(fileName)
	if err != nil {
		return nil, "", err
	}
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, absFilePath, nil, 0)
	if err != nil {
		return nil, "", err
	}
	data, err := os.ReadFile(absFilePath)
	if err != nil {
		return nil, "", err
	}

	lines := make(map[int]Line)
	visitor := &fileVisitor{
		fset:     fSet{fset},
		fileData: data,
		profile:  profile,
		lines:    lines,
	}
	ast.Walk(visitor, parsed)
	return lines, absFilePath, nil
}

func findAbsFile(file string) (string, error) {
	dir, file := filepath.Split(file)
	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		return "", err
	}
	return filepath.Join(pkg.Dir, file), nil
}

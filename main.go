package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/cover"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer cancel()

	t := time.Now()
	fCount, err := run(ctx, os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Sonarqube coverage report generated in %s from %d files.\n",
		time.Since(t).Round(time.Millisecond), fCount)
}

func run(ctx context.Context, in io.Reader, out io.Writer) (int, error) {
	profiles, err := cover.ParseProfilesFromReader(in)
	if err != nil {
		return 0, err
	}

	lines, fCount, err := getCoverage(ctx, profiles)
	if err != nil {
		return 0, err
	}

	fmt.Fprintf(out, xml.Header)
	encoder := xml.NewEncoder(out)
	encoder.Indent("", "\t")
	err = encoder.Encode(makeCoverage(lines))
	if err != nil {
		return 0, err
	}

	return fCount, encoder.Flush()
}

func getCoverage(ctx context.Context, profiles []*cover.Profile) (filesWithLines, int, error) {
	lines := make(map[string]map[int]Line)
	fCount := 0
	mu := sync.Mutex{}
	eg, _ := errgroup.WithContext(ctx)
	eg.SetLimit(max(len(profiles), 50))

	wd, err := os.Getwd()
	if err != nil {
		return nil, 0, err
	}

	for _, profile := range profiles {
		fCount++
		eg.Go(func() error {
			newLines, file, err := getFileCoverage(profile)
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(wd, file)
			if err != nil {
				return err
			}

			mu.Lock()
			lines[relPath] = newLines
			mu.Unlock()

			return nil
		})
	}

	return lines, fCount, eg.Wait()
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

	// free memory
	clear(data)
	clear(visitor.fileData)
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

package main

import (
	"context"
	"encoding/xml"
	"flag"
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

var (
	supportedFormats = []string{"sonar-cover-report"}
	format           = flag.String("format", "reserved", "output format")
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer cancel()
	flag.Parse()

	if err := run(ctx, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, in io.Reader, out io.Writer) error {
	t := time.Now()
	profiles, err := cover.ParseProfilesFromReader(in)
	if err != nil {
		return err
	}

	report, fCount, err := getCoverage(ctx, profiles)
	if err != nil {
		return err
	}

	switch *format {
	case "sonar-cover-report":
		fmt.Fprintf(out, xml.Header)
		encoder := xml.NewEncoder(out)
		encoder.Indent("", "\t")
		err = report.ToSonarCoverage(encoder)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Sonarqube coverage report generated in %s from %d files.\n",
			time.Since(t).Round(time.Millisecond), fCount)
		return encoder.Flush()
	default:
		return fmt.Errorf("unknown format %q: must be one of %v", *format, supportedFormats)
	}
}

func getCoverage(ctx context.Context, profiles []*cover.Profile) (report, int, error) {
	files := make([]file, 0, len(profiles))
	fCount := 0
	mu := sync.Mutex{}
	eg, _ := errgroup.WithContext(ctx)
	eg.SetLimit(max(len(profiles), 50))

	wd, err := os.Getwd()
	if err != nil {
		return report{}, 0, err
	}

	for _, profile := range profiles {
		fCount++
		eg.Go(func() error {
			newLines, fileP, err := getFileCoverage(profile)
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(wd, fileP)
			if err != nil {
				return err
			}

			mu.Lock()
			files = append(files, file{path: relPath, lines: newLines.lines})
			mu.Unlock()

			return nil
		})
	}

	return report{files: files}, fCount, eg.Wait()
}

func getFileCoverage(profile *cover.Profile) (file, string, error) {
	fileName := profile.FileName
	absFilePath, err := findAbsFile(fileName)
	if err != nil {
		return file{}, "", err
	}
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, absFilePath, nil, 0)
	if err != nil {
		return file{}, "", err
	}
	data, err := os.ReadFile(absFilePath)
	if err != nil {
		return file{}, "", err
	}

	lines := make(map[int]Line)
	visitor := &fileVisitor{
		fset:     fSet{fset},
		fileData: data,
		profile:  profile,
		file:     file{lines: lines},
	}
	ast.Walk(visitor, parsed)

	// free memory
	clear(data)
	clear(visitor.fileData)
	return visitor.file, absFilePath, nil
}

func findAbsFile(file string) (string, error) {
	dir, file := filepath.Split(file)
	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		return "", err
	}
	return filepath.Join(pkg.Dir, file), nil
}

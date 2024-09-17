package main

import (
	"cmp"
	"encoding/xml"
	"log"
	"os"
	"reflect"
)

type line struct {
	LineNumber      int  `xml:"lineNumber,attr"`
	Covered         bool `xml:"covered,attr"`
	BranchesToCover *int `xml:"branchesToCover,attr,omitempty"`
	CoveredBranches *int `xml:"coveredBranches,attr,omitempty"`
}

type file struct {
	Path        string `xml:"path,attr"`
	LineToCover []line `xml:"lineToCover,omitempty"`
}

type coverage struct {
	XMLName xml.Name `xml:"coverage"`
	Version int      `xml:"version,attr"`
	File    []file   `xml:"file,omitempty"`
}

func main() {
	f1, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	f2, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	var coverage1, coverage2 coverage
	err = xml.NewDecoder(f1).Decode(&coverage1)
	if err != nil {
		log.Fatal(err)
	}

	err = xml.NewDecoder(f2).Decode(&coverage2)
	if err != nil {
		log.Fatal(err)
	}

	for _, file1 := range coverage1.File {
		fileFound := false
		for _, file2 := range coverage2.File {
			if file1.Path == file2.Path {
				fileFound = true
				for _, line1 := range file1.LineToCover {
					lineFound := false
					for _, line2 := range file2.LineToCover {
						if line1.LineNumber == line2.LineNumber {
							lineFound = true
							if !reflect.DeepEqual(line1, line2) {
								log.Printf("Difference found in file %s, line %d\n", file1.Path, line1.LineNumber)
								log.Printf("First coverage report: covered=%v, branchesToCover=%v, coveredBranches=%v\n",
									line1.Covered, *cmp.Or(line1.BranchesToCover, ptr(0)), *cmp.Or(line1.CoveredBranches, ptr(0)))
								log.Printf("Second coverage report: covered=%v, branchesToCover=%v, coveredBranches=%v\n",
									line2.Covered, *cmp.Or(line2.BranchesToCover, ptr(0)), *cmp.Or(line2.CoveredBranches, ptr(0)))
								os.Exit(1)
							}
						}
					}
					if !lineFound {
						log.Fatalf("Line %d not found in second coverage report\n", line1.LineNumber)
					}
				}
			}
		}
		if !fileFound {
			log.Fatalf("File %s not found in second coverage report\n", file1.Path)
		}
	}
	log.Printf("No differences found between %s and %s\n", os.Args[1], os.Args[2])
}

func ptr[Type any](v Type) *Type {
	return &v
}

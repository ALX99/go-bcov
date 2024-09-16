package main

import (
	"encoding/xml"
)

type filesWithLines map[string]map[int]Line

// file represents the "file" element
type file struct {
	Path        string `xml:"path,attr"`
	LineToCover []Line `xml:"lineToCover,omitempty"`
}

// coverage represents the root "coverage" element
type coverage struct {
	XMLName xml.Name `xml:"coverage"`
	Version int      `xml:"version,attr"`
	File    []file   `xml:"file,omitempty"`
}

type Line struct {
	Number          int  `xml:"lineNumber,attr"`
	Covered         bool `xml:"covered,attr"`
	BranchesToCover *int `xml:"branchesToCover,attr,omitempty"`
	CoveredBranches *int `xml:"coveredBranches,attr,omitempty"`
}

func makeCoverage(filesAndLines filesWithLines) coverage {
	newCov := coverage{Version: 1}
	for fileName, fileLines := range filesAndLines {
		newFile := file{Path: fileName}
		for _, line := range fileLines {
			newFile.LineToCover = append(newFile.LineToCover, line)
		}
		newCov.File = append(newCov.File, newFile)
	}
	return newCov
}

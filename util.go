package main

import (
	"encoding/xml"
)

type filesWithLines map[string]map[int]Line

type Line struct {
	IsSingleIf      bool
	CoveredCount    int
	BranchesToCover *int
	CoveredBranches *int
	IfBodyStartLine int
}

type report struct {
	files []file
}

type file struct {
	path  string
	lines map[int]Line
}

func (r report) toSonarCoverage(e *xml.Encoder) error {
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

	newCov := coverage{Version: 1}
	for _, rFile := range r.files {
		newFile := file{Path: rFile.path}
		for index, lineInfo := range rFile.lines {
			newFile.LineToCover = append(newFile.LineToCover, line{
				LineNumber:      index,
				Covered:         lineInfo.CoveredCount > 0,
				BranchesToCover: lineInfo.BranchesToCover,
				CoveredBranches: lineInfo.CoveredBranches,
			})
		}
		newCov.File = append(newCov.File, newFile)
	}
	return e.Encode(newCov)
}

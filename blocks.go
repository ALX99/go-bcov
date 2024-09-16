package main

import "golang.org/x/tools/cover"

type blocks []cover.ProfileBlock

func (blocks blocks) allLinesCovered(startLine, endLine, startCol, endCol int) bool {
	for _, b := range blocks {
		if b.StartLine > endLine || (b.StartLine == endLine && b.StartCol >= endCol) {
			break
		}
		if b.EndLine < startLine || (b.EndLine == startLine && b.EndCol <= startCol) {
			continue
		}
		return b.Count > 0
	}
	return false
}

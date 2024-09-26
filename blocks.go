package main

import (
	"go/token"

	"golang.org/x/tools/cover"
)

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

func (blocks blocks) getCoveredCount(pos token.Position) int {
	for _, b := range blocks {
		if b.StartLine > pos.Line || (b.StartLine == pos.Line && b.StartCol >= pos.Column) {
			break
		}
		if b.EndLine < pos.Line || (b.EndLine == pos.Line && b.EndCol <= pos.Column) {
			continue
		}
		return b.Count
	}
	return 0
}

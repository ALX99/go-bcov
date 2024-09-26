package main

import (
	"cmp"
	"go/ast"

	"golang.org/x/tools/cover"
)

type fileVisitor struct {
	fset     fSet
	fileData []byte
	profile  *cover.Profile
	file     file
}

func (v *fileVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case ast.Stmt:
		start := v.fset.getPos(n.Pos())
		end := v.fset.getPos(n.End())
		startLine := start.Line
		startCol := start.Column
		endLine := end.Line
		endCol := end.Column

		for _, b := range v.profile.Blocks {
			if b.StartLine > endLine || (b.StartLine == endLine && b.StartCol >= endCol) {
				break
			}
			if b.EndLine < startLine || (b.EndLine == startLine && b.EndCol <= startCol) {
				continue
			}
			for i := b.StartLine; i <= b.EndLine; i++ {
				line, ok := v.file.lines[i]
				if !ok {
					line = Line{}
				}
				line.CoveredCount = max(b.Count, line.CoveredCount)
				v.file.lines[i] = line
			}
		}
		switch n := n.(type) {
		case *ast.IfStmt:
			startLine = v.fset.getPos(n.If).Line

			line, ok := v.file.lines[startLine]
			if !ok {
				line = Line{}
			}
			// nil protection
			line.BranchesToCover = cmp.Or(line.BranchesToCover, ptr(0))
			line.CoveredBranches = cmp.Or(line.CoveredBranches, ptr(0))
			*line.CoveredBranches = v.fset.checkIfBranchCovered(n, v.profile.Blocks)
			*line.BranchesToCover = 2

			v.file.lines[startLine] = line
		case *ast.SwitchStmt:
			startLine = v.fset.getPos(n.Switch).Line

			line, ok := v.file.lines[startLine]
			if !ok {
				line = Line{}
			}
			// nil protection
			line.BranchesToCover = cmp.Or(line.BranchesToCover, ptr(0))
			line.CoveredBranches = cmp.Or(line.CoveredBranches, ptr(0))

			branches, covered := v.fset.checkSwitchBranchCovered(n, v.profile.Blocks)

			*line.BranchesToCover += branches
			*line.CoveredBranches += covered

			v.file.lines[startLine] = line

		}
	}
	return v
}

func ptr[Type any](v Type) *Type {
	return &v
}

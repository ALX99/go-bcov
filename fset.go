package main

import (
	"go/ast"
	"go/token"
)

type fSet struct {
	*token.FileSet
}

func (f fSet) getPos(pos token.Pos) token.Position {
	if !pos.IsValid() {
		panic("Node position is not valid, is your AST correct?")
	}
	return f.Position(pos)
}

func (f fSet) checkIfBranchCovered(ifStmt *ast.IfStmt, blocks blocks) int {
	covered := 0

	stmtscovered := true
	for _, stmt := range ifStmt.Body.List {
		start := f.getPos(stmt.Pos())
		end := f.getPos(stmt.End())
		if !blocks.allLinesCovered(start.Line, end.Line, start.Column, end.Column) {
			stmtscovered = false
			break
		}
	}

	if stmtscovered {
		covered++
	}

	if len(ifStmt.Body.List) > 0 &&
		blocks.getCoveredCount(f.getPos(ifStmt.Cond.Pos())) > blocks.getCoveredCount(f.getPos(ifStmt.Body.List[0].Pos())) {
		covered++
	}

	return covered
}

func (f fSet) checkCaseCoverage(caseStmt *ast.CaseClause, magicMap map[int]int, blocks blocks) int {
	covered := 0

	stmtscovered := true
	for _, stmt := range caseStmt.Body {
		start := f.getPos(stmt.Pos())
		end := f.getPos(stmt.End())
		if !blocks.allLinesCovered(start.Line, end.Line, start.Column, end.Column) {
			stmtscovered = false
			break
		}
	}

	if stmtscovered {
		covered++
	}

	if len(caseStmt.Body) > 0 &&
		magicMap[f.getPos(caseStmt.Case).Line] > blocks.getCoveredCount(f.getPos(caseStmt.Body[0].Pos())) {
		covered++
	}

	return covered
}

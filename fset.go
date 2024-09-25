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

func (f fSet) checkIfBranchCovered(ifStmt *ast.IfStmt, blocks blocks) (branches, covered int) {
	branches = 2

	// Check the body of the if statement
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

	if ifStmt.Else == nil {
		return
	}

	// Check the body of the else statement
	switch elseStmt := ifStmt.Else.(type) {
	case *ast.BlockStmt:
		for _, stmt := range elseStmt.List {
			start := f.getPos(stmt.Pos())
			end := f.getPos(stmt.End())
			if blocks.allLinesCovered(start.Line, end.Line, start.Column, end.Column) {
				covered++
				break
			}
		}
	case *ast.IfStmt:
		// Handle else-if case
		branchesElse, coveredElse := f.checkIfBranchCovered(elseStmt, blocks)
		branches += branchesElse - 1 // subtract 1 because the else-if adds an extra branch
		covered += coveredElse
	}

	return
}

func (f fSet) checkSwitchBranchCovered(switchStmt *ast.SwitchStmt, blocks blocks) (branches, covered int) {
	for _, stmt := range switchStmt.Body.List {
		caseClause, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}
		branches++

		stmtscovered := true
		for _, stmt := range caseClause.Body {
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
	}
	return
}

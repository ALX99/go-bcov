package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/cover"
)

func Test_fSet_checkSwitchBranchCovered(t *testing.T) {
	fset := token.NewFileSet()
	dummySrc := `package main

import "fmt"

func main() {
	i,_ := fmt.Scanf("%d", nil)
	switch i {
	case 1:
		fmt.Println("1")
	case 2:
		fmt.Println("2")
	case 3:
		fmt.Println("3")
		fmt.Println("3")
	default:
		fmt.Println("default")
	}
}`

	file, err := parser.ParseFile(fset, "main.go", dummySrc, 0)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	switchStmt := file.Decls[1].(*ast.FuncDecl).Body.List[1].(*ast.SwitchStmt)

	type fields struct {
		FileSet *token.FileSet
	}
	type args struct {
		switchStmt *ast.SwitchStmt
		blocks     blocks
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantBranches int
		wantCovered  int
	}{
		{
			name: "Default branch covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				switchStmt: switchStmt,
				blocks: []cover.ProfileBlock{
					{StartLine: 16, StartCol: 3, EndLine: 16, EndCol: 25, Count: 1},
				},
			},
			wantBranches: 4,
			wantCovered:  1,
		},
		{
			name: "Case 1 covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				switchStmt: switchStmt,
				blocks: []cover.ProfileBlock{
					{StartLine: 9, StartCol: 3, EndLine: 9, EndCol: 21, Count: 1},
				},
			},
			wantBranches: 4,
			wantCovered:  1,
		},
		{
			name: "Case 2 covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				switchStmt: switchStmt,
				blocks: []cover.ProfileBlock{
					{StartLine: 11, StartCol: 3, EndLine: 11, EndCol: 21, Count: 1},
				},
			},
			wantBranches: 4,
			wantCovered:  1,
		},
		{
			name: "Case 3 covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				switchStmt: switchStmt,
				blocks: []cover.ProfileBlock{
					{StartLine: 13, StartCol: 3, EndLine: 14, EndCol: 21, Count: 1},
				},
			},
			wantBranches: 4,
			wantCovered:  1,
		},
		{
			name: "Case 3 partially covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				switchStmt: switchStmt,
				blocks: []cover.ProfileBlock{
					{StartLine: 13, StartCol: 3, EndLine: 13, EndCol: 21, Count: 1},
				},
			},
			wantBranches: 4,
			wantCovered:  0,
		},
		{
			name: "All cases covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				switchStmt: switchStmt,
				blocks: []cover.ProfileBlock{
					{StartLine: 9, StartCol: 3, EndLine: 9, EndCol: 21, Count: 1},
					{StartLine: 11, StartCol: 3, EndLine: 11, EndCol: 21, Count: 1},
					{StartLine: 13, StartCol: 3, EndLine: 14, EndCol: 21, Count: 1},
					{StartLine: 16, StartCol: 3, EndLine: 16, EndCol: 25, Count: 1},
				},
			},
			wantBranches: 4,
			wantCovered:  4,
		},
		{
			name: "No cases covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				switchStmt: switchStmt,
				blocks:     []cover.ProfileBlock{},
			},
			wantBranches: 4,
			wantCovered:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fSet{
				FileSet: tt.fields.FileSet,
			}
			gotBranches, gotCovered := f.checkSwitchBranchCovered(tt.args.switchStmt, tt.args.blocks)
			if gotBranches != tt.wantBranches {
				t.Errorf("fSet.checkSwitchBranchCovered() gotBranches = %v, want %v", gotBranches, tt.wantBranches)
			}
			if gotCovered != tt.wantCovered {
				t.Errorf("fSet.checkSwitchBranchCovered() gotCovered = %v, want %v", gotCovered, tt.wantCovered)
			}
		})
	}
}

func Test_fSet_checkIfBranchCovered(t *testing.T) {
	fset := token.NewFileSet()
	dummySrc := `package main

import "fmt"

func main() {
	_, err := fmt.Scanf("%d", nil)
	if err != nil {
		println("true")
		println("true")
	} else if err == nil {
		println("else if")
	} else {
		println("false")
	}
}`

	file, err := parser.ParseFile(fset, "main.go", dummySrc, 0)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	ifStmt := file.Decls[1].(*ast.FuncDecl).Body.List[1].(*ast.IfStmt)

	type fields struct {
		FileSet *token.FileSet
	}
	type args struct {
		ifStmt *ast.IfStmt
		blocks blocks
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantBranches int
		wantCovered  int
	}{
		{
			name: "If branch covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				ifStmt: ifStmt,
				blocks: []cover.ProfileBlock{
					{StartLine: 8, StartCol: 3, EndLine: 9, EndCol: 18, Count: 1},
				},
			},
			wantBranches: 3,
			wantCovered:  1,
		},
		{
			name: "Else if branch covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				ifStmt: ifStmt,
				blocks: []cover.ProfileBlock{
					{StartLine: 11, StartCol: 3, EndLine: 11, EndCol: 21, Count: 1},
				},
			},
			wantBranches: 3,
			wantCovered:  1,
		},
		{
			name: "Else branch covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				ifStmt: ifStmt,
				blocks: []cover.ProfileBlock{
					{StartLine: 13, StartCol: 3, EndLine: 13, EndCol: 20, Count: 1},
				},
			},
			wantBranches: 3,
			wantCovered:  1,
		},
		{
			name: "Both if and else if branches covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				ifStmt: ifStmt,
				blocks: []cover.ProfileBlock{
					{StartLine: 8, StartCol: 3, EndLine: 9, EndCol: 18, Count: 1},
					{StartLine: 11, StartCol: 3, EndLine: 11, EndCol: 21, Count: 1},
				},
			},
			wantBranches: 3,
			wantCovered:  2,
		},
		{
			name: "All branches covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				ifStmt: ifStmt,
				blocks: []cover.ProfileBlock{
					{StartLine: 8, StartCol: 3, EndLine: 9, EndCol: 18, Count: 1},
					{StartLine: 11, StartCol: 3, EndLine: 11, EndCol: 21, Count: 1},
					{StartLine: 13, StartCol: 3, EndLine: 13, EndCol: 20, Count: 1},
				},
			},
			wantBranches: 3,
			wantCovered:  3,
		},
		{
			name: "No branches covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				ifStmt: ifStmt,
				blocks: []cover.ProfileBlock{},
			},
			wantBranches: 3,
			wantCovered:  0,
		},
		{
			name: "If branch partially covered",
			fields: fields{
				FileSet: fset,
			},
			args: args{
				ifStmt: ifStmt,
				blocks: []cover.ProfileBlock{
					{StartLine: 8, StartCol: 3, EndLine: 8, EndCol: 18, Count: 1},
				},
			},
			wantBranches: 3,
			wantCovered:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fSet{
				FileSet: tt.fields.FileSet,
			}
			gotBranches, gotCovered := f.checkIfBranchCovered(tt.args.ifStmt, tt.args.blocks)
			if gotBranches != tt.wantBranches {
				t.Errorf("fSet.checkIfBranchCovered() gotBranches = %v, want %v", gotBranches, tt.wantBranches)
			}
			if gotCovered != tt.wantCovered {
				t.Errorf("fSet.checkIfBranchCovered() gotCovered = %v, want %v", gotCovered, tt.wantCovered)
			}
		})
	}
}

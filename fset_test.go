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

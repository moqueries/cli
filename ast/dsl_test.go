package ast_test

import (
	"go/token"
	"reflect"
	"testing"

	"github.com/dave/dst"

	"github.com/myshkin5/moqueries/ast"
)

var (
	id1 = dst.NewIdent("id1")
	id2 = dst.NewIdent("id2")
	id3 = dst.NewIdent("id3")
	id4 = dst.NewIdent("id4")
	id5 = dst.NewIdent("id5")
	id6 = dst.NewIdent("id6")

	field1 = &dst.Field{Type: id1}
	field2 = &dst.Field{Type: id2}

	params  = &dst.FieldList{List: []*dst.Field{{Type: id3}, {Type: id4}}}
	results = &dst.FieldList{List: []*dst.Field{{Type: id5}, {Type: id6}}}

	assign1 = &dst.AssignStmt{Lhs: []dst.Expr{id1}}
	assign2 = &dst.AssignStmt{Lhs: []dst.Expr{id2}}
	assign3 = &dst.AssignStmt{Lhs: []dst.Expr{id3}}
	assign4 = &dst.AssignStmt{Lhs: []dst.Expr{id4}}
)

func TestAssign(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.AssignStmt{}

		// ACT
		actual := ast.Assign().Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		decs := dst.AssignStmtDecorations{
			NodeDecs: dst.NodeDecs{Before: dst.EmptyLine},
		}
		expected := &dst.AssignStmt{
			Lhs:  []dst.Expr{id1, id2},
			Tok:  token.ASSIGN,
			Rhs:  []dst.Expr{id3, id4},
			Decs: decs,
		}

		// ACT
		actual := ast.Assign(id1, id2).Tok(token.ASSIGN).Rhs(id3, id4).Decs(decs).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestAssignDecs(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := dst.AssignStmtDecorations{
			NodeDecs: dst.NodeDecs{Before: dst.EmptyLine},
		}

		// ACT
		actual := ast.AssignDecs(dst.EmptyLine).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := dst.AssignStmtDecorations{
			NodeDecs: dst.NodeDecs{
				Before: dst.EmptyLine,
				After:  dst.NewLine,
			},
		}

		// ACT
		actual := ast.AssignDecs(dst.EmptyLine).After(dst.NewLine).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestBin(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.BinaryExpr{X: id1}

		// ACT
		actual := ast.Bin(id1).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.BinaryExpr{X: id1, Op: token.EQL, Y: id2}

		// ACT
		actual := ast.Bin(id1).Op(token.EQL).Y(id2).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestBlock(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.BlockStmt{}

		// ACT
		actual := ast.Block().Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.BlockStmt{List: []dst.Stmt{assign1, assign2}}

		// ACT
		actual := ast.Block(assign1, assign2).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestBreak(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.BranchStmt{Tok: token.BREAK}

		// ACT
		actual := ast.Break()

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestCall(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.CallExpr{Fun: id1}

		// ACT
		actual := ast.Call(id1).Args().Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.CallExpr{
			Fun:      id1,
			Args:     []dst.Expr{id2, id3},
			Ellipsis: true,
		}

		// ACT
		actual := ast.Call(id1).Args(id2, id3).Ellipsis(true).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestComp(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.CompositeLit{Type: id1}

		// ACT
		actual := ast.Comp(id1).Elts().Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		decs := dst.CompositeLitDecorations{Lbrace: []string{"ciao"}}
		expected := &dst.CompositeLit{
			Type: id1,
			Elts: []dst.Expr{id2, id3},
			Decs: decs,
		}

		// ACT
		actual := ast.Comp(id1).Elts(id2, id3).Decs(decs).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestContinue(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.BranchStmt{Tok: token.CONTINUE}

		// ACT
		actual := ast.Continue()

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestEllipsis(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.Ellipsis{Elt: id1}

		// ACT
		actual := ast.Ellipsis(id1)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestExpr(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.ExprStmt{X: id1}

		// ACT
		actual := ast.Expr(id1).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		decs := dst.ExprStmtDecorations{
			NodeDecs: dst.NodeDecs{After: dst.EmptyLine},
		}
		expected := &dst.ExprStmt{
			X:    id1,
			Decs: decs,
		}

		// ACT
		actual := ast.Expr(id1).Decs(decs).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestExprDecs(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := dst.ExprStmtDecorations{
			NodeDecs: dst.NodeDecs{After: dst.EmptyLine},
		}

		// ACT
		actual := ast.ExprDecs(dst.EmptyLine).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestField(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.Field{Type: id1}

		// ACT
		actual := ast.Field(id1).Names().Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		decs := dst.FieldDecorations{
			NodeDecs: dst.NodeDecs{Before: dst.EmptyLine, After: dst.NewLine},
		}
		expected := &dst.Field{
			Type:  id1,
			Names: []*dst.Ident{id2, id3},
			Decs:  decs,
		}

		// ACT
		actual := ast.Field(id1).Names(id2, id3).Decs(decs).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestFieldDecs(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := dst.FieldDecorations{
			NodeDecs: dst.NodeDecs{Before: dst.EmptyLine, After: dst.NewLine},
		}

		// ACT
		actual := ast.FieldDecs(dst.EmptyLine, dst.NewLine).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestFieldList(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.FieldList{}

		// ACT
		actual := ast.FieldList()

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.FieldList{List: []*dst.Field{field1, field2}}

		// ACT
		actual := ast.FieldList(field1, field2)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestFn(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.FuncDecl{
			Name: dst.NewIdent("fn1"),
			Type: &dst.FuncType{
				Params:  &dst.FieldList{},
				Results: &dst.FieldList{},
			},
			Body: &dst.BlockStmt{},
		}

		// ACT
		actual := ast.Fn("fn1").Params().Results().Body().Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		decs := dst.FuncDeclDecorations{NodeDecs: dst.NodeDecs{Before: dst.EmptyLine}}
		expected := &dst.FuncDecl{
			Name: dst.NewIdent("fn1"),
			Type: &dst.FuncType{Params: params, Results: results},
			Recv: &dst.FieldList{List: []*dst.Field{{Type: id1}, {Type: id2}}},
			Body: &dst.BlockStmt{List: []dst.Stmt{assign1, assign2}},
			Decs: decs,
		}

		// ACT
		actual := ast.Fn("fn1").
			Recv(ast.Field(id1).Obj, ast.Field(id2).Obj).
			Params(ast.Field(id3).Obj, ast.Field(id4).Obj).
			Results(ast.Field(id5).Obj, ast.Field(id6).Obj).
			Body(assign1, assign2).
			Decs(decs).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete - field lists", func(t *testing.T) {
		// ASSEMBLE
		decs := dst.FuncDeclDecorations{NodeDecs: dst.NodeDecs{Before: dst.EmptyLine}}
		expected := &dst.FuncDecl{
			Name: dst.NewIdent("fn1"),
			Type: &dst.FuncType{Params: params, Results: results},
			Recv: &dst.FieldList{List: []*dst.Field{{Type: id1}, {Type: id2}}},
			Body: &dst.BlockStmt{List: []dst.Stmt{assign1, assign2}},
			Decs: decs,
		}

		// ACT
		actual := ast.Fn("fn1").
			Recv(ast.Field(id1).Obj, ast.Field(id2).Obj).
			ParamList(params).
			ResultList(results).
			Body(assign1, assign2).
			Decs(decs).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestFnLit(t *testing.T) {
	fnType := ast.FnType(params).Obj

	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.FuncLit{Type: fnType}

		// ACT
		actual := ast.FnLit(fnType).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.FuncLit{
			Type: fnType,
			Body: &dst.BlockStmt{List: []dst.Stmt{assign1, assign2}},
		}

		// ACT
		actual := ast.FnLit(fnType).Body(assign1, assign2).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestFnType(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.FuncType{Params: params}

		// ACT
		actual := ast.FnType(params).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.FuncType{Params: params, Results: results}

		// ACT
		actual := ast.FnType(params).Results(results).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestFor(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.ForStmt{Init: assign1, Body: &dst.BlockStmt{}}

		// ACT
		actual := ast.For(assign1).Body().Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.ForStmt{
			Init: assign1,
			Cond: id1,
			Post: assign2,
			Body: &dst.BlockStmt{List: []dst.Stmt{assign3, assign4}},
		}

		// ACT
		actual := ast.For(assign1).Cond(id1).Post(assign2).Body(assign3, assign4).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestId(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := dst.NewIdent("hola")

		// ACT
		actual := ast.Id("hola")

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestIdf(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := dst.NewIdent("hola 1 two")

		// ACT
		actual := ast.Idf("hola %d %s", 1, "two")

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestIdPath(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.Ident{Name: "str1", Path: "str2"}

		// ACT
		actual := ast.IdPath("str1", "str2")

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestIncStmt(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.IncDecStmt{X: id1, Tok: token.INC}

		// ACT
		actual := ast.IncStmt(id1)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestIndex(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.IndexExpr{X: id1}

		// ACT
		actual := ast.Index(id1).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.IndexExpr{X: id1, Index: id2}

		// ACT
		actual := ast.Index(id1).Sub(id2).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestIf(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.IfStmt{Cond: id1, Body: &dst.BlockStmt{}}

		// ACT
		actual := ast.If(id1).Body().Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		decs := dst.IfStmtDecorations{
			NodeDecs: dst.NodeDecs{After: dst.EmptyLine},
		}
		expected := &dst.IfStmt{
			Cond: id1,
			Body: &dst.BlockStmt{List: []dst.Stmt{assign1, assign2}},
			Else: assign3,
			Decs: decs,
		}

		// ACT
		actual := ast.If(id1).Body(assign1, assign2).Else(assign3).Decs(decs).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("block else", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.IfStmt{
			Cond: id1,
			Else: &dst.BlockStmt{List: []dst.Stmt{assign3, assign4}},
		}

		// ACT
		actual := ast.If(id1).Else(assign3, assign4).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestIfDecs(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := dst.IfStmtDecorations{
			NodeDecs: dst.NodeDecs{After: dst.EmptyLine},
		}

		// ACT
		actual := ast.IfDecs(dst.EmptyLine).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestKeyValue(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.KeyValueExpr{Key: id1}

		// ACT
		actual := ast.Key(id1).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		decs := dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{Before: dst.EmptyLine},
		}
		expected := &dst.KeyValueExpr{Key: id1, Value: id2, Decs: decs}

		// ACT
		actual := ast.Key(id1).Value(id2).Decs(decs).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestKeyValueDecs(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{Before: dst.EmptyLine},
		}

		// ACT
		actual := ast.KeyValueDecs(dst.EmptyLine).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{
				Before: dst.EmptyLine,
				After:  dst.NewLine,
			},
		}

		// ACT
		actual := ast.KeyValueDecs(dst.EmptyLine).After(dst.NewLine).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestLitInt(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.BasicLit{Kind: token.INT, Value: "42"}

		// ACT
		actual := ast.LitInt(42)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestLitString(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.BasicLit{Kind: token.STRING, Value: "\"bonjour\""}

		// ACT
		actual := ast.LitString("bonjour")

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestLitStringf(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.BasicLit{Kind: token.STRING, Value: "\"bonjour 1 two\""}

		// ACT
		actual := ast.LitStringf("bonjour %d %s", 1, "two")

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestMapType(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.MapType{Key: id1}

		// ACT
		actual := ast.MapType(id1).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.MapType{Key: id1, Value: id2}

		// ACT
		actual := ast.MapType(id1).Value(id2).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestNodeDecsf(t *testing.T) {
	t.Run("short", func(t *testing.T) {
		// ASSEMBLE
		expected := dst.NodeDecs{
			Before: dst.NewLine,
			Start: dst.Decorations{
				"// here's a short line",
			},
		}

		// ACT
		actual := ast.NodeDecsf("// here's a %s line", "short")

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("long", func(t *testing.T) {
		// ASSEMBLE
		expected := dst.NodeDecs{
			Before: dst.NewLine,
			Start: dst.Decorations{
				"// here's a long line that seems to go on forever and never quite finishes in a",
				"// single line",
			},
		}

		// ACT
		actual := ast.NodeDecsf("// here's a %s line",
			"long line that seems to go on forever and never quite finishes in a single")

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("nothing", func(t *testing.T) {
		// ASSEMBLE
		expected := dst.NodeDecs{
			Before: dst.NewLine,
			Start:  dst.Decorations{},
		}

		// ACT
		actual := ast.NodeDecsf("")

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestParen(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.ParenExpr{X: id1}

		// ACT
		actual := ast.Paren(id1)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestRange(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.RangeStmt{X: id1, Body: &dst.BlockStmt{}}

		// ACT
		actual := ast.Range(id1).Body().Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.RangeStmt{
			X:     id1,
			Key:   id2,
			Value: id3,
			Tok:   token.DEFINE,
			Body:  &dst.BlockStmt{List: []dst.Stmt{assign1, assign2}},
		}

		// ACT
		actual := ast.Range(id1).
			Key(id2).
			Value(id3).
			Tok(token.DEFINE).
			Body(assign1, assign2).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestReturn(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.ReturnStmt{}

		// ACT
		actual := ast.Return()

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.ReturnStmt{Results: []dst.Expr{id1, id2}}

		// ACT
		actual := ast.Return(id1, id2)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestSel(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.SelectorExpr{X: id1}

		// ACT
		actual := ast.Sel(id1).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.SelectorExpr{X: id1, Sel: id2}

		// ACT
		actual := ast.Sel(id1).Dot(id2).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestSliceExpr(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.SliceExpr{X: id1}

		// ACT
		actual := ast.SliceExpr(id1).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.SliceExpr{X: id1, Low: id2, High: id3}

		// ACT
		actual := ast.SliceExpr(id1).Low(id2).High(id3).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestSliceType(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.ArrayType{Elt: id1}

		// ACT
		actual := ast.SliceType(id1)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestStar(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.StarExpr{X: id1}

		// ACT
		actual := ast.Star(id1)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestStruct(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.StructType{
			Fields: &dst.FieldList{Opening: true, Closing: true},
		}

		// ACT
		actual := ast.Struct()

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.StructType{
			Fields: &dst.FieldList{List: []*dst.Field{field1, field2}},
		}

		// ACT
		actual := ast.Struct(field1, field2)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestStructFromList(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.StructType{
			Fields: &dst.FieldList{Opening: true, Closing: true},
		}

		// ACT
		actual := ast.StructFromList(nil)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("empty", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.StructType{
			Fields: &dst.FieldList{Opening: true, Closing: true},
		}
		emptyList := &dst.FieldList{}

		// ACT
		actual := ast.StructFromList(emptyList)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.StructType{Fields: params}

		// ACT
		actual := ast.StructFromList(params)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestTypeSpec(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.TypeSpec{Name: dst.NewIdent("typ")}

		// ACT
		actual := ast.TypeSpec("typ").Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.TypeSpec{Name: dst.NewIdent("typ"), Type: id2}

		// ACT
		actual := ast.TypeSpec("typ").Type(id2).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

//nolint:dupl // separate tests are preferred here
func TestTypeDecl(t *testing.T) {
	typ := &dst.TypeSpec{Name: id1}

	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.GenDecl{Tok: token.TYPE, Specs: []dst.Spec{typ}}

		// ACT
		actual := ast.TypeDecl(typ).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		decs := dst.GenDeclDecorations{
			NodeDecs: dst.NodeDecs{Before: dst.EmptyLine},
		}
		expected := &dst.GenDecl{
			Tok:   token.TYPE,
			Specs: []dst.Spec{typ},
			Decs:  decs,
		}

		// ACT
		actual := ast.TypeDecl(typ).Decs(decs).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestUn(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.UnaryExpr{Op: token.AND, X: id1}

		// ACT
		actual := ast.Un(token.AND, id1)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestValue(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.ValueSpec{Type: id1}

		// ACT
		actual := ast.Value(id1).Names().Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.ValueSpec{
			Type:   id1,
			Names:  []*dst.Ident{id2, id3},
			Values: []dst.Expr{id4, id5},
		}

		// ACT
		actual := ast.Value(id1).Names(id2, id3).Values(id4, id5).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

func TestVar(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.DeclStmt{Decl: &dst.GenDecl{Tok: token.VAR}}

		// ACT
		actual := ast.Var()

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		typ1 := &dst.TypeSpec{Name: id1}
		typ2 := &dst.TypeSpec{Name: id2}
		expected := &dst.DeclStmt{Decl: &dst.GenDecl{
			Tok:   token.VAR,
			Specs: []dst.Spec{typ1, typ2},
		}}

		// ACT
		actual := ast.Var(typ1, typ2)

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

//nolint:dupl // separate tests are preferred here
func TestVarDecl(t *testing.T) {
	typ := &dst.TypeSpec{Name: id1}

	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		expected := &dst.GenDecl{Tok: token.VAR, Specs: []dst.Spec{typ}}

		// ACT
		actual := ast.VarDecl(typ).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})

	t.Run("complete", func(t *testing.T) {
		// ASSEMBLE
		decs := dst.GenDeclDecorations{
			NodeDecs: dst.NodeDecs{Before: dst.EmptyLine},
		}
		expected := &dst.GenDecl{
			Tok:   token.VAR,
			Specs: []dst.Spec{typ},
			Decs:  decs,
		}

		// ACT
		actual := ast.VarDecl(typ).Decs(decs).Obj

		// ASSERT
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %#v, want %#v", actual, expected)
		}
	})
}

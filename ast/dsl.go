package ast

import (
	"fmt"
	"go/token"
	"strconv"
	"strings"

	"github.com/dave/dst"
)

// AssignDSL translates to a dst.AssignStmt
type AssignDSL struct{ Obj *dst.AssignStmt }

// Assign creates a new AssignDSL
func Assign(lhs ...dst.Expr) AssignDSL {
	return AssignDSL{Obj: &dst.AssignStmt{Lhs: lhs}}
}

// Tok specifies the token used in an assignment
func (d AssignDSL) Tok(tok token.Token) AssignDSL {
	d.Obj.Tok = tok
	return d
}

// Rhs specifies the right-hand expressions in the assignment
func (d AssignDSL) Rhs(rhs ...dst.Expr) AssignDSL {
	d.Obj.Rhs = append(d.Obj.Rhs, rhs...)
	return d
}

// Decs adds decorations to an AssignDSL
func (d AssignDSL) Decs(decs dst.AssignStmtDecorations) AssignDSL {
	d.Obj.Decs = decs
	return d
}

// AssignDecsDSL translates to a dst.AssignStmtDecorations
type AssignDecsDSL struct{ Obj dst.AssignStmtDecorations }

// AssignDecs creates a new AssignDecsDSL
func AssignDecs(before dst.SpaceType) AssignDecsDSL {
	return AssignDecsDSL{Obj: dst.AssignStmtDecorations{
		NodeDecs: dst.NodeDecs{Before: before},
	}}
}

// After adds whitespace after an assign statement
func (d AssignDecsDSL) After(after dst.SpaceType) AssignDecsDSL {
	d.Obj.After = after
	return d
}

// BinDSL translates to a dst.BinaryExpr
type BinDSL struct{ Obj *dst.BinaryExpr }

// Bin creates a new BinDSL
func Bin(x dst.Expr) BinDSL {
	return BinDSL{Obj: &dst.BinaryExpr{X: x}}
}

// Op specifies the operator
func (d BinDSL) Op(op token.Token) BinDSL {
	d.Obj.Op = op
	return d
}

// Y specifies the right side expression
func (d BinDSL) Y(y dst.Expr) BinDSL {
	d.Obj.Y = y
	return d
}

// BlockDSL translates to a dst.BlockStmt
type BlockDSL struct{ Obj *dst.BlockStmt }

// Block creates a new BlockDSL for a list of statements
func Block(list ...dst.Stmt) BlockDSL {
	return BlockDSL{Obj: &dst.BlockStmt{List: list}}
}

// Break creates a break dst.BranchStmt
func Break() *dst.BranchStmt {
	return &dst.BranchStmt{Tok: token.BREAK}
}

// CallDSL translates to a dst.CallExpr
type CallDSL struct{ Obj *dst.CallExpr }

// Call creates a new CallDSL
func Call(fun dst.Expr) CallDSL {
	return CallDSL{Obj: &dst.CallExpr{Fun: fun}}
}

// Args specifies the arguments to a call
func (d CallDSL) Args(args ...dst.Expr) CallDSL {
	d.Obj.Args = args
	return d
}

// Ellipsis specifies if the last argument is variadic
func (d CallDSL) Ellipsis(ellipsis bool) CallDSL {
	d.Obj.Ellipsis = ellipsis
	return d
}

// CompDSL translates to a dst.CompositeLit
type CompDSL struct{ Obj *dst.CompositeLit }

// Comp creates a new CompDSL
func Comp(typ dst.Expr) CompDSL {
	return CompDSL{Obj: &dst.CompositeLit{Type: typ}}
}

// Elts defines the elements of a CompDSL
func (d CompDSL) Elts(elts ...dst.Expr) CompDSL {
	d.Obj.Elts = elts
	return d
}

// Decs adds decorations to a CompDSL
func (d CompDSL) Decs(decs dst.CompositeLitDecorations) CompDSL {
	d.Obj.Decs = decs
	return d
}

// Continue creates a continue dst.BranchStmt
func Continue() *dst.BranchStmt {
	return &dst.BranchStmt{Tok: token.CONTINUE}
}

// Ellipsis translates an expression to a dst.Ellipsis
func Ellipsis(elt dst.Expr) *dst.Ellipsis {
	return &dst.Ellipsis{Elt: elt}
}

// ExprDSL translates to a dst.ExprStmt
type ExprDSL struct{ Obj *dst.ExprStmt }

// Expr returns a new ExprDSL
func Expr(x dst.Expr) ExprDSL {
	return ExprDSL{Obj: &dst.ExprStmt{X: x}}
}

// Decs adds decorations to a ExprDSL
func (d ExprDSL) Decs(decs dst.ExprStmtDecorations) ExprDSL {
	d.Obj.Decs = decs
	return d
}

// ExprDecsDSL translates to a dst.ExprStmtDecorations
type ExprDecsDSL struct{ Obj dst.ExprStmtDecorations }

// ExprDecs creates a new ExprDecsDSL
func ExprDecs(after dst.SpaceType) ExprDecsDSL {
	return ExprDecsDSL{Obj: dst.ExprStmtDecorations{
		NodeDecs: dst.NodeDecs{After: after},
	}}
}

// FieldDSL translates to a dst.Field
type FieldDSL struct{ Obj *dst.Field }

// Field creates a new FieldDSL
func Field(typ dst.Expr) FieldDSL {
	return FieldDSL{Obj: &dst.Field{Type: typ}}
}

// Names sets the names of a field
func (d FieldDSL) Names(names ...*dst.Ident) FieldDSL {
	d.Obj.Names = names
	return d
}

// Decs adds decorations to a FieldDSL
func (d FieldDSL) Decs(decs dst.FieldDecorations) FieldDSL {
	d.Obj.Decs = decs
	return d
}

// FieldDecsDSL translates to a dst.FieldDecorations
type FieldDecsDSL struct{ Obj dst.FieldDecorations }

// FieldDecs creates a new FieldDecsDSL
func FieldDecs(before, after dst.SpaceType) FieldDecsDSL {
	return FieldDecsDSL{Obj: dst.FieldDecorations{
		NodeDecs: dst.NodeDecs{Before: before, After: after},
	}}
}

// FieldList translates to a dst.FieldList
func FieldList(fields ...*dst.Field) *dst.FieldList {
	return &dst.FieldList{List: fields}
}

// FnDSL translates to a dst.GenDecl containing a function type
type FnDSL struct{ Obj *dst.FuncDecl }

// Fn creates a new FnDSL
func Fn(name string) FnDSL {
	return FnDSL{Obj: &dst.FuncDecl{Name: Id(name), Type: &dst.FuncType{}}}
}

// Recv specifies the receiver for a function
func (d FnDSL) Recv(fields ...*dst.Field) FnDSL {
	d.Obj.Recv = FieldList(fields...)
	return d
}

// ParamList specifies a parameter FieldList for a function
func (d FnDSL) ParamList(fieldList *dst.FieldList) FnDSL {
	d.Obj.Type.Params = fieldList
	return d
}

// Params specifies parameters for a function
func (d FnDSL) Params(fields ...*dst.Field) FnDSL {
	d.Obj.Type.Params = FieldList(fields...)
	return d
}

// ResultList specifies a result FieldList for a function
func (d FnDSL) ResultList(fieldList *dst.FieldList) FnDSL {
	d.Obj.Type.Results = fieldList
	return d
}

// Results specifies results for a function
func (d FnDSL) Results(fields ...*dst.Field) FnDSL {
	d.Obj.Type.Results = FieldList(fields...)
	return d
}

func (d FnDSL) TypeParams(fieldList *dst.FieldList) FnDSL {
	d.Obj.Type.TypeParams = fieldList
	return d
}

// Body specifies the body for a function
func (d FnDSL) Body(list ...dst.Stmt) FnDSL {
	d.Obj.Body = Block(list...).Obj
	return d
}

// Decs adds decorations to a FnDSL
func (d FnDSL) Decs(decs dst.FuncDeclDecorations) FnDSL {
	d.Obj.Decs = decs
	return d
}

// FnLitDSL translates to a dst.FuncLit
type FnLitDSL struct{ Obj *dst.FuncLit }

// FnLit creates a new FnLitDSL
func FnLit(typ *dst.FuncType) FnLitDSL {
	return FnLitDSL{Obj: &dst.FuncLit{Type: typ}}
}

// Body specifies a body
func (d FnLitDSL) Body(list ...dst.Stmt) FnLitDSL {
	d.Obj.Body = Block(list...).Obj
	return d
}

// FnTypeDSL translates to a dst.FuncType
type FnTypeDSL struct{ Obj *dst.FuncType }

// FnType creates a new FnTypeDSL
func FnType(paramFieldList *dst.FieldList) FnTypeDSL {
	return FnTypeDSL{Obj: &dst.FuncType{Params: paramFieldList}}
}

// Results adds a result field list
func (d FnTypeDSL) Results(resultFieldList *dst.FieldList) FnTypeDSL {
	d.Obj.Results = resultFieldList
	return d
}

// ForDSL translatesto a dst.ForStmt
type ForDSL struct{ Obj *dst.ForStmt }

// For returns a new ForDSL
func For(init dst.Stmt) ForDSL {
	return ForDSL{Obj: &dst.ForStmt{Init: init}}
}

// Cond specifies the condition of a for statement
func (d ForDSL) Cond(cond dst.Expr) ForDSL {
	d.Obj.Cond = cond
	return d
}

// Post specifies the post statement of a for statement
func (d ForDSL) Post(post dst.Stmt) ForDSL {
	d.Obj.Post = post
	return d
}

// Body defines the body of a for statement
func (d ForDSL) Body(list ...dst.Stmt) ForDSL {
	d.Obj.Body = Block(list...).Obj
	return d
}

// Id returns a named dst.Ident
func Id(name string) *dst.Ident {
	return dst.NewIdent(name)
}

// Idf returns a formatted dst.Ident
func Idf(format string, a ...interface{}) *dst.Ident {
	return Id(fmt.Sprintf(format, a...))
}

// IdPath returns a dst.Ident with a name and path
func IdPath(name, path string) *dst.Ident {
	return &dst.Ident{Name: name, Path: path}
}

// IfDSL translates to a dst.IfStmt
type IfDSL struct{ Obj *dst.IfStmt }

// If creates a new If
func If(cond dst.Expr) IfDSL {
	return IfDSL{Obj: &dst.IfStmt{Cond: cond}}
}

// Body specifies the body of the If
func (d IfDSL) Body(list ...dst.Stmt) IfDSL {
	d.Obj.Body = Block(list...).Obj
	return d
}

// Else specifies an optional else branch of the If
func (d IfDSL) Else(els ...dst.Stmt) IfDSL {
	if len(els) == 1 {
		d.Obj.Else = els[0]
	} else {
		d.Obj.Else = Block(els...).Obj
	}
	return d
}

// Decs adds decorations to a IfDSL
func (d IfDSL) Decs(decs dst.IfStmtDecorations) IfDSL {
	d.Obj.Decs = decs
	return d
}

// IfDecsDSL translates to a dst.IfStmtDecorations
type IfDecsDSL struct{ Obj dst.IfStmtDecorations }

// IfDecs creates a new IfDecsDSL
func IfDecs(after dst.SpaceType) IfDecsDSL {
	return IfDecsDSL{Obj: dst.IfStmtDecorations{
		NodeDecs: dst.NodeDecs{After: after},
	}}
}

// IncStmt creates a dst.IncDecStmt for incrementing an expression
func IncStmt(x dst.Expr) *dst.IncDecStmt {
	return &dst.IncDecStmt{X: x, Tok: token.INC}
}

// IndexDSL translates to a dst.IndexListExpr
type IndexDSL struct{ Obj *dst.IndexListExpr }

// Index creates a new IndexDSL
func Index(x dst.Expr) IndexDSL {
	return IndexDSL{Obj: &dst.IndexListExpr{X: x}}
}

// Sub specifies the sub-expression
func (d IndexDSL) Sub(index ...dst.Expr) IndexDSL {
	d.Obj.Indices = index
	return d
}

// KeyValueDSL translates to a dst.KeyValueExpr
type KeyValueDSL struct{ Obj *dst.KeyValueExpr }

// Key creates a new KeyValueDSL
func Key(key dst.Expr) KeyValueDSL {
	return KeyValueDSL{Obj: &dst.KeyValueExpr{Key: key}}
}

// Value specifies the value
func (d KeyValueDSL) Value(value dst.Expr) KeyValueDSL {
	d.Obj.Value = value
	return d
}

// Decs adds decorations to a KeyValueDSL
func (d KeyValueDSL) Decs(decs dst.KeyValueExprDecorations) KeyValueDSL {
	d.Obj.Decs = decs
	return d
}

// KeyValueDecsDSL translates to a dst.KeyValueExprDecorations
type KeyValueDecsDSL struct{ Obj dst.KeyValueExprDecorations }

// KeyValueDecs creates a new KeyValueDecsDSL
func KeyValueDecs(before dst.SpaceType) KeyValueDecsDSL {
	return KeyValueDecsDSL{Obj: dst.KeyValueExprDecorations{
		NodeDecs: dst.NodeDecs{Before: before},
	}}
}

// After adds decorations after the KeyValueDSL
func (d KeyValueDecsDSL) After(after dst.SpaceType) KeyValueDecsDSL {
	d.Obj.After = after
	return d
}

// LitInt returns a dst.BasicLit with a literal int value
func LitInt(value int) *dst.BasicLit {
	return &dst.BasicLit{Kind: token.INT, Value: strconv.Itoa(value)}
}

// LitString returns a dst.BasicLit with a literal string value
func LitString(value string) *dst.BasicLit {
	return &dst.BasicLit{Kind: token.STRING, Value: "\"" + value + "\""}
}

// LitStringf returns a formatted dst.BasicLit with a literal string value
func LitStringf(format string, a ...interface{}) *dst.BasicLit {
	return &dst.BasicLit{
		Kind:  token.STRING,
		Value: fmt.Sprintf("\""+format+"\"", a...),
	}
}

// MapTypeDSL translates to a dst.MapType
type MapTypeDSL struct{ Obj *dst.MapType }

// MapType returns a new MapTypeDSL
func MapType(key dst.Expr) MapTypeDSL {
	return MapTypeDSL{Obj: &dst.MapType{Key: key}}
}

// Value specifies the value expression of a dst.MapType
func (d MapTypeDSL) Value(value dst.Expr) MapTypeDSL {
	d.Obj.Value = value
	return d
}

// NodeDecsf formats a string into a standard node decoration, chopping down
// the line as necessary. The format param should be a single long line
// starting with a line comment (//). The format line is broken into multiple
// lines each with at least one word.
func NodeDecsf(format string, a ...interface{}) dst.NodeDecs {
	const (
		sep     = " "
		maxLine = 79
	)
	full := fmt.Sprintf(format, a...)
	words := strings.Split(full, sep)
	comment := words[0]
	start := dst.Decorations{}
	line := comment
	curWords := 0
	for _, w := range words[1:] {
		if len(line)+len(sep)+len(w) > maxLine && curWords > 0 {
			start = append(start, line)
			line = comment
			curWords = 0
		}
		line += sep + w
		curWords++
	}
	if len(line) > 0 {
		start = append(start, line)
	}
	return dst.NodeDecs{
		Before: dst.NewLine,
		Start:  start,
	}
}

// Paren translates to a dst.ParenExpr
func Paren(x dst.Expr) *dst.ParenExpr {
	return &dst.ParenExpr{X: x}
}

// RangeDSL translates to a dst.RangeStmt
type RangeDSL struct{ Obj *dst.RangeStmt }

// Range returns a new RangeDSL
func Range(x dst.Expr) RangeDSL {
	return RangeDSL{Obj: &dst.RangeStmt{X: x}}
}

// Key sets the key of a range statement
func (d RangeDSL) Key(key dst.Expr) RangeDSL {
	d.Obj.Key = key
	return d
}

// Value sets the value of a range statement
func (d RangeDSL) Value(value dst.Expr) RangeDSL {
	d.Obj.Value = value
	return d
}

// Tok sets the token of a range statement
func (d RangeDSL) Tok(tok token.Token) RangeDSL {
	d.Obj.Tok = tok
	return d
}

// Body defines the body of a range statement
func (d RangeDSL) Body(list ...dst.Stmt) RangeDSL {
	d.Obj.Body = Block(list...).Obj
	return d
}

// Return returns a dst.ReturnStmt
func Return(results ...dst.Expr) *dst.ReturnStmt {
	return &dst.ReturnStmt{Results: results}
}

// SelDSL translates to a dst.SelectorExpr
type SelDSL struct{ Obj *dst.SelectorExpr }

// Sel creates a new SelDSL
func Sel(x dst.Expr) SelDSL {
	return SelDSL{Obj: &dst.SelectorExpr{X: x}}
}

// Dot specifies what is selected
func (d SelDSL) Dot(sel *dst.Ident) SelDSL {
	d.Obj.Sel = sel
	return d
}

// SliceExprDSL translates to a dst.SliceExpr
type SliceExprDSL struct{ Obj *dst.SliceExpr }

// SliceExpr creates a new slice expression
func SliceExpr(x dst.Expr) SliceExprDSL {
	return SliceExprDSL{Obj: &dst.SliceExpr{X: x}}
}

// Low specifies the low expression of a slice expression
func (d SliceExprDSL) Low(low dst.Expr) SliceExprDSL {
	d.Obj.Low = low
	return d
}

// High specifies the high expression of a slice expression
func (d SliceExprDSL) High(high dst.Expr) SliceExprDSL {
	d.Obj.High = high
	return d
}

// SliceType returns a dst.ArrayType representing a slice
func SliceType(elt dst.Expr) *dst.ArrayType {
	return &dst.ArrayType{Elt: elt}
}

// Star returns a dst.StarExpr
func Star(x dst.Expr) *dst.StarExpr {
	return &dst.StarExpr{X: x}
}

// Struct returns a dst.StructType
func Struct(fields ...*dst.Field) *dst.StructType {
	return StructFromList(FieldList(fields...))
}

// StructFromList returns a dst.StructType given a dst.FieldList
func StructFromList(fieldList *dst.FieldList) *dst.StructType {
	if fieldList == nil || len(fieldList.List) == 0 {
		fieldList = emptyFieldList()
	}
	return &dst.StructType{Fields: fieldList}
}

// TypeSpecDSL translates to a dst.TypeSpec
type TypeSpecDSL struct{ Obj *dst.TypeSpec }

// TypeSpec creates a new TypeSpecDSL
func TypeSpec(name string) TypeSpecDSL {
	return TypeSpecDSL{Obj: &dst.TypeSpec{Name: Id(name)}}
}

// Type creates a new TypeSpecDSL
func (d TypeSpecDSL) Type(typ dst.Expr) TypeSpecDSL {
	d.Obj.Type = typ
	return d
}

// TypeParams adds type parameters to TypeDeclDSL
func (d TypeSpecDSL) TypeParams(typeParams *dst.FieldList) TypeSpecDSL {
	d.Obj.TypeParams = typeParams
	return d
}

// TypeDeclDSL translates various types into a dst.GenDecl
type TypeDeclDSL struct{ Obj *dst.GenDecl }

// TypeDecl creates a new TypeDeclDSL
func TypeDecl(typeSpec *dst.TypeSpec) TypeDeclDSL {
	return TypeDeclDSL{
		Obj: &dst.GenDecl{Tok: token.TYPE, Specs: []dst.Spec{typeSpec}},
	}
}

// Decs adds decorations to a TypeDeclDSL
func (d TypeDeclDSL) Decs(decs dst.GenDeclDecorations) TypeDeclDSL {
	d.Obj.Decs = decs
	return d
}

// Un returns a dst.UnaryExpr
func Un(op token.Token, x dst.Expr) *dst.UnaryExpr {
	return &dst.UnaryExpr{Op: op, X: x}
}

// ValueDSL translates to a dst.ValueSpec
type ValueDSL struct{ Obj *dst.ValueSpec }

// Value creates a new ValueDSL
func Value(typ dst.Expr) ValueDSL {
	return ValueDSL{Obj: &dst.ValueSpec{Type: typ}}
}

// Names sets the names of a value
func (d ValueDSL) Names(names ...*dst.Ident) ValueDSL {
	d.Obj.Names = names
	return d
}

// Values sets the values of the spec
func (d ValueDSL) Values(values ...dst.Expr) ValueDSL {
	d.Obj.Values = values
	return d
}

// Var returns a var dst.DeclStmt
func Var(specs ...dst.Spec) *dst.DeclStmt {
	return &dst.DeclStmt{Decl: VarDecl(specs...).Obj}
}

// VarDeclDSL translates variable declaration into a dst.GenDecl
type VarDeclDSL struct{ Obj *dst.GenDecl }

// VarDecl creates a new VarDeclDSL
func VarDecl(specs ...dst.Spec) VarDeclDSL {
	return VarDeclDSL{Obj: &dst.GenDecl{Tok: token.VAR, Specs: specs}}
}

// Decs adds decorations to a VarDeclDSL
func (d VarDeclDSL) Decs(decs dst.GenDeclDecorations) VarDeclDSL {
	d.Obj.Decs = decs
	return d
}

// emptyFieldList returns an empty field list (renders as `struct {}` with no
// new lines)
func emptyFieldList() *dst.FieldList {
	return &dst.FieldList{Opening: true, Closing: true}
}

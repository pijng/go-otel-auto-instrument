package main

import (
	"fmt"
	"go/token"
	"strconv"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/pijng/goinject"
)

// lnOffset is used to calculate the line number of the function where instrumentation occurs.
// The current offset value MUST equal the number of lines added by the instrumentation + 1.
// The calculated value will then be used as a /*line */ directive.
const lnOffset = 5 + 1

type autoinstrumentModifier struct{}

func (mm autoinstrumentModifier) Modify(f *dst.File, dec *decorator.Decorator, res *decorator.Restorer) *dst.File {
	for _, decl := range f.Decls {
		funcDecl, isFunc := decl.(*dst.FuncDecl)
		if !isFunc {
			continue
		}

		var contextArgName string
		// for _, param := range funcDecl.Type.Params.List {
		// 	paramIdent, isIdent := param.Type.(*dst.Ident)
		// 	if !isIdent {
		// 		continue
		// 	}

		// 	if paramIdent.Path == "context" && paramIdent.Name == "Context" && param.Names[0].Name != "_" {
		// 		contextArgName = param.Names[0].Name
		// 		break
		// 	}
		// }

		for _, param := range funcDecl.Type.Params.List {
			starExpr, isPointer := param.Type.(*dst.StarExpr)
			if !isPointer {
				continue
			}
			paramIdent, isIdent := starExpr.X.(*dst.Ident)
			if !isIdent {
				continue
			}

			if paramIdent.Path == "net/http" && paramIdent.Name == "Request" && param.Names[0].Name != "_" {
				contextArgName = param.Names[0].Name + ".Context()"
			}
		}

		funcName := f.Name.Name + "." + funcDecl.Name.Name

		line, err := getLine(dec, funcDecl.Type)
		if err != nil {
			panic(err)
		}

		addLineDirectiveToFuncDecl(funcDecl, line, lnOffset)

		spanStmt := buildSpan(funcName, contextArgName)
		funcDecl.Body.List = append(spanStmt.List, funcDecl.Body.List...)
	}

	return f
}

func main() {
	goinject.Process(autoinstrumentModifier{})
}

func buildSpan(funcName string, contextArgName string) dst.BlockStmt {
	spanFunc := "StartSpan"
	if contextArgName != "" {
		spanFunc = "StartSpanCtx"
	}

	args := make([]dst.Expr, 0)
	if contextArgName != "" {
		args = append(args, &dst.Ident{Name: contextArgName})
	}
	args = append(args, &dst.BasicLit{Kind: token.STRING, Value: strconv.Quote(funcName)})

	lhs := make([]dst.Expr, 0)

	newCtxName := contextArgName
	if newCtxName == "" || strings.Contains(newCtxName, ".") {
		newCtxName = "_"
	}
	oldCtxName := "goai_oldCtx"
	lhs = append(lhs, &dst.Ident{Name: newCtxName})
	lhs = append(lhs, &dst.Ident{Name: oldCtxName})
	lhs = append(lhs, &dst.Ident{Name: "goai_Span"})

	return dst.BlockStmt{
		List: []dst.Stmt{
			// gomSpan := gootelinstrument.StartSpan(funcName)
			&dst.AssignStmt{
				Lhs: lhs,
				Tok: token.DEFINE, // :=
				Rhs: []dst.Expr{
					&dst.CallExpr{
						Fun:  &dst.Ident{Path: "github.com/pijng/go_otel_auto_instrument", Name: spanFunc},
						Args: args,
					},
				},
			},
			// defer func() { gomSpan.End() }()
			&dst.DeferStmt{
				Call: &dst.CallExpr{
					Fun: &dst.FuncLit{
						Type: &dst.FuncType{
							Params: &dst.FieldList{},
						},
						Body: &dst.BlockStmt{
							List: []dst.Stmt{
								&dst.ExprStmt{
									X: &dst.CallExpr{
										Fun: &dst.SelectorExpr{
											X:   &dst.Ident{Name: "goai_Span"},
											Sel: &dst.Ident{Name: "End"},
										},
									},
								},
								&dst.ExprStmt{
									X: &dst.CallExpr{
										Fun: &dst.Ident{Path: "github.com/pijng/go_otel_auto_instrument", Name: "SetContext"},
										Args: []dst.Expr{
											&dst.Ident{Name: oldCtxName},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func addLineDirectiveToFuncDecl(funcDecl *dst.FuncDecl, funcLine int, offset int) {
	funcDecl.Decs.NodeDecs.Start.Append(fmt.Sprintf("/*line :%d:1*/", funcLine-offset))
}

func getLine(dec *decorator.Decorator, funcType *dst.FuncType) (int, error) {
	astNode := dec.Map.Ast.Nodes[funcType]
	if astNode != nil {
		pos := astNode.Pos()
		return dec.Fset.Position(pos).Line, nil
	}

	return 0, fmt.Errorf("could not find a position for: %v", funcType)
}

package analysis

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/types"
)

func HasSideEffect(n ast.Node, info *Info) bool {
	v := hasSideEffectVisitor{info: info}
	ast.Walk(&v, n)
	return v.hasSideEffect
}

type hasSideEffectVisitor struct {
	info          *Info
	hasSideEffect bool
}

func (v *hasSideEffectVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if v.hasSideEffect {
		return nil
	}
	switch n := node.(type) {
	case *ast.CallExpr:
		if _, isSig := v.info.Types[n.Fun].Type.(*types.Signature); isSig { // skip conversions
			v.hasSideEffect = true
			return nil
		}
	case *ast.UnaryExpr:
		if n.Op == token.ARROW {
			v.hasSideEffect = true
			return nil
		}
	}
	return v
}
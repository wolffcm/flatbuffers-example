package evaluation

import (
	"errors"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/wolffcm/flatbuffers-example/tree"
)

// EvalFromBytes evaluates the expression tree contained in the given flatbuffer.
func EvalFromBytes(bs []byte) (float64, error) {
	root := tree.GetRootAsRoot(bs, 0)

	unionTable := new(flatbuffers.Table)
	if hasExpr := root.Expr(unionTable); !hasExpr {
		return 0.0, errors.New("missing root expr")
	}
	return dispatchFromUnionTable(root.ExprType(), unionTable)
}

func dispatchFromUnionTable(ty tree.Node, tbl *flatbuffers.Table) (float64, error) {
	var v float64
	var err error
	switch ty {
	case tree.NodeConstant:
		c := new(tree.Constant)
		c.Init(tbl.Bytes, tbl.Pos)
		v, err = evalConstant(c)
	case tree.NodeBinOp:
		bo := new(tree.BinOp)
		bo.Init(tbl.Bytes, tbl.Pos)
		v, err = evalBinOp(bo)
	default:
		return 0.0, errors.New("unhandled node type: " + tree.EnumNamesNode[ty])
	}
	if err != nil {
		return 0.0, err
	}
	return v, nil

}

func evalConstant(c *tree.Constant) (float64, error) {
	return c.Value(), nil
}

func evalBinOp(bo *tree.BinOp) (float64, error) {

	var lhs float64
	{
		fbt := new(flatbuffers.Table)
		if hasLHS := bo.Lhs(fbt); !hasLHS {
			return 0.0, errors.New("missing LHS")
		}
		var err error
		if lhs, err = dispatchFromUnionTable(bo.LhsType(), fbt); err != nil {
			return 0.0, err
		}
	}

	var rhs float64
	{
		fbt := new(flatbuffers.Table)
		if hasRHS := bo.Rhs(fbt); !hasRHS {
			return 0.0, errors.New("missing RHS")
		}
		var err error
		if rhs, err = dispatchFromUnionTable(bo.RhsType(), fbt); err != nil {
			return 0.0, err
		}
	}

	switch o := bo.Op(); o {
	case tree.OpAdd:
		return lhs + rhs, nil
	case tree.OpSubtract:
		return lhs - rhs, nil
	case tree.OpMultiply:
		return lhs * rhs, nil
	case tree.OpDivide:
		return lhs / rhs, nil
	default:
		return 0.0, errors.New("unknown bin op: " + tree.EnumNamesOp[o])
	}
}
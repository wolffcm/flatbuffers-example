package evaluation

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/wolffcm/flatbuffers-example/tree"
	"math/rand"
)

// Generate creates an expression tree with the given number of levels.
// The produced tree will have 2^levels nodes.
func Generate(levels int) ([]byte, error) {
	b := flatbuffers.NewBuilder(1024)
	expr, err := generateHelper(b, levels)
	if err != nil {
		return nil, err
	}

	tree.RootStart(b)
	tree.RootAddExprType(b, tree.NodeBinOp)
	tree.RootAddExpr(b, expr)
	root := tree.RootEnd(b)

	b.Finish(root)
	bs := b.FinishedBytes()
	return bs, nil
}

func genConst(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	val := float64(rand.Intn(100))
	tree.ConstantStart(b)
	tree.ConstantAddValue(b, val)
	return tree.ConstantEnd(b)
}

func generateHelper(b *flatbuffers.Builder, levels int) (flatbuffers.UOffsetT, error) {
	if levels < 1 {
		return genConst(b), nil
	}

	lhs, err := generateHelper(b, levels - 1)
	if err != nil {
		return 0, nil
	}
	rhs, err := generateHelper(b, levels - 1)
	if err != nil {
		return 0, nil
	}

	var operandType tree.Node
	if levels > 1 {
		operandType = tree.NodeBinOp
	} else {
		operandType = tree.NodeConstant
	}

	op := tree.Op(rand.Intn(4))
	tree.BinOpStart(b)
	tree.BinOpAddLhsType(b, operandType)
	tree.BinOpAddLhs(b, lhs)
	tree.BinOpAddRhsType(b, operandType)
	tree.BinOpAddRhs(b, rhs)
	tree.BinOpAddOp(b, op)
	return tree.BinOpEnd(b), nil
}

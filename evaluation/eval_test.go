package evaluation_test

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/wolffcm/flatbuffers-example/evaluation"
	"github.com/wolffcm/flatbuffers-example/tree"
	"testing"
)

func TestEvalFromBytes(t *testing.T) {
	b := flatbuffers.NewBuilder(1024)

	tree.ConstantStart(b)
	tree.ConstantAddValue(b, 3.0)
	lhs := tree.ConstantEnd(b)

	tree.ConstantStart(b)
	tree.ConstantAddValue(b, 7.0)
	rhs := tree.ConstantEnd(b)

	tree.BinOpStart(b)
	tree.BinOpAddLhsType(b, tree.NodeConstant)
	tree.BinOpAddLhs(b, lhs)
	tree.BinOpAddRhsType(b, tree.NodeConstant)
	tree.BinOpAddRhs(b, rhs)
	tree.BinOpAddOp(b, tree.OpAdd)
	add := tree.BinOpEnd(b)

	tree.RootStart(b)
	tree.RootAddExprType(b, tree.NodeBinOp)
	tree.RootAddExpr(b, add)
	root := tree.RootEnd(b)

	b.Finish(root)
	bs := b.FinishedBytes()

	ans, err := evaluation.EvalFromBytes(bs)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := 10.0, ans; want != got {
		t.Fatalf("got wrong answer, wanted %v, got %v", want, got)
	}
}

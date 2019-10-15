package evaluation_test

import (
	"math"
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/wolffcm/flatbuffers-example/evaluation"
	"github.com/wolffcm/flatbuffers-example/tree"
)

func TestEvalFromBytesInRust(t *testing.T) {
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

	ans, err := evaluation.EvalFromBytesInRust(bs)
	if err != nil {
		t.Fatal(err)
	}
	if ans != 10.0 {
		t.Fatalf("wanted 10.0, got %v", ans)
	}
}

func TestGetBytesFromRust(t *testing.T) {
	bs, free, err := evaluation.GetBytesFromRust()
	if err != nil {
		t.Fatal(err)
	}
	defer free()
	ans, err := evaluation.EvalFromBytes(bs)
	if ans != 10.0 {
		t.Fatalf("wanted 10.0, got %v", ans)
	}
}

func BenchmarkEvalFromBytesInRust_Big(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs, err := evaluation.Generate(10)
		if err != nil {
			b.Fatal(err)
		}
		_, err = evaluation.EvalFromBytesInRust(bs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestEvalFromBytesInRust_Big(t *testing.T) {
	bs, err := evaluation.Generate(10)
	if err != nil {
		t.Fatal(err)
	}
	goResult, err := evaluation.EvalFromBytes(bs)
	if err != nil {
		t.Fatal(err)
	}
	rustResult, err := evaluation.EvalFromBytesInRust(bs)
	if err != nil {
		t.Fatal(err)
	}
	if goResult != rustResult && (!math.IsNaN(goResult) || !math.IsNaN(rustResult)) {
		t.Fatalf("go result was %v; rust result was %v", goResult, rustResult)
	}
}


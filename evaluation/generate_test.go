package evaluation_test

import (
	"fmt"
	"testing"

	"github.com/wolffcm/flatbuffers-example/evaluation"
)

func TestGenerate(t *testing.T) {
	bs, err := evaluation.Generate(5)
	if err != nil {
		t.Fatal(err)
	}
	f, err := evaluation.EvalFromBytes(bs)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(f)
}
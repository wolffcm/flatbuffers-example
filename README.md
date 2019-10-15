# Flatbuffers Example

This repo contains tests that exercise passing flatbuffers between Go and Rust.

A simple flatbuffers schema for an expression tree is contained in `tree.fbs`.
Generated code produced by `flatc --rust` is in `./rust/expr_tree/src/tree_generated.rs`.
The Go counterpart is in `./tree/*.go`.

Go code in `evaluation` can build simple expression trees as flatbuffers and evaluate
them.  
Likewise, Rust code in `rust/expr_tree/src/lib.rs` does the same.

To build:
```
cd rust/expr_tree
cargo build --lib --release
cd -
go test ./...
```
Relevant tests are in `./evaluation/rust_eval_test.go`.

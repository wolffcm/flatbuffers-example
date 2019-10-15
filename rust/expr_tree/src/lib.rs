
extern crate flatbuffers;

mod tree_generated;

use std::os::raw::{c_char, c_int};
use tree_generated::tree::{
    self,
    BinOp, BinOpArgs,
    Constant, ConstantArgs,
    Node,
    Op,
    Root, RootArgs,
};

/// Evaluates the tree assumed to be contained in the storage pointed at by buf.
/// Does not assume ownership or release any memory.
#[no_mangle]
pub fn eval_from_c(buf: *const c_char, len:c_int) -> f64 {
    let bytes = unsafe {
        let u8ptr = buf as *const u8;
        std::slice::from_raw_parts(u8ptr, len as usize)
    };
    eval_from_bytes(bytes)
}

/// Produces a simple expression tree contained in a flatbuffer, and returns a pointer to it.
/// This function releases ownership of the memory returned; to free it, call free_expr_tree().
/// Output parameter len will contain the total size of the referenced memory.
/// Output parameter offset will contain the offset within the buffer where data actually begins.
#[no_mangle]
pub fn get_expr_tree(len: *mut c_int, offset: *mut c_int) -> *const c_char {
    let mut builder = flatbuffers::FlatBufferBuilder::new_with_capacity(1024);
    build_simple_tree(&mut builder);

    // Destroy the builder and take ownership of its buffer
    let (vec, vec_offset) = builder.collapse();
    let sl = vec.into_boxed_slice();
    let sl_len = sl.len();
    let raw = std::boxed::Box::into_raw(sl);
    unsafe {
        *len = sl_len as c_int;
        *offset = vec_offset as c_int;
    }
    raw as *const c_char
}

/// Reclaims ownership and releases memory previously allocated in get_expr_tree().
#[no_mangle]
pub fn free_expr_tree(buf: *mut c_char, len: c_int) {
    let vec = unsafe {
        Vec::from_raw_parts(buf, len as usize, len as usize)
    };
    drop(vec)
}

/// Evaluates an expression tree contained in a flatbuffer.
pub fn eval_from_bytes(buf:&[u8]) -> f64 {
    let root = tree::get_root_as_root(buf);
    match root.expr_type() {
        tree::Node::Constant => eval_constant(root.expr_as_constant().unwrap()),
        tree::Node::BinOp => eval_bin_op(root.expr_as_bin_op().unwrap()),
        _ => panic!(),
    }
}

fn eval_constant(c:tree::Constant) -> f64 {
    c.value()
}

fn eval_bin_op(binop:tree::BinOp) -> f64 {
    let lhs = match binop.lhs_type() {
        Node::Constant => eval_constant(binop.lhs_as_constant().unwrap()),
        Node::BinOp => eval_bin_op(binop.lhs_as_bin_op().unwrap()),
        _ => panic!(),
    };
    let rhs = match binop.rhs_type() {
        Node::Constant => eval_constant(binop.rhs_as_constant().unwrap()),
        Node::BinOp => eval_bin_op(binop.rhs_as_bin_op().unwrap()),
        _ => panic!(),
    };
    match binop.op() {
        Op::Add => lhs + rhs,
        Op::Subtract => lhs - rhs,
        Op::Multiply => lhs * rhs,
        Op::Divide => lhs / rhs,
    }
}

// Builds a simple tree for the expression `3.0 + 7.0`.
fn build_simple_tree(builder: &mut flatbuffers::FlatBufferBuilder) {
    let c1 = crate::Constant::create(builder, &crate::ConstantArgs{value: 3.0});
    let c2 = crate::Constant::create(builder, &crate::ConstantArgs{value: 7.0});
    let add = crate::BinOp::create( builder, &crate::BinOpArgs{
        lhs_type: crate::Node::Constant,
        lhs: Some(c1.as_union_value()),
        rhs_type: crate::Node::Constant,
        rhs: Some(c2.as_union_value()),
        op: crate::Op::Add,
    });
    let root = crate::Root::create(builder, &crate::RootArgs{
        expr_type: crate::Node::BinOp,
        expr: Some(add.as_union_value()),
    });
    builder.finish(root, None);
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        let mut builder = flatbuffers::FlatBufferBuilder::new_with_capacity(1024);
        crate::build_simple_tree(&mut builder);
        let buf = builder.finished_data();
        let got = crate::eval_from_bytes(buf);
        assert_eq!(got, 10.0);
    }
}

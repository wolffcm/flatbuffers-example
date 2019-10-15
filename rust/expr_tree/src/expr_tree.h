
// Evaluates the expression tree flatbuffer pointed at by buf.
double eval_from_c(const char* buf, int len);

// Returns pointer to a simple expression tree contained in a flatbuffer.
//   - len is an output parameter that will contain the total size of the buffer
//   - offset is an output parameter that will contain the offset
//     within the buffer that data begins.
// (Flatbuffers seems to write a buffer starting at higher-addressed memory and going backwards,
// hence the offset to the first valid byte must be returned.)
const char* get_expr_tree(int *len, int *offset);

// Frees memory previously allocated from get_expr_tree().
void free_expr_tree(char* buf, int len);


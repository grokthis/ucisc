module page_decoder(
  input [15:0] instruction,
  output source_immediate,
  output source_memory,
  output [6:0] immediate,
  output [4:0] alu_code,
  output pre_increment,
  output post_increment,
  output decrement,
  output [2:0] source_select,
  output [2:0] destination_select,
  output destination_pc,
  output destination_mem,
  output destination_reg,
  output [2:0] effect
);

// Controls if source of 4 is interpreted as flags or immediate
assign source_immediate = 1;
assign source_memory = (source_select[0] | source_select[1]);
assign alu_code = 5'b0;
assign decrement = 0;
assign pre_increment = 0;
assign post_increment = 0;
assign effect[2:0] = 3'b11;
assign immediate = {instruction[5], instruction[5:0]};

// whether the destination is memory or not changes how a
// few things are interpreted.
assign source_select = instruction[9:7];
assign destination_select = 3'b0;
assign destination_pc = 0;
assign destination_mem = 0;
assign destination_reg = 0;

endmodule

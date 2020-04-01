module alu_decoder(
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
assign source_immediate = 0;
assign source_memory = ~source_select[2] & (source_select[0] | source_select[1]);
assign alu_code = instruction[6:2];
assign decrement = sign;
wire sign = instruction[13];

// whether the destination is memory or not changes how a
// few things are interpreted.
assign destination_mem = destination_select[0] | destination_select[1];

assign immediate = 7'b0;

assign pre_increment = 0;
assign post_increment = instruction[12];

assign source_select = instruction[9:7];
assign destination_select = {0,instruction[11:10]};

assign destination_pc = instruction[11] ~& instruction[10];
assign destination_reg = 0;

assign effect[2] = ~post_increment & sign;
assign effect[1:0] = instruction[14:13];

endmodule

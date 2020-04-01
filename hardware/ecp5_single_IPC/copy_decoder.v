module copy_decoder(
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
assign source_memory = ~source_select[2] & (source_select[0] | source_select[1]);
// Copy is always ALU code 0
assign alu_code = 5'b00000;
// Copy always decrements, never increments
assign decrement = 1;
wire control = instruction[6];

// whether the destination is memory or not changes how a
// few things are interpreted.
assign destination_mem = ~destination_select[2] && (destination_select[0] | destination_select[1]);

assign immediate =
  (destination_mem) ? {instruction[5],instruction[5:0]} :
  (~destination_mem) ? instruction[6:0]:
  7'h0;

assign pre_increment = destination_mem & control;
assign post_increment = 0;

assign source_select = instruction[9:7];
assign destination_select = instruction[12:10];

assign destination_pc = destination_select == 3'b000;
assign destination_reg = destination_select[2];

assign effect[2] = 0;
assign effect[1:0] = instruction[14:13];

endmodule

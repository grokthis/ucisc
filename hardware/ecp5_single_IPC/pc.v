// program counter

module pc(
  input clock,
  input reset,
  input store_enabled,
  input destination_pc,
  input [15:0] source_value,
  output reg [15:0] current_pc,
  output reg [15:0] next_pc
);

wire [15:0] inc_pc = current_pc + 1;

always @(negedge clock)
  if (reset)
    next_pc <= 16'h0;
  else if (destination_pc & store_enabled)
    next_pc <= source_value;
  else
    next_pc <= inc_pc;

always @(posedge clock)
  current_pc <= next_pc;

endmodule

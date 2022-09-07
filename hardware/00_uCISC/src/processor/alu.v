module alu(
  input [15:0] source,
  input [15:0] destination,
  input [3:0] op_code,
  input [15:0] flags,
  output [15:0] result_out,
  output overflow,
  output carry
);

wire [31:0] mult_result = source * destination;

wire [31:0] signed_shift_destination = { {16{flags[8] & destination[15]}}, destination };
assign result_out = result[15:0];
assign carry = result[16];
assign overflow = carry;
wire [16:0] result =
    //Bit operations
    op_code == 4'h0 ? source : // Copy
    op_code == 4'h1 ? source & destination : // AND
    op_code == 4'h2 ? source | destination : // OR
    op_code == 4'h3 ? source ^ destination : // XOR
    op_code == 4'h4 ? ~source + 1 : // 2's compliment

    //Shift operations
    op_code == 4'h5 ? { 1'h0, (source[15:4] == 12'h0 ? destination << source[3:0] : 16'h0) } : // Shift left
    // shift signed/unsigned by sign flag
    op_code == 4'h6 ? source > 16'hF ? {16{flags[8] & destination[15]}} : signed_shift_destination >> source[3:0] :
    op_code == 4'h7 ? { source[7:0], source[15:8] } : // Swap
    op_code == 4'h8 ? { source[15:8], 8'h00 } : // High byte
    op_code == 4'h9 ? { 8'h00, source[7:0] } : // Low byte

    op_code == 4'hA ? destination + source : // Add
    op_code == 4'hB ? destination + ~source + 1 : // Subtract
    op_code == 4'hC ? mult_result[15:0] : // Multiply
    op_code == 4'hD ? mult_result[31:16] : // Multiply upper word
    op_code == 4'hE ? destination + source + flags[2] : // Add with carry
    0; // TBD

endmodule

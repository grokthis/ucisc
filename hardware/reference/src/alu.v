module alu(
  input [15:0] source,
  input [15:0] destination,
  input [3:0] op_code,
  input [15:0] flags,
  output [15:0] result_out,
  output [15:0] flags_out,
  output write_flags
);

wire [31:0] mult_result = source * destination;

wire [31:0] signed_shift_destination = { {16{flags[8] & destination[15]}}, destination };
assign result_out = result[15:0];
wire carry = result[16];
wire [16:0] result =
    //Bit operations
    op_code == 4'h0 ? source : // Copy
    op_code == 4'h1 ? source & destination : // AND
    op_code == 4'h2 ? source | destination : // OR
    op_code == 4'h3 ? source ^ destination : // XOR
    op_code == 4'h4 ? ~source : // inverse

    //Shift operations
    op_code == 4'h5 ? destination << source : // Shift left
    // shift signed/unsigned by sign flag
    op_code == 4'h6 ? source > 16'hF ? {16{flags[8] & destination[15]}} : signed_shift_destination >> source[3:0] :
    op_code == 4'h7 ? { source[7:0], source[15:8] } : // Swap
    op_code == 4'h8 ? { source[15:8], 8'h00 } : // High byte
    op_code == 4'h9 ? { 8'h00, source[7:0] } : // Low byte

    op_code == 4'hA ? source + destination : // Add
    op_code == 4'hB ? source - destination : // Subtract
    op_code == 4'hC ? 16'h0 : //mult_result[16:0] : // Multiply
    op_code == 4'hD ? 16'h0 : //source / destination : // Divide
    op_code == 4'hE ? source & overflow : source & destination;

wire zero = result[15:0] == 15'h0;
wire negative = result[15];
wire overflow = carry;
wire divide_error = op_code == 4'hD && source == 16'h0;

assign flags_out = { flags[15:5], divide_error, overflow, carry, negative, zero };
assign write_flags = op_code != 4'h0;
endmodule

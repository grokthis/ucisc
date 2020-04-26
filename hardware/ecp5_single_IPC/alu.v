module alu(
  input [15:0] a,
  input [15:0] b,
  input [4:0] alu_code,
  input carry_in,
  output [15:0] result_out,
  output carry_out,
  output overflow_out,
  output zero_out,
  output negative_out
);

wire signed_op =
  alu_code == 5'h07 ||
  alu_code == 5'h0B ||
  alu_code == 5'h0D ||
  alu_code == 5'h0F;

wire add_result =
  alu_code == 5'h0A ||
  alu_code == 5'h0B ||
  alu_code == 5'h0C ||
  alu_code == 5'h0D;

// We are going to extend the inputs by 1 bit to handle
// signed and unsigned operations in the same logic
wire [16:0] input_a = signed_op ? $signed(a) : a;
wire [16:0] input_b = signed_op ? $signed(b) : b;

wire [31:0] mult_result_16 = input_a * input_b;
// For multiply signed, the high bits including the last bit
// of the valid result (bit index 15) must all be 0
// (positive) or 1 (negative) as a sign carrier for the lower
// bits. If there is a mix of 1's and 0's it encodes numeric
// values and an overflow occured.
assign zero_upper_word = ~(mult_result_16[31] | mult_result_16[30] | mult_result_16[29] | mult_result_16[28] | mult_result_16[27] | mult_result_16[26] | mult_result_16[25] | mult_result_16[24] | mult_result_16[23] | mult_result_16[22] | mult_result_16[21] | mult_result_16[20] | mult_result_16[19] | mult_result_16[18] | mult_result_16[17] | mult_result_16[16]);
assign ones_upper_word = (mult_result_16[31] & mult_result_16[30] & mult_result_16[29] & mult_result_16[28] & mult_result_16[27] & mult_result_16[26] & mult_result_16[25] & mult_result_16[24] & mult_result_16[23] & mult_result_16[22] & mult_result_16[21] & mult_result_16[20] & mult_result_16[19] & mult_result_16[18] & mult_result_16[17] & mult_result_16[16]);
assign mult_16_overflow = signed_op ? (zero_upper_word & ~mult_result_16[15] | ones_upper_word & mult_result_16[15]) : zero_upper_word;

wire [15:0] pre_sum = input_b + carry_in;
wire [15:0] invert_a = ~input_a;

wire [7:0] low_a = input_a[7:0];
wire [7:0] high_a = input_a[15:8];
wire [7:0] low_b = input_b[7:0];
wire [7:0] high_b = input_b[15:8];

wire [15:0] second_adder_arg = alu_code[4:1] == 4'h3 ? invert_a : input_a;
wire [16:0] sum = pre_sum + second_adder_arg;

reg [15:0] result;
always @* begin
  case (alu_code)
    //Bit operations
    5'h00: result = input_a; // copy
    5'h01: result = input_a & input_b; // AND
    5'h02: result = input_a | input_b; // OR
    5'h03: result = input_a ^ input_b; // XOR
    5'h04: result = invert_a; // INV

    //Shift operations
    5'h05: result = input_b << input_a; // Shift left
    5'h06: result = input_b >> input_a; // Shift right zero extend
    5'h07: result = input_b >>> input_a; // Shift right sign extend

    5'h0A: result = sum[15:0]; // Add unsigned
    5'h0B: result = sum[15:0]; // Add Signed
    5'h0C: result = sum[15:0]; // Subtract unsigned
    5'h0D: result = sum[15:0]; // Subtract signed

    5'h0E: result = mult_result_16[15:0]; // Multiply unsigned
    5'h0F: result = mult_result_16[15:0]; // Multiply signed

    5'h12: result = {low_a, high_a}; // Swap MSB and LSB bytes
    5'h13: result = {8'h00, high_a}; // MSB to LSB
    5'h14: result = {high_a, 8'h00}; // Zero LSB
    5'h15: result = {8'h00, low_a}; // Zero MSB
    5'h16: result = {high_a, low_b}; // Write MSB only
    5'h17: result = {high_b, low_a}; // Write LSB only
    5'h18: result = {low_a, low_b}; // Write LSB as MSB
    5'h19: result = {high_b, high_a}; // Write MSB as LSBs
  endcase
end

wire add_overflow = signed_op ? (input_b[15] ~^ input_a[15]) & (input_a[15] ^ result[15]) : sum[16];
wire overflow_out = alu_code[4:1] == 4'h7 ? mult_16_overflow : // multiply
                alu_code[3] ? add_overflow : 0; // assums no mod codes, since not supported

// Carry indicates there is data to chain to the next math op
// This is not the same as an overflow when working with negative
// numbers since 2's compliment relies on carries to work
assign carry_out = add_result ? sum[16] : 0;
assign result_out = result;
assign zero_out = result == 16'h0;
assign negative_out = result[15];

endmodule

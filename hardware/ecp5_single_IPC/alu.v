module alu(
  input clock,
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

// width of 17 to preserve a sign bit
reg [15:0] a0;
reg [15:0] b0;
reg second_step;
reg zero_continued;

wire signed_op =
  alu_code == 5'h07 ||
  alu_code == 5'h09 ||
  alu_code == 5'h0B ||
  alu_code == 5'h0D ||
  alu_code == 5'h0F ||
  alu_code == 5'h11;

// We are going to extend the inputs by 1 bit to handle
// signed and unsigned operations in the same logic
wire [16:0] input_a = signed_op ? $signed(a) : a;
wire [16:0] input_b = signed_op ? $signed(b) : b;

wire [31:0] mult_result_16 = input_a * input_b;
wire [63:0] mult_result_32 = {input_a, a0} * {input_b, b0};

//wire [15:0] div_result_16 = input_a / input_b;
//wire [31:0] div_result_32 = {input_a, a0} / {input_b, b0};

//wire [15:0] mod_result_16 = input_a % input_b;
//wire [31:0] mod_result_32 = {input_a, a0} % {input_b, b0};

always @(posedge clock) begin
  a0 <= input_a;
  b0 <= input_b;
  zero_continued <= zero_out;
  if (alu_code == 5'h08 || alu_code == 5'h09
    || alu_code == 5'h0E || alu_code == 5'h0F
    || alu_code == 5'h10 || alu_code == 5'h11)
    second_step <= ~second_step;
  else
    second_step <= 1'b0;
end

reg [16:0] result;
always @* begin
  case (alu_code)
    //Bit operations
    5'h00: result = input_a; // copy
    5'h01: result = input_a & input_b; // AND
    5'h02: result = input_a | input_b; // OR
    5'h03: result = input_a ^ input_b; // XOR
    5'h04: result = ~input_a; // INV

    //Shift operations
    5'h05: result = input_b << input_a; // Shift left
    5'h06: result = input_b >> input_a; // Shift right zero extend
    5'h07: result = input_b >>> input_a; // Shift right sign extend

    // Arithmetic operations
    //5'h08: result = second_step ? {0, mod_result_32[31:16]} : {0, mod_result_16[15:0]}; // Mod unsigned
    //5'h09: result = second_step ? $signed(mod_result_32[31:16]) : $signed(mod_result_16[15:0]); // Mod signed

    5'h0A: result = input_b + input_a + carry_in; // Add unsigned
    5'h0B: result = input_b + input_a + carry_in; // Add Signed
    5'h0C: result = input_b - input_a + carry_in; // Subtract unsigned
    5'h0D: result = input_b - input_a + carry_in; // Subtract signed

    5'h0E: result = second_step ? {0,mult_result_32[31:16]} : {0,mult_result_16[15:0]}; // Multiply unsigned
    5'h0F: result = second_step ? $signed(mult_result_32[31:16]) : {0,mult_result_16[15:0]}; // Multiply signed
    //5'h10: result = second_step ? {0,div_result_32[31:16]} : {0,div_result_16}; // Divide unsigned
    //5'h11: result = second_step ? $signed(div_result_32[31:16]) : {0,div_result_16}; // Divide signed

    5'h12: result = {input_a[7:0], input_a[15:8]}; // Swap MSB and LSB bytes
    5'h13: result = {8'h00,input_a[15:8]}; // MSB to LSB
    5'h14: result = {input_a[15:8], 8'h00}; // Zero LSB
    5'h15: result = {8'h00, input_a[7:0]}; // Zero MSB
    5'h16: result = {input_a[15:8], input_b[7:0]}; // Write MSB only
    5'h17: result = {input_b[15:8], input_a[7:0]}; // Write LSB only
    5'h18: result = {input_a[7:0], input_b[7:0]}; // Write LSB as MSB
    5'h19: result = {input_b[15:8], input_a[15:8]}; // Write MSB as LSBs
  endcase
end

// Carry indicates there is data to chain to the next math op
// This is not the same as an overflow when working with negative
// numbers since 2's compliment relies on carries to work
assign carry_out = overflow_calc | result[16] > 0;
assign result_out = result[15:0];
assign zero_out = second_step ?
  zero_continued & result_out == 16'h0 :
  result_out == 16'h0;
assign negative_out = result[15];
assign overflow_out = overflow_calc;

// Overflow indicates if the result fits in 16 bits or not
reg overflow_calc;
always @* begin
  case (alu_code)
    //Bit operations
    5'h00: overflow_calc = 0; // copy
    5'h01: overflow_calc = 0; // AND
    5'h02: overflow_calc = 0; // OR
    5'h03: overflow_calc = 0; // XOR
    5'h04: overflow_calc = 0; // INV

    //Shift operations - not yet implemented, need to write a barrel shifter
    5'h05: overflow_calc = 0; // Shift left
    5'h06: overflow_calc = 0; // Shift right zero extend
    5'h07: overflow_calc = 0; // Shift right sign extend

    // Arithmetic operations
    5'h08: overflow_calc = 0; // Mod unsigned
    5'h09: overflow_calc = 0; // Mod signed

    // for add unsigned, overflow is just did it carry
    5'h0A: overflow_calc = result[16]; // Add unsigned
    // for add signed, overflow is if inputs have same sign
    // and output has a different sign
    5'h0B: overflow_calc = (input_b[15] ~^ input_a[15]) & (input_a[15] ^ result[15]); // Add Signed
    // for sub unsigned, overflow is if the result is negative. 
    5'h0C: overflow_calc = result[16]; // Subtract unsigned
    // for sub signed, overflow is if inputs have different sign
    // before negating input_a and output has a different sign
    5'h0D: overflow_calc = (input_b[15] ^ input_a[15]) & (input_b[15] ^ result[15]); // Subtract signed

    // For multiply unsigned, any non-zero bits in position 16 or
    // greater indicate an overflow
    5'h0E: overflow_calc = second_step ? mult_result_32[63:32] != 32'h0 : mult_result_16[31:15] != 16'h0; // Multiply unsigned
    // For multiply signed, the high bits including the last bit
    // of the valid result (bit index 15) must all be 0
    // (positive) or 1 (negative) as a sign carrier for the lower
    // bits. If there is a mix of 1's and 0's it encodes numeric
    // values and an overflow occured.
    5'h0F: overflow_calc = second_step ? (mult_result_32[63:31] != 33'h0 & mult_result_32[63:31] != 33'h1FFFFFFFF) : (mult_result_16[31:15] != 16'h0 & mult_result_16[31:15] != 17'h1FFFF); // Multiply signed
    5'h10: overflow_calc = 0; // Divide unsigned
    5'h11: overflow_calc = 0; // Divide signed

    5'h12: overflow_calc = 0; // Swap MSB and LSB bytes
    5'h13: overflow_calc = 0; // MSB to LSB
    5'h14: overflow_calc = 0; // Zero LSB
    5'h15: overflow_calc = 0; // Zero MSB
    5'h16: overflow_calc = 0; // Write MSB only
    5'h17: overflow_calc = 0; // Write LSB only
    5'h18: overflow_calc = 0; // Write LSB as MSB
    5'h19: overflow_calc = 0; // Write MSB as LSBs
  endcase
end
endmodule

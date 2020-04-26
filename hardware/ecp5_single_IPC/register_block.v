module register_block(
  input clock,
  input [15:0] pc,
  input source_immediate,
  input [6:0] immediate,
  input [15:0] destination_write,
  input [2:0] destination_select,
  input [2:0] source_select, 
  input [15:0] flags_value,
  input set_flags,
  input store_value,
  input pre_increment,
  input post_increment,
  input decrement,
  output [15:0] source_out,
  output [15:0] destination_out,
  output reg [15:0] flags_out
);

reg [15:0] r1;
reg [15:0] r2;
reg [15:0] r3;

reg [15:0] reg_current_value;
always @* begin
  case(source_select[1:0])
    3'h0: reg_current_value = source_select[2] ? (flags_value & {16{source_immediate}}) : pc;
    3'h1: reg_current_value = r1;
    3'h2: reg_current_value = r2;
    3'h3: reg_current_value = r3;
  endcase
end

reg [15:0] dest_base_value;
always @* begin
  case(destination_select[1:0])
    3'h0: dest_base_value = destination_select[2] ? flags_value : pc;
    3'h1: dest_base_value = r1;
    3'h2: dest_base_value = r2;
    3'h3: dest_base_value = r3;
  endcase
end

wire [15:0] full_immediate = {{9{immediate[6]}}, immediate};
assign source_out = reg_current_value + full_immediate;

wire [15:0] pre_increment_value = {{15{decrement & pre_increment}}, pre_increment};
assign destination_out = dest_base_value + pre_increment_value;

wire [15:0] post_increment_value = {{15{decrement & post_increment}}, post_increment};
wire [15:0] reg_incremented = source_out + post_increment_value;
wire should_store_reg = post_increment;

wire store_register = (store_value & destination_select[2]);
wire [15:0] dest_incremented = store_register ? destination_write : destination_out + post_increment_value;
wire should_store_dest = post_increment | store_register;

wire [15:0] r1_result = (should_store_dest & destination_select == 2'h1) ? dest_incremented :
                        (should_store_reg & source_select == 2'h1) ? source_incremented : r1;

wire [15:0] r2_result = (should_store_dest & destination_select == 2'h2) ? dest_incremented :
                        (should_store_reg & source_select == 2'h2) ? source_incremented : r2;

wire [15:0] r3_result = (should_store_dest & destination_select == 2'h3) ? dest_incremented :
                        (should_store_reg & source_select == 2'h3) ? source_incremented : r3;

wire [15:0] flags_result = set_flags ? flags_value :
                           (store_value & destination_select == 3'h4) ? destination_write : flags_out;

always @(posedge clock)
begin
  r1 <= r1_result;
  r2 <= r2_result;
  r3 <= r3_result;
  flags_out <= flags_result;
end

endmodule

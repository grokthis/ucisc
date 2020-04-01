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
  output [15:0] flags_out
);

reg [15:0] flags;
reg [15:0] r1;
reg [15:0] r2;
reg [15:0] r3;

wire [15:0] four_source = (source_immediate) ? $signed(immediate) : flags_value;
reg [15:0] reg_base_value;
always @* begin
  case(source_select)
    3'h0: reg_base_value = pc;
    3'h1: reg_base_value = r1;
    3'h2: reg_base_value = r2;
    3'h3: reg_base_value = r3;
    3'h4: reg_base_value = four_source;
    3'h5: reg_base_value = r1;
    3'h6: reg_base_value = r2;
    3'h7: reg_base_value = r3;
  endcase
end

reg [15:0] dest_base_value;
always @* begin
  case(destination_select)
    3'h0: dest_base_value = pc;
    3'h1: dest_base_value = r1;
    3'h2: dest_base_value = r2;
    3'h3: dest_base_value = r3;
    3'h4: dest_base_value = flags_value;
    3'h5: dest_base_value = r1;
    3'h6: dest_base_value = r2;
    3'h7: dest_base_value = r3;
  endcase
end

assign source_out = reg_base_value + immediate;
assign flags_out = flags;

wire [15:0] increment_value = decrement ? -1 : 1;
wire [15:0] pre_increment_value = pre_increment ? increment_value : 16'h0;
assign destination_out = dest_base_value + pre_increment_value;

always @(posedge clock)
begin
  if (set_flags)
    flags <= flags_value;
  else if (store_value & destination_select == 3'b100)
    flags <= destination_write;
end

wire destination_mem;
assign destination_mem = ~destination_select[2] & (destination_select[0] | destination_select[1]);

wire [15:0] post_increment_value = post_increment ? increment_value : 16'h0;
wire [15:0] r1_incremented = r1 + post_increment_value;
wire [15:0] r2_incremented = r2 + post_increment_value;
wire [15:0] r3_incremented = r3 + post_increment_value;

always @(posedge clock)
begin
  if (post_increment) begin
    if (destination_select == 3'h1 || source_select == 3'h1)
      r1 <= r1_incremented;
    else if (destination_select == 3'h2 || source_select == 3'h2)
      r2 <= r2_incremented;
    else if (destination_select == 3'h3 || source_select == 3'h3)
      r3 <= r3_incremented;
  end
  if (store_value & destination_select[2]) begin
    if (destination_select[1:0] == 2'h1)
      r1 <= destination_write;
    else if (destination_select[1:0] == 2'h2)
      r2 <= destination_write;
    else if (destination_select[1:0] == 2'h3)
      r3 <= destination_write;
  end
end

endmodule

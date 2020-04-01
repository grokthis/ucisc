module memory_block_16kx16 (
  input clock,
  input write_enable,
  input [13:0] portA_address,
  input [13:0] portB_address,
  input [15:0] data_in,
  input [255:0] row_data,
  input row_write,
  output [255:0] row_data_out,
  output reg [15:0] portA_out,
  output reg [15:0] portB_out
);

wire block0_addr = portB_address[1:0] == 2'b00;
wire block0_enable = row_write | (write_enable & block0_addr);
wire [15:0] block0_portA_out;
wire [15:0] block0_portB_out;
wire [63:0] block0_row_out;
memory_block_04kx16 block0(
  .clock(clock),
  .write_enable(block0_enable),
  .portA_address(portA_address[13:2]),
  .portB_address(portB_address[13:2]),
  .data_in(data_in),
  .row_data(row_data[63:0]),
  .row_write(row_write),
  .row_data_out(block0_row_out),
  .portA_out(block0_portA_out),
  .portB_out(block0_portB_out)
);

wire block1_addr = portB_address[1:0] == 2'b01;
wire block1_enable = row_write | (write_enable & block1_addr);
wire [15:0] block1_portA_out;
wire [15:0] block1_portB_out;
wire [63:0] block1_row_out;
memory_block_04kx16 block1(
  .clock(clock),
  .write_enable(block1_enable),
  .portA_address(portA_address[13:2]),
  .portB_address(portB_address[13:2]),
  .data_in(data_in),
  .row_data(row_data[127:64]),
  .row_write(row_write),
  .row_data_out(block1_row_out),
  .portA_out(block1_portA_out),
  .portB_out(block1_portB_out)
);

wire block2_addr = portB_address[1:0] == 2'b10;
wire block2_enable = row_write | (write_enable & block2_addr);
wire [15:0] block2_portA_out;
wire [15:0] block2_portB_out;
wire [63:0] block2_row_out;
memory_block_04kx16 block2(
  .clock(clock),
  .write_enable(block2_enable),
  .portA_address(portA_address[13:2]),
  .portB_address(portB_address[13:2]),
  .data_in(data_in),
  .row_data(row_data[191:128]),
  .row_write(row_write),
  .row_data_out(block2_row_out),
  .portA_out(block2_portA_out),
  .portB_out(block2_portB_out)
);

wire block3_addr = portB_address[1:0] == 2'b11;
wire block3_enable = row_write | (write_enable & block3_addr);
wire [15:0] block3_portA_out;
wire [15:0] block3_portB_out;
wire [63:0] block3_row_out;
memory_block_04kx16 block3(
  .clock(clock),
  .write_enable(block3_enable),
  .portA_address(portA_address[13:2]),
  .portB_address(portB_address[13:2]),
  .data_in(data_in),
  .row_data(row_data[255:192]),
  .row_write(row_write),
  .row_data_out(block3_row_out),
  .portA_out(block3_portA_out),
  .portB_out(block3_portB_out)
);

assign row_data_out = {block0_row_out, block1_row_out, block2_row_out, block3_row_out};
always @* begin
  case (portA_address[1:0])
    2'b00:
      portA_out <= block0_portA_out;
    2'b01:
      portA_out <= block1_portA_out;
    2'b10:
      portA_out <= block2_portA_out;
    2'b11:
      portA_out <= block3_portA_out;
  endcase
  case (portB_address[1:0])
    2'b00:
      portB_out <= block0_portB_out;
    2'b01:
      portB_out <= block1_portB_out;
    2'b10:
      portB_out <= block2_portB_out;
    2'b11:
      portB_out <= block3_portB_out;
  endcase
end

endmodule

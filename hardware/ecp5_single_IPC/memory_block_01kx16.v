module memory_block_01kx16 (
  input clock,
  input write_enable,
  input [9:0] portA_address,
  input [9:0] portB_address,
  input [15:0] data_in,
  output reg [15:0] portA_out,
  output reg [15:0] portB_out
);

reg [15:0] mem[1024:0];

always @(posedge clock) begin
  portA_out <= mem[portA_address[9:0]];
  if (write_enable)
    mem[portB_address[9:0]] <= data_in;
end

always @(negedge clock) begin
  portB_out <= mem[portB_address[9:0]];
end

endmodule

module basic_memory_block (
  input clock,
  input write_enable,
  input [WIDTH-1:0] read_address,
  input [WIDTH-1:0] write_address,
  input [15:0] data_in,
  output reg [15:0] data_out,
);
  parameter WIDTH = 8;
  parameter mem_depth = 1 << WIDTH;

  reg [15:0] mem[mem_depth-1:0];

  always @(posedge clock) begin
      data_out <= mem[read_address];

      if(write_enable)
        mem[write_address[WIDTH-1:0]] <= data_in;
  end
endmodule

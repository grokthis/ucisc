module memory_block (
  input clock,
  input write_enable,
  input [15:0] read_address,
  input [15:0] write_address,
  input [15:0] data_in,
  output reg [15:0] data_out
);

parameter WIDTH = 13;
parameter mem_depth = 1 << WIDTH;
parameter MEM_INIT_FILE = "prog.hex";

reg [15:0] mem[mem_depth-1:0];

initial begin
  if (MEM_INIT_FILE != "") begin
    $readmemh(MEM_INIT_FILE, mem);
  end
end

always @(posedge clock) begin
    if(write_enable)
        mem[write_address[WIDTH-1:0]] <= data_in;

    data_out <= mem[read_address[WIDTH-1:0]];
end

endmodule
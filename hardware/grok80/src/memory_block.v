module memory_block (
  input clock,
  input write_enable,
  input [15:0] read_address,
  input [15:0] write_address,
  input [15:0] data_in,
  output [15:0] data_out,
  output reg [31:0] read_data
);

parameter WORDS = 65536;
parameter ADDRESS_WIDTH = 16;
parameter MEM_INIT_FILE = "prog.hex";

reg [31:0] mem[(WORDS-1)/2:0];

initial begin
  if (MEM_INIT_FILE != "") begin
    $readmemh(MEM_INIT_FILE, mem);
  end
end

assign data_out = read_address[0] ? read_data[31:16] : read_data[15:0];

assign data_out = read_address[0] ? read_data[31:16] : read_data[15:0];

always @(negedge clock) begin
    read_data <= mem[read_address];

    if(write_enable)
      if (write_address[0])
        mem[write_address][31:16] <= data_in;
      else
        mem[write_address][15:0] <= data_in;
end

endmodule

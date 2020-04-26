// Write clock must flip @negedge of clock to setup next address
// Read + Write Timing:
// t0 = @negedge clock, write_clock
// t1 = @posedge clock, write_clock - write to portB if write_enabled
// t2 = @negedge clock, ~write_clock - portB_out valid
// t3 = @posedge clock, ~write_clock - portA_out valid
module memory_block (
  input clock,
  input write_clock,
  input write_enable,
  input [address_width-1:0] portA_address,
  input [address_width-1:0] portB_address,
  input [data_width-1:0] data_in,
  output reg [data_width-1:0] portA_out,
  output reg [data_width-1:0] portB_out
);

parameter data_width = 16;
parameter address_width = 10;
parameter mem_depth = 1 << address_width;

reg [data_width-1:0] mem[mem_depth-1:0];

wire [address_width-1:0] read_address = write_clock ? portB_address : portA_address;

always@(negedge clock) begin
  if(~write_clock)
    portB_out = portA_out;
end

always @(posedge clock) begin
  portA_out <= mem[read_address];
  if(write_clock && write_enable)
    mem[read_address] <= data_in;
end

endmodule

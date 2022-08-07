module memory_block (
  input clock,
  input write_enable,
  input [WIDTH-1:0] read_address,
  input [WIDTH-1:0] write_address,
  input [15:0] data_in,
  output [15:0] data_out,
  output reg [15:0] data_out_a,
  output reg [15:0] data_out_b,
  output reg primary_out_select
);

parameter WIDTH = 15;
parameter mem_depth = 1 << (WIDTH - 1);
parameter MEM_INIT_FILEA = "prog.hex.a.hex";
parameter MEM_INIT_FILEB = "prog.hex.b.hex";

reg [15:0] memA[mem_depth-1:0];
reg [15:0] memB[mem_depth-1:0];

initial begin
  $readmemh(MEM_INIT_FILEA, memA);
  $readmemh(MEM_INIT_FILEB, memB);
end

always@(posedge clock) begin
    primary_out_select <= read_address[0];
end
assign data_out = primary_out_select ? data_out_b : data_out_a;

wire [WIDTH-2:0] readA = read_address[0] ? read_address[WIDTH-1:1] + 1 : read_address[WIDTH-1:1];
always @(posedge clock) begin
    data_out_a <= memA[readA];

    if(write_enable & ~write_address[0])
      memA[write_address[WIDTH-1:1]] <= data_in;
end

wire [WIDTH-2:0] readB = read_address[WIDTH-1:1];
always @(posedge clock) begin
    data_out_b <= memB[readB];

    if(write_enable & write_address[0])
      memB[write_address[WIDTH-1:1]] <= data_in;
end

endmodule

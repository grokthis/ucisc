module dff(
  input clock,
  input [WIDTH-1:0] d,
  input async_reset,
  input enable,
  output reg [WIDTH-1:0] q
);
    parameter WIDTH = 1;
    parameter INIT = {WIDTH{1'b0}};

    initial q = INIT;

    always @ (posedge clock or posedge async_reset)
        if (async_reset) begin
            q <= INIT;
        end else if (enable == 1) begin
            q <= d;
        end
endmodule
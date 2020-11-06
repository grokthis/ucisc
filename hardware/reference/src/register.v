module register(
    input [WIDTH-1:0] data_in,
    input [WIDTH-1:0] increment,
    input capture,
    input write_enable,
    input increment_enable,
    input reset,
    output [WIDTH-1:0] value
);
    parameter WIDTH = 1;
    parameter INIT = {WIDTH{1'b0}};

    wire increment_select = increment_enable & ~write_enable;
    wire data_enable = increment_enable | write_enable;

    wire [WIDTH-1:0] incremented = increment + value;

    wire [WIDTH-1:0] d = increment_select ? incremented : data_in;

    dff #(.WIDTH(WIDTH), .INIT(INIT)) data(
        .clock(capture),
        .d(d),
        .async_reset(reset),
        .enable(data_enable),
        .q(value)
    );
endmodule
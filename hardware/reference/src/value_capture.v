module value_capture (
    input clock,
    input [1:0] current_step,
    input [1:0] capture_on,
    input [WIDTH-1:0] input_value,
    output [WIDTH-1:0] captured_out
);
    parameter WIDTH = 1;

    wire should_capture = current_step == capture_on;
    wire clock_inv = ~clock;

    wire [WIDTH-1:0] dff_out;
    wire do_capture = clock_inv & should_capture;
    dff #(.WIDTH(WIDTH)) dff(
        .clock(do_capture),
        .d(input_value),
        .async_reset(1'h0),
        .enable(1'h1),
        .q(dff_out)
    );

    // Use input as output leading up to capture, then switch to captured value
    assign captured_out = should_capture ^ do_capture ? input_value : dff_out;

endmodule
module arg_loader(
    input clock,
    input [1:0] step,
    input [1:0] capture_on,
    input [15:0] immediate,
    input [15:0] register_value,
    input [15:0] memory_in,
    input is_mem,
    output [15:0] read_address,
    output [15:0] source_value
);

    wire [15:0] register_sum = register_value + immediate;
    assign read_address = register_sum;

    wire [15:0] to_capture = is_mem ? memory_in : register_sum;

    value_capture #(.WIDTH(16)) value_capture (
        .clock(clock),
        .current_step(step),
        .capture_on(capture_on),
        .input_value(to_capture),
        .captured_out(source_value)
    );

endmodule
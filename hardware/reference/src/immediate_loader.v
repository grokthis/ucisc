module immediate_loader (
    input clock,
    input [1:0] step,
    input [15:0] memory_in,
    input [15:0] pc,
    input destination_mem,
    output [15:0] read_address,
    output [15:0] immediate,
    output [3:0] offset,
    output can_halt
);

    assign read_address = pc + 1'h1;

    wire [15:0] captured_immediate;
    wire [15:0] short_immediate = $signed(captured_immediate[10:0]);
    assign immediate = destination_mem ? short_immediate : captured_immediate;

    assign can_halt = captured_immediate == 16'h0;

    assign offset = destination_mem ? captured_immediate[15:12] : 4'h0;

    value_capture #(.WIDTH(16)) value_capture (
        .clock(clock),
        .current_step(step),
        .capture_on(2'h1),
        .input_value(memory_in),
        .captured_out(captured_immediate)
    );

endmodule
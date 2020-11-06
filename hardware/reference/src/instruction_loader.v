module instruction_loader (
    input clock,
    input [1:0] step,
    input [15:0] memory_in,
    input [15:0] pc,
    output [3:0] op,
    output [2:0] effect,
    output [3:0] destination,
    output [3:0] source,
    output push,
    output pop,
    output destination_mem,
    output source_mem,
    output could_halt,
    output [15:0] read_address
);
    wire [15:0] instruction;
    value_capture #(.WIDTH(16)) value_capture(
        .clock(clock),
        .current_step(step),
        .capture_on(2'h0),
        .input_value(memory_in),
        .captured_out(instruction)
    );

    assign read_address = pc;

    assign op = instruction[3:0];
    assign effect = instruction[6:4];
    wire increment = instruction[7];
    assign destination = instruction[11:8];
    assign source = instruction[15:12];

    assign destination_mem = (destination[0] | destination[1]) & ~destination[2];
    assign source_mem = (source[0] | source[1]) & ~source[2];
    assign push = increment & destination_mem;
    assign pop = increment & source_mem & ~destination_mem;
    assign could_halt = (source == 4'h0) & (destination == 4'h0);

endmodule

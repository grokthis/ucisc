module cpu_tb;

    // Setup registers for the inputs
    reg clock;
    reg reset;

    // Setup a wire for the outputs
    wire [1:0] step;
    wire [15:0] pc;
    wire [15:0] r1;
    wire [15:0] mem_read_address;
    wire [15:0] memory_result;
    wire [15:0] source_value;
    wire [15:0] destination_value;
    wire [15:0] result;
    wire [15:0] immediate;
    wire push;
    wire pop;

    // Instantiate the unit under test
    cpu #(.MEM_INIT_FILE("test/test_prog.hex")) cpu(
        .clock_input(clock),
        .reset(reset),
        .step(step),
        .r1_peek(r1),
        .pc_peek(pc),
        .mem_read_address_peek(mem_read_address),
        .memory_result_peek(memory_result),
        .source_value_peek(source_value),
        .destination_value_peek(destination_value),
        .result_peek(result),
        .immediate_peek(immediate),
        .push_peek(push),
        .pop_peek(pop)
    );

    initial begin
        // Initialize standard setup
        clock = 0;

        #20 reset = 1;
        #20 reset = 0;
        #1 `assert(pc, 16'h0000);
        #1 `assert(step, 2'h3);
        #1 `assert(mem_read_address, 16'h0);

        // BEGIN Instruction 1
        #20 clock = 1;
        #1 `assert(memory_result, 16'h4540);
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #1 `assert(memory_result, 16'h0000);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #1 `assert(source_value, 16'h0000);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #1 `assert(result, 16'h0000);
        #20 clock = 1;
        // BEGIN Instruction 2
        #1 `assert(step, 2'h0);
        #1 `assert(memory_result, 16'h41C0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #1 `assert(memory_result, 16'h07FF);
        #20 clock = 0;
        #1 `assert(immediate, 16'h07FF);
        #1 `assert(push, 1'b1);
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #1 `assert(mem_read_address, 16'h0000);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFF);
        // BEGIN Instruction 3
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'h07FF);
        // BEGIN Instruction 4
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        // BEGIN Instruction 5
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        // BEGIN Instruction 6
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFF);
        // BEGIN Instruction 7
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFF);
        // BEGIN Instruction 8
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        // BEGIN Instruction 9
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFF);
        // BEGIN Instruction 10
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'h0000);
        // BEGIN Instruction 11
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFF);
        // BEGIN Instruction 12
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'h0000);
        // BEGIN Instruction 13
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        // BEGIN Instruction 14
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        // BEGIN Instruction 15
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFE);
        // BEGIN Instruction 16
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        // BEGIN Instruction 17
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFE);
        // BEGIN Instruction 18
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'h000A);
        // BEGIN Instruction 19
    end

    initial begin
        $monitor("clock=%d, step=%01x, pc=%04x, r1=%04x, read=%04x, mem=%04x, immediate=%04x, src=%04x, dest=%04x",
        clock, step, pc, r1, mem_read_address, memory_result, immediate, source_value, destination_value);
    end
endmodule
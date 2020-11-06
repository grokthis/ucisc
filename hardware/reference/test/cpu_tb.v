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

    // Instantiate the unit under test
    cpu #(.MEM_INIT_FILE("test/test_prog.hex")) cpu(
        .clock_input(clock),
        .reset(reset),
        .step(step),
        .r1_peek(r1),
        .pc_peek(pc),
        .mem_read_address_peek(mem_read_address),
        .memory_result_peek(memory_result)
    );

    initial begin
        // Initialize standard setup
        clock = 0;

        #20 reset = 1;
        #20 reset = 0;
        #1 `assert(pc, 16'h0000);
        #1 `assert(step, 2'h3);
        #1 `assert(mem_read_address, 16'h0);

        // Setup to load instruction from memory 4540
        #20 clock = 1;
        #1 `assert(memory_result, 16'h4540);
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(memory_result, 16'h000F);
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        #20 clock = 0;
        #20 clock = 1;
        #20 clock = 0;
        #20 clock = 1;
        #20 clock = 0;
        #20 clock = 1;
        #20 clock = 0;
        #20 clock = 1;
        #20 clock = 0;
        #20 clock = 1;

    end

    initial begin
        $monitor("clock=%d, step=%01x, pc=%04x, r1=%04x, read=%04x, mem=%04x",
        clock, step, pc, r1, mem_read_address, memory_result);
    end
endmodule
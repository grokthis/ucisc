module cpu_tb;

    // Setup registers for the inputs
    reg clock;
    reg reset;

    // Setup a wire for the outputs
    wire [1:0] step;
    wire [15:0] pc;
    wire [15:0] r1;

    // Instantiate the unit under test
    cpu #(.MEM_INIT_FILE("test/test_prog.hex")) cpu(
        .clock_input(clock),
        .reset(reset),
        .step(step),
        .r1_peek(r1),
        .pc_peek(pc)
    );

    initial begin
        // Initialize standard setup
        clock = 0;

        #20 reset = 1;
        #20 reset = 0;
        #1 `assert(pc, 16'h0000);

        // BEGIN Instruction 1
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
        // BEGIN Instruction 2
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
        // BEGIN Instruction 4
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'h07FF);
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
        // BEGIN Instruction 7
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFF);
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        // BEGIN Instruction 8
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFF);
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
        // BEGIN Instruction 10
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFF);
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        // BEGIN Instruction 11
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'h0000);
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        // BEGIN Instruction 12
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFF);
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        // BEGIN Instruction 13
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'h0000);
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
        // BEGIN Instruction 16
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFE);
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
        // BEGIN Instruction 18
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'hFFFE);
        #1 `assert(step, 2'h1);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h2);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(step, 2'h3);
        #20 clock = 0;
        #20 clock = 1;
        // BEGIN Instruction 19
        #1 `assert(step, 2'h0);
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(r1, 16'h000A);
    end

    initial begin
        $monitor("clock=%d, step=%01x, pc=%04x, r1=%04x",
        clock, step, pc, r1);
    end
endmodule
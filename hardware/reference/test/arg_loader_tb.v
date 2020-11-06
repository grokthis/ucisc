module arg_loader_tb;

    // Setup registers for the inputs
    reg clock;
    reg [1:0] step;
    reg [1:0] capture_on;
    reg [15:0] immediate;
    reg [15:0] register_value;
    reg [15:0] memory_in;
    reg is_mem;

    // Setup a wire for the outputs
    wire [15:0] read_address;
    wire [15:0] source_value;

    // Instantiate the unit under test
    arg_loader arg_loader(
        .clock(clock),
        .step(step),
        .capture_on(capture_on),
        .immediate(immediate),
        .register_value(register_value),
        .memory_in(memory_in),
        .is_mem(is_mem),
        .read_address(read_address),
        .source_value(source_value)
    );

    initial begin
        // Initialize standard setup
        clock = 0;
        step = 2'h0;
        capture_on = 2'h0;
        immediate = 16'h0;
        register_value = 16'h0;
        memory_in = 4'h0;
        is_mem = 0;

        #1 `assert(read_address, 16'h0);
        // TODO Implement these tests
    end

    initial begin
        $monitor("clock=%04x, step=%x, capture_on=%x, immediate=%04x, register_value=%04x, memory_in=%04x, is_mem=%d, read_address=%04x, source_value=%04x",
        clock, step, capture_on, immediate, register_value, memory_in, is_mem, read_address, source_value);
    end
endmodule
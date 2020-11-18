module memory_block_tb;

    // Setup registers for the inputs
    reg clock;
    reg write_enable;
    reg [15:0] read_address;
    reg [15:0] write_address;
    reg [15:0] data_in;

    // Setup a wire for the output
    wire [15:0] data_out;

    // Instantiate the unit under test
    memory_block #(.WIDTH(14), .MEM_INIT_FILE("test/test_prog.hex")) memory_block(
        .clock(clock),
        .write_enable(write_enable),
        .read_address(read_address),
        .write_address(write_address),
        .data_in(data_in),
        .data_out(data_out)
    );

    initial begin
        // Initialize standard setup
        clock = 0;
        write_enable = 0;
        read_address = 16'h0;
        write_address = 16'h0;
        data_in = 16'h0;

        // Read initial data from file
        #20 clock = 1;
        #1 `assert(data_out, 16'h4540);
        #20 clock = 0;

        // Set output to 0, a known value
        #20 write_enable = 1;
        #20 clock = 1;
        // data written
        #20 clock = 0;
        #20 clock = 1;
        // data read
        #1 `assert(data_out, 16'h0000);

        #20 clock = 0;
        #20 clock = 1;
        #1 write_enable = 0;
        #1 data_in = 16'hAAAA;
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(data_out, 16'h0000);

        // Memory is pseudo dual port, doesn't immediately write through
        #1 data_in = 16'hAAAA;
        #1 write_enable = 1;
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(data_out, 16'h0000);
        #1 write_enable = 0;
        #20 clock = 0;
        // Doesn't read on negative edge
        #1 `assert(data_out, 16'h0000);
        #20 clock = 1;
        // Only reads on positive edge
        #1 `assert(data_out, 16'hAAAA);
    end

    initial begin
        $monitor(
            "clock=%04x, write_enable=%d, read_address=%04x, write_address=%04x, data_in=%04x, data_out=%04x",
            clock, write_enable, read_address, write_address, data_in, data_out);
    end
endmodule
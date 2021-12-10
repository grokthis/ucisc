module memory_block_tb;

    // Setup registers for the inputs
    reg clock;
    reg write_enable;
    reg [15:0] read_address;
    reg [15:0] write_address;
    reg [15:0] data_in;

    // Setup a wire for the output
    wire [15:0] data_out;
    wire [31:0] full_width_data_out;

    // Instantiate the unit under test
    memory_block #(
      .WIDTH(16),
      .MEM_INIT_FILEA("test/test_prog.hex.a.hex"),
      .MEM_INIT_FILEB("test/test_prog.hex.b.hex")
    ) memory_block(
        .clock(clock),
        .write_enable(write_enable),
        .read_address(read_address),
        .write_address(write_address),
        .data_in(data_in),
        .data_out(data_out),
        .full_width_data_out(full_width_data_out)
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
        #1 `assert(full_width_data_out, 32'h45400000);
        #20 clock = 0;
        #1 read_address = 16'h1;
        #20 clock = 1;
        #1 `assert(data_out, 16'h0000);
        #1 `assert(full_width_data_out, 32'h000041c0);
        #20 clock = 0;
        #1 read_address = 16'h3;
        #20 clock = 1;
        #1 `assert(data_out, 16'h07FF);
        #1 `assert(full_width_data_out, 32'h07ff1540);
        #20 clock = 0;

        // Set output to 0, a known value
        #20 write_enable = 1;
        #20 clock = 1;
        // data written
        #20 write_enable = 0;
        #20 clock = 0;
        #1 read_address = 16'h0;
        #20 clock = 1;
        // data read
        #1 `assert(data_out, 16'h0000);
        #1 `assert(full_width_data_out, 32'h00000000);

        // Set output to 0, a known value
        #20 clock = 0;
        #20 write_enable = 1;
        #1 data_in = 16'h1111;
        #1 write_address = 16'h3;
        #20 clock = 1;
        // data written
        #20 write_enable = 0;
        #20 clock = 0;
        #1 read_address = 16'h3;
        #20 clock = 1;
        // data read
        #1 `assert(data_out, 16'h1111);
        #1 `assert(full_width_data_out, 32'h11111540);

        #20 clock = 0;
        #20 clock = 1;
        #1 write_enable = 0;
        #1 data_in = 16'hAAAA;
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(data_out, 16'h1111);
        #1 `assert(full_width_data_out, 32'h11111540);

        // Memory is pseudo dual port, doesn't immediately write through
        #1 data_in = 16'hAAAA;
        #1 write_enable = 1;
        #20 clock = 0;
        #20 clock = 1;
        #1 `assert(data_out, 16'h1111);
        #1 `assert(full_width_data_out, 32'h11111540);
        #1 write_enable = 0;
        #20 clock = 0;
        // Doesn't read on negative edge
        #1 `assert(data_out, 16'h1111);
        #1 `assert(full_width_data_out, 32'h11111540);
        #20 clock = 1;
        // Only reads on positive edge
        #1 `assert(data_out, 16'hAAAA);
        #1 `assert(full_width_data_out, 32'haaaa1540);
        #20 clock = 0;
        #1 read_address = 16'h2;
        #20 clock = 1;
        #1 `assert(data_out, 16'h41c0);
        #1 `assert(full_width_data_out, 32'h41c0aaaa);
    end

    initial begin
        $monitor(
            "clock=%x, write_enable=%d, read_address=%04x, write_address=%04x, data_in=%04x, data_out=%04x, full_width_data_out=%08x",
            clock, write_enable, read_address, write_address, data_in, data_out, full_width_data_out);
    end
endmodule
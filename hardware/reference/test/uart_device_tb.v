module uart_tx_tb;

    // Setup registers for the inputs
    reg clock;
    reg rx;
    reg [3:0] control_address;
    reg control_write;
    reg [15:0] data_in;

    // Setup a wire for the outputs
    wire [7:0] flags;
    wire [15:0] control_read;
    wire tx;

    // Instantiate the unit under test
    uart_device uart_device(
        .clock(clock),
        .Rx(rx),
        .control_address(control_address),
        .control_write(control_write),
        .data_in(data_in),
        .flags(flags),
        .control_read(control_read),
        .Tx(tx)
    );

    initial begin
        // Initialize standard setup
        clock = 0;
        rx = 0;
        control_address = 4'h2;
        control_write = 0;
        data_in = 16'h1;

        #1 `assert(control_read, 16'd12);
        #1 control_write = 1;
        #1 clock = 1;
        #1 `assert(control_read, 16'h1);
        #1 clock = 0;

        #4 control_address = 16'h3; data_in = 16'h00AA;

        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 1);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 1);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 1);
        #1 control_write = 1;
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 1);
        // Begin bit 1 (begin signal)
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        // Begin bit 2 (0)
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        // Begin bit 3 (1)
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        // Begin bit 4 (0)
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        // Begin bit 5 (1)
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        // Begin bit 6 (0)
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        // Begin bit 7 (1)
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        // Begin bit 8 (0)
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        // Begin bit 9 (1)
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        // Begin bit 10 (1 - end bit)
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 0);
        // Tx done, data ready now
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 1);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(flags[0], 1);
        // Never disabled write enable, so picks up next byte
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(flags[0], 0);
    end

    initial begin
        $monitor("clock=%04x, rx=%x, ctrl_addr=%02x, ctrl_write=%h, data_in=%04x, flags=%02x, ctrl_read=%04x, tx=%x",
        clock, rx, control_address, control_write, data_in, flags, control_read, tx);
    end
endmodule
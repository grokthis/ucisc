module uart_tx_tb;

    // Setup registers for the inputs
    reg clock;
    reg [15:0] clock_divider;
    reg [7:0] data_in;
    reg write_en;

    // Setup a wire for the outputs
    wire data_ready;
    wire tx;

    // Instantiate the unit under test
    uart_tx #(.WIDTH(8)) uart_tx(
        .clock(clock),
        .clock_divider(clock_divider),
        .data_in(data_in),
        .write_en(write_en),
        .data_ready(data_ready),
        .tx(tx)
    );

    initial begin
        // Initialize standard setup
        clock = 0;
        clock_divider = 16'h1;
        data_in = 8'hAA;
        write_en = 0;

        #1 `assert(tx, 1);
        #1 `assert(data_ready, 1);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 1);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 1);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 1);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 1);

        #1 write_en = 1;
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 1);
        // Begin bit 1 (begin signal)
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        // Begin bit 2 (0)
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        // Begin bit 3 (1)
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        // Begin bit 4 (0)
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        // Begin bit 5 (1)
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        // Begin bit 6 (0)
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        // Begin bit 7 (1)
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        // Begin bit 8 (0)
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        // Begin bit 9 (1)
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        // Begin bit 10 (1 - end bit)
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 0);
        // Tx done, data ready now
        #4 clock = 1;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 1);
        #4 clock = 0;
        #1 `assert(tx, 1);
        #1 `assert(data_ready, 1);
        // Never disabled write enable, so picks up next byte
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 0;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #1 `assert(tx, 0);
        #1 `assert(data_ready, 0);
    end

    initial begin
        $monitor("clock=%04x, clock_divider=%04x, data_in=%02x, write_en=%x, data_ready=%x, tx=%x",
        clock, clock_divider, data_in, write_en, data_ready, tx);
    end
endmodule
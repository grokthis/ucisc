module uart_tx_tb;

    // Setup registers for the inputs
    reg clock;
    reg [15:0] clock_divider;
    reg read_en;
    reg rx;

    // Setup a wire for the outputs
    wire data_ready;
    wire [7:0] data_out;
    wire [7:0] buffer_out;

    // Instantiate the unit under test
    uart_rx #(.WIDTH(8)) uart_rx(
        .clock(clock),
        .clock_divider(clock_divider),
        .read_en(read_en),
        .rx(rx),
        .data_ready(data_ready),
        .data_out(data_out)
    );

    initial begin
        clock = 0;
        clock_divider = 16'h10;
        read_en = 0;
        rx = 1;

        #1 `assert(data_ready, 0);
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        // Begin signal
        #1 rx = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #1 rx = 1; // bit 1
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #1 rx = 0; // bit 2
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #1 rx = 1; // bit 3
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #1 rx = 0; // bit 4
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #1 rx = 1; // bit 5
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #1 rx = 0; // bit 6
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #1 rx = 1; // bit 7
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #1 rx = 0; // bit 8
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #1 rx = 1; // halt bit
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #1 `assert(data_ready, 1);
        #1 `assert(data_out, 8'h55);
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #4 clock = 1;
        #4 clock = 0;
        #1 `assert(data_ready, 1);
        #1 `assert(data_out, 8'h55);
    end

    initial begin
        $monitor("clock=%04x, clock_divider=%04x, read_en=%x, rx=%x, data_ready=%x, data_out=%02x",
        clock, clock_divider, read_en, rx, data_ready, data_out);
    end
endmodule
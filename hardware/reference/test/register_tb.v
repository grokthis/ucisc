module register_tb;

    // Setup registers for the inputs
    reg [15:0] data_in;
    reg [15:0] increment;
    reg capture;
    reg w_en;
    reg inc_en;
    reg reset;

    // Setup a wire for the output
    wire [15:0] value;

    // Instantiate the unit under test
    register #(.WIDTH(16)) register(
        .data_in(data_in),
        .increment(increment),
        .capture(capture),
        .write_enable(w_en),
        .increment_enable(inc_en),
        .reset(reset),
        .value(value)
    );

    initial begin
        // Initialize standard setup
        data_in = 16'h0;
        increment = 16'h0;
        capture = 0;
        w_en = 1;
        inc_en = 0;
        reset = 0;

        // Set output to 0, a known value
        #20 capture = 1;
        #20 capture = 0;
        #1 `assert(value, 16'h0000);

        // Captures input as output on capture if w_en
        #20 data_in = 16'hFFFF;
        #20 capture = 1;
        #1 `assert(value, 16'hFFFF);

        // If w_en is low at positive capture edge, output doesn't change
        #20 data_in = 16'hEEEE;
        #1 `assert(value, 16'hFFFF);
        #20 w_en = 0;
        #1 capture = 1;
        #1 `assert(value, 16'hFFFF);
        #20 w_en = 1;
        #1 `assert(value, 16'hFFFF);

        // If w_en is low, inc_en supports increment
        #20 increment = 16'h1;
        #20 w_en = 0;
        #20 inc_en = 1;
        #20 `assert(value, 16'hFFFF);
        #20 capture = 0;
        #20 capture = 1;
        #20 capture = 0;
        #20 `assert(value, 16'h0000);
        #20 capture = 1;
        #20 capture = 0;
        #20 `assert(value, 16'h0001);
        #20 capture = 1;
        #1 capture = 0;
        #1 `assert(value, 16'h0002);
        #20 increment = 16'h0002;
        #20 capture = 1;
        #1 capture = 0;
        #1 `assert(value, 16'h0004);
        #20 increment = 16'hFFFF;
        #20 capture = 1;
        #1 capture = 0;
        #1 `assert(value, 16'h0003);
        #20 capture = 1;
        #1 capture = 0;
        #1 `assert(value, 16'h0002);

        // If w_en and inc_en are both asserted, w_en wins
        #20 w_en = 1; inc_en = 1;
        #20 capture = 1;
        #1 capture = 0;
        #1 `assert(value, 16'hEEEE);

        // Reset trumps all
        #20 reset = 1;
        #1 `assert(value, 16'h0000);
        #20 capture = 1;
        #1 capture = 0;
        #1 `assert(value, 16'h0000);
        #20 w_en = 0; inc_en = 1;
        #20 capture = 1;
        #1 capture = 0;
        #1 `assert(value, 16'h0000);
        #20 w_en = 1; inc_en = 0;
        #20 capture = 1;
        #1 capture = 0;
        #1 `assert(value, 16'h0000);
        #20 w_en = 0; inc_en = 0;
        #20 capture = 1;
        #1 capture = 0;
        #1 `assert(value, 16'h0000);
    end

    initial begin
        $monitor(
            "data_in=%04x, increment=%04x, capture=%d, w_en=%d, inc_en=%d reset=%d value=%04x",
            data_in, increment, capture, w_en, inc_en, reset, value);
    end
endmodule
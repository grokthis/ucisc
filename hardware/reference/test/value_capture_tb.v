module value_capture_tb;

    // Setup registers for the inputs
    reg [15:0] input_value;
    reg clock;
    reg [1:0] current_step;
    reg [1:0] capture_on;

    // Setup a wire for the output
    wire [15:0] captured_out;

    // Instantiate the unit under test
    value_capture #(.WIDTH(16)) value_capture(
        .input_value(input_value),
        .clock(clock),
        .current_step(current_step),
        .capture_on(capture_on),
        .captured_out(captured_out)
    );

    initial begin
        // Initialize standard setup
        input_value = 16'h0;
        clock = 0;
        current_step = 2'h0;
        capture_on = 2'h0;

        // Set output to 0, a known value
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'h0000);

        // Should only capture the value on the output on the specified step
        // Step 0
        #20 capture_on = 2'h0;
        #20 input_value = 16'hFFFF;
        #20 current_step = 2'h1;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'h0000);
        #20 current_step = 2'h2;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'h0000);
        #20 current_step = 2'h3;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'h0000);
        #20 current_step = 2'h0;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'hFFFF);

        // Step 1
        #20 capture_on = 2'h1;
        #20 input_value = 16'h0000;
        #20 current_step = 2'h2;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'hFFFF);
        #20 current_step = 2'h3;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'hFFFF);
        #20 current_step = 2'h0;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'hFFFF);
        #20 current_step = 2'h1;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'h0000);


        // Step 2
        #20 capture_on = 2'h2;
        #20 input_value = 16'hFFFF;
        #20 current_step = 2'h3;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'h0000);
        #20 current_step = 2'h0;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'h0000);
        #20 current_step = 2'h1;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'h0000);
        #20 current_step = 2'h2;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'hFFFF);

        // Step 3
        #20 capture_on = 2'h3;
        #20 input_value = 16'h0000;
        #20 current_step = 2'h0;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'hFFFF);
        #20 current_step = 2'h1;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'hFFFF);
        #20 current_step = 2'h2;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'hFFFF);
        #20 current_step = 2'h3;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(captured_out, 16'h0000);

        // Capture happens on negative clock edge.
        // During positive clock, output changes freely
        #20 clock = 1;
        #20 input_value = 16'hEEEE;
        #1 `assert(captured_out, 16'hEEEE);
        #20 input_value = 16'h4242;
        #1 `assert(captured_out, 16'h4242);
        #1 clock = 0;
        #1 input_value = 16'hFFFF;
        #20 `assert(captured_out, 16'h4242);
    end

    initial begin
        $monitor(
            "clock=%04x, input_value=%04x, current_step=%d, capture_on=%d, captured_out=%04x",
            clock, input_value, current_step, capture_on, captured_out);
    end
endmodule
module dff_tb;

    // Setup registers for the DFF inputs
    reg clock;
    reg [15:0] d;
    reg a_reset;
    reg enable;

    // Setup a wire for the dff output "q"
    wire [15:0] q;

    // Instantiate the unit under test
    dff #(.WIDTH(16)) dff(
        .clock(clock),
        .d(d),
        .async_reset(a_reset),
        .enable(enable),
        .q(q)
    );

    initial begin
        // Initialize standard setup
        clock = 0;
        d = 16'h0;
        a_reset = 0;
        enable = 1;

        #20 d = 16'hFFFF;
        // Should capture on positive clock edge
        #20 clock = 1;
        #1 `assert(q, 16'hFFFF);

        // Should not capture without positive clock edge
        #20 d = 16'h1111;
        #1 `assert(q, 16'hFFFF);
        #20 clock = 0;
        #1 `assert(q, 16'hFFFF);

        // Reset sets output and holds output at 0
        #20 a_reset = 1;
        #1 `assert(q, 16'h0000);
        #20 clock = 1;
        #1 `assert(q, 16'h0000);
        #20 a_reset = 0;
        #1 `assert(q, 16'h0000);

        // Enable held at 0 makes DFF unresponsive to clock
        #20 d = 16'hF42F;
        #1 `assert(q, 16'h0000);
        enable = 0;
        #1 `assert(q, 16'h0000);
        #20 clock = 0;
        #1 `assert(q, 16'h0000);
        #20 clock = 1;
        #1 `assert(q, 16'h0000);

        // Enable assertion doesn't capture input
        #20 enable = 1;
        #1 `assert(q, 16'h0000);
        #20 clock = 0;
        #1 `assert(q, 16'h0000);
        // Positive edge captures input
        #20 clock = 1;
        #1 `assert(q, 16'hF42F);
    end

    initial begin
        $monitor("clock=%d, d=%04x, async_reset=%d, enable=%d, q=%04x", clock, d, a_reset, enable, q);
    end
endmodule
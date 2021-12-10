module effect_decoder_tb;

    // Setup registers for the inputs
    reg [15:0] flags;
    reg [2:0] effect;

    // Setup a wire for the outputs
    wire store;

    // Instantiate the unit under test
    effect_decoder effect_decoder(
        .flags(flags),
        .effect(effect),
        .store(store)
    );

    initial begin
        // Initialize standard setup
        flags = 16'h0;
        effect = 3'h0;

        #1 `assert(store, 0);
        // TODO Implement these tests
    end

    initial begin
        $monitor("flags=%04x, effect=%d, store=%x", flags, effect, store);
    end
endmodule
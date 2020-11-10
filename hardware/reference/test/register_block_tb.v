module memory_block_tb;

    // Setup registers for the inputs
    reg clock;
    reg [1:0] step;
    reg reset;
    reg [3:0] desired_source;
    reg [3:0] desired_destination;
    reg [15:0] write_value;
    reg write_enable;
    reg push;
    reg pop;
    reg inc_enable;
    reg [15:0] flags_in;
    reg write_flags;

    // Setup a wire for the output
    wire [15:0] source_value;
    wire [15:0] destination_value;
    wire [15:0] pc;
    wire [15:0] flags;
    wire [15:0] banking;
    wire source_banked;
    wire destination_banked;

    // Instantiate the unit under test
    register_block register_block (
        .clock(clock),
        .step(step),
        .reset(reset),
        .desired_source(desired_source),
        .desired_destination(desired_destination),
        .write_value(write_value),
        .write_enable(write_enable),
        .push(push),
        .pop(pop),
        .inc_enable(inc_enable),
        .flags_in(flags_in),
        .write_flags(write_flags),
        .source_value(source_value),
        .destination_value(destination_value),
        .pc(pc),
        .flags(flags),
        .banking(banking),
        .source_banked(source_banked),
        .destination_banked(destination_banked)
    );

    initial begin
        // Initialize standard setup
        clock = 0;
        step = 2'h1;
        reset = 0;
        desired_source = 4'h0;
        desired_destination = 4'h0;
        write_value = 16'h0;
        write_enable = 1;
        push = 0;
        pop = 0;
        inc_enable = 0;
        flags_in = 16'h0;
        write_flags = 0;

        // Banking should already be initialized
        #1 `assert(banking, 16'h00E0);

        // Set output to 0, a known value
        // Setup the destination on step 1
        #20 clock = 1;
        #20 clock = 0;
        // Write the PC on step 3
        #20 step = 2'h3;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(pc, 16'h0000);
        #1 `assert(source_value, 16'h0000);
        #1 `assert(destination_value, 16'h0000);

        // write to PC
        #20 clock = 1;
        #2 step = 2'h3;
        #2 write_value = 16'hFFFF;
        #2 desired_destination = 4'h0;
        #14 clock = 0;
        #1 `assert(destination_value, 16'hFFFF);
        #1 `assert(pc, 16'hFFFF);
        #20 step = 2'h1;
        #1 clock = 1;
        #1 `assert(destination_value, 16'hFFFF);
        #1 `assert(pc, 16'hFFFF);

        // write to r1 reg, will write
        #20 write_value = 16'hFFF1;
        #2 step = 2'h1; desired_destination = 4'h5;
        #20 clock = 0;
        #20 clock = 1;
        #2 step = 2'h3;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF1);
        // write to r1 mem, doesn`t write
        #2 step = 2'h1; desired_destination = 4'h1;
        #20 clock = 1;
        #20 clock = 0;
        #2 step = 2'h3; write_value = 16'h5555;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF1);

        // write to r2 reg, will write
        #20 write_value = 16'hFFF2;
        #2 step = 2'h1; desired_destination = 4'h6;
        #20 clock = 0;
        #20 clock = 1;
        #2 step = 2'h3;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF2);
        // write to r2 mem, doesn`t write
        #2 step = 2'h1; desired_destination = 4'h2;
        #20 clock = 1;
        #20 clock = 0;
        #2 step = 2'h3; write_value = 16'h5555;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF2);
        #1 `assert(destination_banked, 0);

        // write to r3 reg, will write
        #20 write_value = 16'hFFF3;
        #2 step = 2'h1; desired_destination = 4'h7;
        #20 clock = 0;
        #20 clock = 1;
        #2 step = 2'h3;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF3);
        // write to r3 mem, doesn`t write
        #2 step = 2'h1; desired_destination = 4'h3;
        #20 clock = 1;
        #20 clock = 0;
        #2 step = 2'h3; write_value = 16'h5555;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF3);
        #1 `assert(destination_banked, 0);

        // write to rb1 reg, will write
        #20 write_value = 16'hFFF9;
        #2 step = 2'h1; desired_destination = 4'hD;
        #20 clock = 0;
        #20 clock = 1;
        #2 step = 2'h3;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF9);
        // write to rb1 mem, doesn`t write
        #2 step = 2'h1; desired_destination = 4'h9;
        #20 clock = 1;
        #20 clock = 0;
        #2 step = 2'h3; write_value = 16'h5555;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF9);
        #1 `assert(destination_banked, 1);

        // write to rb2 reg, will write
        #20 write_value = 16'hFFFA;
        #2 step = 2'h1; desired_destination = 4'hE;
        #20 clock = 0;
        #20 clock = 1;
        #2 step = 2'h3;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFFA);
        // write to rb2 mem, doesn`t write
        #2 step = 2'h1; desired_destination = 4'hA;
        #20 clock = 1;
        #20 clock = 0;
        #2 step = 2'h3; write_value = 16'h5555;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFFA);
        #1 `assert(destination_banked, 1);

        // write to rb3 reg, will write
        #20 write_value = 16'hFFFB;
        #2 step = 2'h1; desired_destination = 4'hF;
        #20 clock = 0;
        #20 clock = 1;
        #2 step = 2'h3;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFFB);
        // write to rb2 mem, doesn`t write
        #2 step = 2'h1; desired_destination = 4'hB;
        #20 clock = 1;
        #20 clock = 0;
        #2 step = 2'h3; write_value = 16'h5555;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFFB);
        #1 `assert(destination_banked, 1);

        // Check that registers retained values
        #20 write_enable = 0;
        #2 step = 2'h1; desired_destination = 4'h0; desired_source = 4'h0;
        #20 clock = 1;
        #20 clock = 0;
        // pc should have been incrementing during clock transitions
        #1 `assert(destination_value, 16'h0023);
        #1 `assert(source_value, 16'h0023);

        // Check that registers retained values
        #2 desired_destination = 4'h1; desired_source = 4'h1;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF1);
        #1 `assert(source_value, 16'hFFF1);

        // Check that registers retained values
        #2 desired_destination = 4'h2; desired_source = 4'h2;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF2);
        #1 `assert(source_value, 16'hFFF2);

        // Check that registers retained values
        #2 desired_destination = 4'h3; desired_source = 4'h3;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF3);
        #1 `assert(source_value, 16'hFFF3);

        // Check that registers retained values
        #2 desired_destination = 4'h4; desired_source = 4'h4;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(source_value, 16'h0000);

        // Check that registers retained values
        #2 desired_destination = 4'h5; desired_source = 4'h5;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF1);
        #1 `assert(source_value, 16'hFFF1);

        // Check that registers retained values
        #2 desired_destination = 4'h6; desired_source = 4'h6;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF2);
        #1 `assert(source_value, 16'hFFF2);

        // Check that registers retained values
        #2 desired_destination = 4'h7; desired_source = 4'h7;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF3);
        #1 `assert(source_value, 16'hFFF3);

        // Check that registers retained values
        #2 desired_destination = 4'h9; desired_source = 4'h9;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF9);
        #1 `assert(source_value, 16'hFFF9);

        // Check that registers retained values
        #2 desired_destination = 4'hA; desired_source = 4'hA;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFFA);
        #1 `assert(source_value, 16'hFFFA);

        // Check that registers retained values
        #2 desired_destination = 4'hB; desired_source = 4'hB;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFFB);
        #1 `assert(source_value, 16'hFFFB);

        // Check that registers retained values
        #2 desired_destination = 4'hD; desired_source = 4'hD;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFF9);
        #1 `assert(source_value, 16'hFFF9);

        // Check that registers retained values
        #2 desired_destination = 4'hE; desired_source = 4'hE;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFFA);
        #1 `assert(source_value, 16'hFFFA);

        // Check that registers retained values
        #2 desired_destination = 4'hF; desired_source = 4'hF;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(destination_value, 16'hFFFB);
        #1 `assert(source_value, 16'hFFFB);

        #2 desired_destination = 4'h8; desired_source = 4'h8;
        #20 clock = 1;
        #20 clock = 0;
        #20 flags_in = 16'hFF00; write_enable = 1; write_flags = 1; step = 2'h3;
        #20 clock = 1;
        #20 clock = 0;
        // write_flags wins
        #1 `assert(flags, 16'hFF00);
        #20 flags_in = 16'hFF00; write_enable = 1; write_flags = 0; step = 2'h3;
        #20 clock = 1;
        #20 clock = 0;
        #1 `assert(flags, 16'h5555);

        // TODO test push and pop
    end

    initial begin
        $monitor(
            "clock=%04x, step=%04x, reset=%d, desired_src=%d, desired_dst=%d, write_v=%04x, write_en=%d, push=%d, pop=%d, inc_enable=%d, flags_in=%04x, write_flags=%d, src_val=%04x, dst_val=%04x, pc=%04x, flags=%04x, banking=%04x",
            clock, step, reset, desired_source, desired_destination, write_value, write_enable, push, pop, inc_enable, flags_in, write_flags, source_value, destination_value, pc, flags, banking);
    end
endmodule
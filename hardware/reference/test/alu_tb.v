module alu_tb;

    // Setup registers for the inputs
    reg [15:0] source;
    reg [15:0] destination;
    reg [3:0] op_code;
    reg [15:0] flags;

    // Setup a wire for the outputs
    wire [15:0] result_out;
    wire [15:0] flags_out;
    wire write_flags;

    // Instantiate the unit under test
    alu alu(
        .source(source),
        .destination(destination),
        .op_code(op_code),
        .flags(flags),
        .result_out(result_out),
        .flags_out(flags_out),
        .write_flags(write_flags)
    );

    initial begin
        // Initialize standard setup
        source = 16'h0;
        destination = 16'h0;
        op_code = 4'h0;
        flags = 16'h0100;

        // copy
        #1 `assert(result_out, 16'h0);
        #20 source = 16'hFFFF;
        #1 `assert(result_out, 16'hFFFF);
        #20 destination = 16'h4242;
        #1 `assert(result_out, 16'hFFFF);
        // and
        #20 op_code = 4'h1;
        #1 `assert(result_out, 16'h4242);
        // or
        #20 op_code = 4'h2;
        #1 `assert(result_out, 16'hFFFF);
        // xor
        #20 op_code = 4'h3;
        #1 `assert(result_out, 16'hBDBD);
        // inverse
        #20 op_code = 4'h4;
        #1 `assert(result_out, 16'h0000);
        // shift left
        #20 destination = 16'hFFFF;
        #20 source = 16'h0000;
        #20 op_code = 4'h5;
        #1 `assert(result_out, 16'hFFFF);
        #20 source = 16'h0001;
        #1 `assert(result_out, 16'hFFFE);
        #20 source = 16'h0002;
        #1 `assert(result_out, 16'hFFFC);
        #20 source = 16'h0003;
        #1 `assert(result_out, 16'hFFF8);
        #20 source = 16'h0004;
        #1 `assert(result_out, 16'hFFF0);
        #20 source = 16'h0005;
        #1 `assert(result_out, 16'hFFE0);
        #20 source = 16'h0006;
        #1 `assert(result_out, 16'hFFC0);
        #20 source = 16'h0007;
        #1 `assert(result_out, 16'hFF80);
        #20 source = 16'h0008;
        #1 `assert(result_out, 16'hFF00);
        #20 source = 16'h0009;
        #1 `assert(result_out, 16'hFE00);
        #20 source = 16'h000A;
        #1 `assert(result_out, 16'hFC00);
        #20 source = 16'h000B;
        #1 `assert(result_out, 16'hF800);
        #20 source = 16'h000C;
        #1 `assert(result_out, 16'hF000);
        #20 source = 16'h000D;
        #1 `assert(result_out, 16'hE000);
        #20 source = 16'h000E;
        #1 `assert(result_out, 16'hC000);
        #20 source = 16'h000F;
        #1 `assert(result_out, 16'h8000);
        #20 source = 16'h0010;
        #1 `assert(result_out, 16'h0000);
        #20 source = 16'h1110;
        #1 `assert(result_out, 16'h0000);

        // shift right
        #20 op_code = 4'h6;
        #20 destination = 16'h8000;
        #20 source = 16'h0000;
        #1 `assert(result_out, 16'h8000);
        #20 source = 16'h0001;
        #1 `assert(result_out, 16'hC000);
        #20 source = 16'h0002;
        #1 `assert(result_out, 16'hE000);
        #20 source = 16'h0003;
        #1 `assert(result_out, 16'hF000);
        #20 source = 16'h0004;
        #1 `assert(result_out, 16'hF800);
        #20 source = 16'h0009;
        #1 `assert(result_out, 16'hFFC0);
        #20 source = 16'h000D;
        #1 `assert(result_out, 16'hFFFC);
        #20 source = 16'h000E;
        #1 `assert(result_out, 16'hFFFE);
        #20 source = 16'h000F;
        #1 `assert(result_out, 16'hFFFF);
        #20 source = 16'h0010;
        #1 `assert(result_out, 16'hFFFF);
        #20 source = 16'h1110;
        #1 `assert(result_out, 16'hFFFF);
        #20 source = 16'h111F;
        #1 `assert(result_out, 16'hFFFF);

        // TODO Implement op codes 6+
    end

    initial begin
        $monitor("source=%04x, destination=%04x, op_code=%x, flags=%04x, result_out=%04x, flags_out=%04x, write_flags=%d",
        source, destination, op_code, flags, result_out, flags_out, write_flags);
    end
endmodule
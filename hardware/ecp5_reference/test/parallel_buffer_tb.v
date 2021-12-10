module uart_tx_tb;

    // Setup registers for the inputs
    reg clock;
    reg [7:0] data_in;
    reg write_en;
    reg output_ready;

    // Setup a wire for the outputs
    wire [7:0] data_out;
    wire write_out;
    wire write_ready;

    wire write_p1p2;
    wire ready_p1p2;
    wire [7:0] data_p1p2;

    // Instantiate the unit under test
    parallel_buffer #(.WIDTH(8)) pb1(
        .clock(clock),
        .data_in(data_in),
        .write_in(write_en),
        .next_ready(ready_p1p2),
        .write_ready(write_ready),
        .data_out(data_p1p2),
        .write_out(write_p1p2)
    );
    parallel_buffer #(.WIDTH(8)) pb2(
        .clock(clock),
        .data_in(data_p1p2),
        .write_in(write_p1p2),
        .next_ready(output_ready),
        .write_ready(ready_p1p2),
        .data_out(data_out),
        .write_out(write_out)
    );

    initial begin
        // Initialize standard setup
        clock = 0;
        data_in = 8'h4;
        write_en = 0;
        output_ready = 0;

        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 clock = 1;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 clock = 0;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 clock = 1;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 clock = 0;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 clock = 1;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 clock = 0;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 clock = 1;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 write_en = 1;
        #4 clock = 0;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 clock = 1;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 data_in = 8'h2;
        #4 clock = 0;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 clock = 1;
        #1 `assert(write_ready, 0);
        #1 `assert(write_out, 1);
        #1 `assert(data_out, 8'h4);
        #1 write_en = 0;
        #4 clock = 0;
        #1 `assert(write_ready, 0);
        #1 `assert(write_out, 1);
        #1 `assert(data_out, 8'h4);
        #4 clock = 1;
        #1 `assert(write_ready, 0);
        #1 `assert(write_out, 1);
        #1 `assert(data_out, 8'h4);
        #4 clock = 0;
        #1 `assert(write_ready, 0);
        #1 `assert(write_out, 1);
        #1 `assert(data_out, 8'h4);
        #4 clock = 1;
        #1 `assert(write_ready, 0);
        #1 `assert(write_out, 1);
        #1 `assert(data_out, 8'h4);
        #4 output_ready = 1;
        #4 clock = 0;
        // Write is ready because the chain is ready to shift
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 1);
        #1 `assert(data_out, 8'h4);
        #4 clock = 1;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 1);
        #1 `assert(data_out, 8'h2);
        #4 clock = 0;
        #4 output_ready = 0;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 1);
        #1 `assert(data_out, 8'h2);
        #4 clock = 1;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 1);
        #1 `assert(data_out, 8'h2);
        #4 clock = 0;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 1);
        #1 `assert(data_out, 8'h2);
        #4 clock = 1;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 1);
        #1 `assert(data_out, 8'h2);
        #4 clock = 0;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 1);
        #1 `assert(data_out, 8'h2);
        #4 output_ready = 1;
        #4 clock = 1;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 clock = 0;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 clock = 1;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);
        #4 clock = 0;
        #1 `assert(write_ready, 1);
        #1 `assert(write_out, 0);

    end

    initial begin
        $monitor("clock=%x, data_in=%02x, write_en=%x, write_ready=%x, output_ready=%x, data_out=%02x, write_out=%x",
        clock, data_in, write_en, write_ready, output_ready, data_out, write_out);
    end
endmodule
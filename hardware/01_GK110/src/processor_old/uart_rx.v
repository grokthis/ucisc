module uart_rx(
    input clock,
    input [15:0] clock_divider,
    input read_en,
    input rx,
    output reg data_ready,
    output reg [WIDTH-1:0] data_out
);
    parameter WIDTH = 8;

    initial data_ready = 0;

    reg reading = 0; // Are we currently reading a byte?
    reg [15:0] clock_counter = 0; // Counts up to the clock divider
    reg [3:0] bit_number; // Which bit are we currently reading?
    reg [7:0] buffer = 8'h0; // Shift buffer for incoming bits
    reg buffer_ready = 0; // 1 if buffer has a valid byte value for read

    reg [4:0] debounce_rx = 5'h1F;
    wire real_rx = debounce_rx[0] ?
        (
            debounce_rx[1] ? (debounce_rx[2] | debounce_rx[3] | debounce_rx[4]) :
            debounce_rx[2] ? (debounce_rx[3] | debounce_rx[4]) :
            debounce_rx[3] & debounce_rx[4]
        ) :
        debounce_rx[1] ? (
            debounce_rx[2] ? (debounce_rx[3] | debounce_rx[4]) :
            debounce_rx[3] & debounce_rx[4]
        ) :
        debounce_rx[2] & debounce_rx[3] & debounce_rx[4];

    always @(posedge clock) begin
       debounce_rx <= { debounce_rx[3:0], rx };

       if (~reading) begin
        // Bit number stays at zero while waiting for the next byte
         bit_number <= 4'h0;
         if (~real_rx) begin
            // RX went low, count to half the baud clock width and sample
            clock_counter <= { 1'b0, clock_divider[15:1] };
            reading <= 1'b1;
         end
       end else begin
         if (clock_counter == clock_divider) begin
             // take a sample
             if (bit_number == 4'h0) begin
                 // We are still in the start bit, init details
                 buffer_ready <= 1'b0; // About to write to this
                 bit_number <= bit_number + 1;
                 // Clock counter starts at 1 indexed so you can calc
                 // the divider value based without accounting for zero bit
                 clock_counter <= 16'h1;
             end else if (bit_number == 4'h9) begin
                 // Stop bit, if real_rx is low, we took too long, read right away
                 reading = ~real_rx;
                 if (~real_rx) begin
                    // We took too long on the last byte, this time
                    // Remove one clock cycle and go straight into the
                    // next start bit
                    clock_counter <= clock_divider[15:2];
                    bit_number <= 4'h0;
                 end
             end else begin
                 if (bit_number == 4'h8) begin
                     // We read all 8 bits, byte is ready to move into
                     // the output value
                     buffer_ready <= 1'b1;
                 end
                 clock_counter <= 16'h1;
                 bit_number <= bit_number + 1;
                 buffer = { real_rx, buffer[7:1] };
             end
         end else begin
             clock_counter <= clock_counter + 1;
         end
       end
    end

    // make sure buffer ready goes low at least once between reads
    // to ensure we don't read the same value twice
    reg new_data = 0;
    always @(posedge clock) begin
        if (~data_ready & buffer_ready & new_data) begin
            data_out <= buffer;
            data_ready <= 1'b1;
            new_data <= 1'b0;
        end else if (read_en) begin
            data_ready <= 1'b0;
        end else if (~buffer_ready) begin
            new_data <= 1'b1;
        end
    end
endmodule
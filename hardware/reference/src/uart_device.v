module uart_device(
    input clock,
    input Rx,
    input [3:0] control_address,
    input control_write,
    input [15:0] data_in,
    output [7:0] flags,
    output [15:0] control_read,
    output Tx
);

    assign control_read =
        control_address == 4'h1 ? { flags, 8'h4 } :
        control_address == 4'h2 ? baud_clock_divider : 16'h0;

    wire in_progress = ~write_ready | ~ready_p1p2 | ~tx_ready;
    assign flags = { 4'h1, 1'b0, in_progress, 1'b0, write_ready };

    reg [15:0] baud_clock_divider;
    initial baud_clock_divider = 16'd12;

    wire [7:0] buffer_data_in = data_in[7:0];
    wire buffer_write_en = control_write & control_address == 4'h3;

    always @(posedge clock) begin
        if (control_write & control_address == 4'h2)
            baud_clock_divider <= data_in;
    end

    wire write_p1p2;
    wire ready_p1p2;
    wire [7:0] data_p1p2;
    wire tx_ready;
    wire [7:0] tx_data;
    wire tx_write;

    wire write_ready;
    // Instantiate the unit under test
    parallel_buffer #(.WIDTH(8)) pb1(
        .clock(clock),
        .data_in(buffer_data_in),
        .write_in(buffer_write_en),
        .next_ready(ready_p1p2),
        .write_ready(write_ready),
        .data_out(data_p1p2),
        .write_out(write_p1p2)
    );

    parallel_buffer #(.WIDTH(8)) pb2(
        .clock(clock),
        .data_in(data_p1p2),
        .write_in(write_p1p2),
        .next_ready(tx_ready),
        .write_ready(ready_p1p2),
        .data_out(tx_data),
        .write_out(tx_write)
    );

    uart_tx uart_tx(
        .clock(clock),
        .clock_divider(baud_clock_divider),
        .data_in(tx_data),
        .write_en(tx_write),
        .data_ready(tx_ready),
        .tx(Tx)
    );


    // TODO uart_rx
endmodule
module uart_device(
    input clock,
    input write_enable,
    input control,
    input [7:0] address,
    input [15:0] data_in,
    input rx,
    output [15:0] data_out,
    output tx,
    output [15:0] peek
);

    parameter DEVICE_ID = 16'h0;
    parameter DEVICE_TYPE = 8'h4;

    wire control_write = control & write_enable;
    wire [3:0] control_address = address[3:0];

    wire [15:0] control_read =
        control_address == 4'h0 ? DEVICE_ID :
        control_address == 4'h1 ? { flags, DEVICE_TYPE } :
        control_address == 4'h2 ? baud_clock_divider :
        16'h0;

    assign data_out = control ? control_read : 16'h0;

    wire in_progress = ~write_ready | ~ready_p1p2 | ~tx_ready;
    assign flags = { 4'h1, 1'b0, in_progress, 1'b0, write_ready };

    reg [15:0] baud_clock_divider;
    initial baud_clock_divider = 16'h341; // 9600 baud

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
        .data_in(data_in[7:0]),
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
        .tx(tx)
    );


    // TODO uart_rx
endmodule
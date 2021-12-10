module devices (
    input ref_clock,
    input cpu_clock,
    input write_enable,
    input [15:0] address,
    input [15:0] data_in,
    input rx,
    output reg [15:0] data_out,
    output tx,
    output [15:0] peek,
    output [15:0] pixel_out,
    output h_sync_signal,
    output v_sync_signal
);

    wire control = address[15:12] == 4'h0;
    wire [7:0] target_device = control ? address[11:4] : address[15:8];

    always @(posedge cpu_clock) begin
        data_out <=
            target_device == 8'h02 ? uart_data_out :
            target_device == 8'h10 ? 16'h0 : //graphics_data_out :
            16'h0;
    end

    wire [15:0] uart_data_out;
    uart_device #(.DEVICE_ID(16'h100)) uart_device(
        .clock(cpu_clock),
        .write_enable(write_enable & target_device == 8'h2),
        .control(control),
        .address(address[7:0]),
        .data_in(data_in),
        .rx(rx),
        .data_out(uart_data_out),
        .tx(tx),
        .peek(peek)
    );

//    wire [15:0] graphics_data_out;
//    vga_device vga_device(
//        .ref_clock(ref_clock),
//        .cpu_clock(cpu_clock),
//        .control(control),
////        .write_enable(write_enable & target_device == 8'h10),
//        .write_enable(1'h0),
////        .write_address(address[7:0]),
//        .write_address(8'h0),
////        .data_in(data_in),
//        .data_in(16'h0061),
//        .read_enable(1'h0),
//        .read_address(address[7:0]),
//        .data_out(graphics_data_out),
//        .pixel_out(pixel_out),
//        .h_sync_signal(h_sync_signal),
//        .v_sync_signal(v_sync_signal)
//    );

endmodule

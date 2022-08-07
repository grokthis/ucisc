module devices (
    input ref_clock,
    input cpu_clock,
    input write_enable,
    input [15:0] address,
    input [15:0] data_in,
    input rx,
    output reg [15:0] data_out,
    output tx,
    output [15:0] peek
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

endmodule

module top(
    input CLK12,
    input DIP1,
    input DIP2,
    input DIP3,
    input DIP4,
    input DIP5,
    input DIP6,
    input DIP7,
    input DIP8,
    input SW2,
    output UART_TX,
    input UART_RX,
//    output [7:0] LED,
    output [15:0] PIXEL,
    output H_SYNC,
    output V_SYNC
);

parameter CLOCK_DIV = 10;
parameter CLOCK_MULT = 85;

//assign LED[7] = ~r1[7];
//assign LED[6] = ~r1[6];
//assign LED[5] = ~r1[5];
//assign LED[4] = ~r1[4];
//assign LED[3] = ~r1[3];
//assign LED[2] = ~r1[2];
//assign LED[1] = ~r1[1];
//assign LED[0] = ~r1[0];

assign UART_TX = tx;
wire rx = UART_RX;

//wire [15:0] r1;
wire tx;
//assign LED = tx;

wire ref_clock = CLK12;
wire cpu_clock;
pll #(
  .CLKI_DIV(CLOCK_DIV),
  .CLKFB_DIV(CLOCK_MULT),
  .CLKOP_DIV(1),
)
pll_cpu (
  .clki(CLK12),
  .clko(cpu_clock)
);

cpu #(
    .CLOCK_DIV(10),
    .CLOCK_MULT(105)
) cpu(
    .clk(CLK12),
    .reset(0),
//    .r1_peek(r1),
    .tx(tx),
    .rx(rx)//,
    //.pixel_out(PIXEL),
    //.h_sync_signal(H_SYNC),
    //.v_sync_signal(V_SYNC)
);

wire control = 1'b1;
wire vga_control_we = ~SW2;
wire [15:0] data_input = { 8'h0, DIP8, DIP7, DIP6, DIP5, DIP4, DIP3, DIP2, DIP1 };

    wire [15:0] graphics_data_out;
    vga_device vga_device(
        .ref_clock(ref_clock),
        .cpu_clock(cpu_clock),
        .control(control), //control),
//        .write_enable(write_enable & target_device == 8'h10),
        .write_enable(vga_control_we),
//        .write_address(address[7:0]),
        .write_address(16'h0002),
//        .data_in(data_in),
        .data_in(data_input),
        .read_enable(1'h0),
        .read_address(8'h0), //address[7:0]),
        //.data_out(graphics_data_out),
        .pixel_out(PIXEL), //pixel_out),
        .h_sync_signal(H_SYNC), //h_sync_signal),
        .v_sync_signal(V_SYNC), //v_sync_signal)
    );

endmodule

module vga_device(
  input ref_clock,
  input cpu_clock,
  input control,
  input write_enable,
  input [7:0] write_address,
  input [15:0] data_in,
  input read_enable,
  input [7:0] read_address,
  output [15:0] data_out,
  output [15:0] pixel_out,
  output h_sync_signal,
  output v_sync_signal
);
    parameter DEVICE_ID = 16'h0;
    parameter DEVICE_TYPE = 8'h7;

    reg [15:0] address_page = 16'h0;
    reg [15:0] graphics_mode = 16'h0054;

    // Aspect ratio 1 = 16x9, 0 = 4x3
    reg aspect_ratio = 1'h1; //graphics_mode[3];
    // Text mode 1 = text, 0 = graphics
    reg text_mode = 1'h1; //graphics_mode[2];
    // Resolution 0 = lowest, 3 = highest
    reg [1:0] resolution = graphics_mode[1:0]; //2'h0; //graphics_mode[1:0];

    wire [15:0] control_read =
        read_address[3:0] == 4'h0 ? DEVICE_ID :
        read_address[3:0] == 4'h1 ? { 8'h13, DEVICE_TYPE } :
        read_address[3:0] == 4'h2 ? graphics_mode :
        read_address[3:0] == 4'h3 ? address_page :
        16'h0;

    wire control_write_enable = control & write_enable;
    wire mem_write_enable = ~control & write_enable;
    always @(posedge cpu_clock) begin
        if (control_write_enable & write_address[3:0] == 4'h2)
            graphics_mode <= data_in;
        if (control_write_enable & write_address[3:0] == 4'h3)
            address_page <= data_in;
    end

    wire [15:0] graphics_mem_out;
    assign data_out = control ? control_read : graphics_mem_out;

    wire clk_1280_960;
    pll #(
        .CLKI_DIV(10),
        .CLKFB_DIV(85),
        .CLKOP_DIV(2),
      )
      pll_1280_960_pixel (
        .clki(ref_clock),
        .clko(clk_1280_960)
      );

    wire clk_1920_1080;
    pll #(
        .CLKI_DIV(10),
        .CLKFB_DIV(124),
        .CLKOP_DIV(1),
      )
      pll_1920_1080_pixel (
        .clki(ref_clock),
        .clko(clk_1920_1080)
      );

    vga vga(
      .clk_1920_1080(clk_1920_1080),
      .clk_1280_960(clk_1280_960),
      .pixel_in(pixel_data),
      .aspect_ratio(aspect_ratio),
      .text_mode(text_mode),
      .resolution(resolution),
      .h_sync_signal(h_sync_signal),
      .v_sync_signal(v_sync_signal),
      .pixel_change_clk(pixel_change_clk),
      .pixel_out(pixel_out),
      .v_sync_en(v_sync_en),
      .pixel_width(pixel_width),
      .screen_width(screen_width),
      .screen_height(screen_height)
    );

    wire pixel_change_clk;
    wire [2:0] pixel_width;
    wire v_sync_en;
    reg [4:0] char_pixel_clk = 5'h0;
    wire [20:0] pixel_read;
    wire [15:0] pixel_data;
    wire [10:0] screen_width;
    wire [10:0] screen_height;

    wire [23:0] graphics_write_addr = { address_page, write_address };
    wire [23:0] graphics_read_addr = { address_page, read_address };
    graphics_memory graphics_memory(
      .pixel_clk(pixel_change_clk),
      .cpu_clk(cpu_clock),
//      .pixel_width(pixel_width),
      .text_mode_in(text_mode),
      .aspect_ratio_in(aspect_ratio),
      .resolution_in(resolution),
//      .p_width(screen_width),
//      .p_height(screen_height),
      .write_data(data_in),
      .write_address(graphics_write_addr[16:0]),
      .write_enable(mem_write_enable),
      .read_address(graphics_read_addr[16:0]),
      .read_enable(read_enable),
      .data_out(graphics_mem_out),
      .pixel_out(pixel_data),
      .v_sync_en(v_sync_en)
    );

endmodule

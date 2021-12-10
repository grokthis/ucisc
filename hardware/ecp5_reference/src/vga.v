module vga(
  input clk_1920_1080,
  input clk_1280_960,
  input [15:0] pixel_in,
  input aspect_ratio,
  input text_mode,
  input [1:0] resolution,
  output reg h_sync_signal,
  output reg v_sync_signal,
  output pixel_change_clk,
  output reg [15:0] pixel_out,
  output v_sync_en,
  output [10:0] screen_width,
  output [10:0] screen_height,
  output [2:0] pixel_width
);

wire pixel_clk = aspect_ratio ? clk_1920_1080 : clk_1280_960;

assign screen_width = aspect_ratio ? 11'd1920 : 11'd1280;
wire [15:0] h_fp = aspect_ratio ? 16'd2008 : 16'd1360;
wire [15:0] h_sync = aspect_ratio ? 16'd2052 : 16'd1496;
wire [15:0] h_bp = aspect_ratio ? 16'd2200 : 16'd1712;

assign screen_height = aspect_ratio ? 16'd1080 : 16'd960;
wire [15:0] v_fp = aspect_ratio ? 16'd1084 : 16'd961;
wire [15:0] v_sync = aspect_ratio ? 16'd1089 : 16'd964;
wire [15:0] v_bp = aspect_ratio ? 16'd1125 : 16'd994;

wire h_polarity = aspect_ratio ? 1'b1 : 1'b0;
wire v_polarity = 1'b1;

assign pixel_width =
  resolution == 2'h0 ? 3'h3 :
  resolution == 2'h1 ? 3'h2 :
  resolution == 2'h2 ? 3'h1 :
  3'h0;

reg [15:0] h_count = 16'h0;
reg [15:0] v_count = 16'h0;

wire h_pixel_en = h_count < screen_width;
wire h_sync_en = (h_count >= h_fp) && (h_count < h_sync);

wire v_pixel_en = v_count < screen_height;
assign v_sync_en = (v_count >= v_fp) && (v_count < v_sync);

wire pixel_out_enabled = (h_count < screen_width && v_count < screen_height);
//always @(negedge pixel_clk) begin
//  pixel_out_enabled <= (h_count < screen_width && v_count < screen_height);
//end

wire [15:0] h_next = h_count + 1'b1;
wire [15:0] v_next = v_count + 1'b1;

reg pixel_clk_enabled;
assign pixel_change_clk = pixel_clk_enabled & ~pixel_clk;

always @(negedge pixel_clk) begin
  pixel_clk_enabled <= h_count < screen_width && v_count < screen_height;
  if (h_next == h_bp) begin
    h_count <= 16'h0;

    if (v_next == v_bp)
      v_count <= 16'h0;
    else begin
      v_count <= v_next;
    end
  end else begin
    h_count <= h_next;
  end

  h_sync_signal <= h_polarity ^ ~h_sync_en;
  v_sync_signal <= v_polarity ^ ~v_sync_en;

  pixel_out <= (h_pixel_en & v_pixel_en) ? pixel_in : 16'h0000;
end

endmodule

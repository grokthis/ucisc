module graphics_modes(
  input cpu_clk,
  input v_sync_en,
  input [15:0] graphics_mode,
  output reg aspect_ratio,
  output reg text_mode,
  output reg [1:0] resolution_mode,
  output reg [2:0] color_mode,
  output reg color16_active,
  output reg color8_active,
  output reg color4_active,
  output reg color2_active,
  output reg color1_active,
  output reg color0_active,
  output reg [10:0] screen_width,
  output reg [10:0] screen_height,
  output reg [7:0] row_width,
  output reg [7:0] row_height,
  output reg [2:0] pixel_width,
  output reg [2:0] pixel_divider,
  output reg [15:0] total_chars,
  output reg x1_slices,
  output reg x15_slices,
  output reg x2_slices
);

parameter SLICE_SIZE = 20'h4000; // 16k
parameter SLICE_WIDTH = 14;
parameter ADDRESS_WIDTH = 17;

wire [2:0] pixel_width_next =
    graphics_mode[1:0] == 2'h0 ? 3'h3 :
    graphics_mode[1:0] == 2'h1 ? 3'h2 :
    graphics_mode[1:0] == 2'h2 ? 3'h1 :
    3'h0;
wire text_mode_next = graphics_mode[2];
wire [1:0] resolution_mode_next = graphics_mode[1:0];
wire aspect_ratio_next = graphics_mode[3];

wire color16_ok;
wire color8_ok;
wire color4_ok;
wire color2_ok;
wire color1_ok;
// half-bit color allowed for all modes
wire color0 = 1'b1;


// Color modes indicate the color depth that can be used in a given mode
// 15-bit color for all text modes except 16x9 high-res
// 15-bit color for 16x9 low-res, 4x3 lower-res
wire color16_ok = text_mode_next & !(aspect_ratio_next & resolution_mode_next == 2'h3) |
  resolution_mode_next == 1'h0 |
  (~aspect_ratio_next & resolution_mode_next == 2'h1);

// 7-bit color for all text modes, all lower-res graphics
wire color8_ok = text_mode_next | resolution_mode_next <= 2'h1;

// 3-bit color for all text modes, all lower-res graphics, 4x3 higher-res
wire color4_ok = text_mode_next | resolution_mode_next <= 2'h1 |
  (~aspect_ratio_next & resolution_mode_next == 2'h2);

// 1-bit color for all text modes, all higher-res graphics
wire color2_ok = text_mode_next | resolution_mode_next <= 2'h3;

// 0-bit color for all text modes, all higher-res graphics, 4x3 high-res
wire color1_ok = text_mode_next | resolution_mode_next <= 2'h3 |
  (~aspect_ratio_next & resolution_mode_next == 2'h3);

wire [10:0] screen_width_next = aspect_ratio_next ? 11'd1920 : 11'd1280;
wire [10:0] screen_height_next = aspect_ratio_next ? 16'd1080 : 16'd960;

always @(posedge cpu_clk & v_sync_en) begin
  // Graphics mode breakdown
  // color_mode 0 = 0.5 bit, 1 = 1 bit, 2 = 2 bit, 3 = 4 bit, 4 = 8 bit, 5 = 16 bit
  color_mode <= graphics_mode[6:4];
  // Aspect ratio 1 = 16x9, 0 = 4x3
  aspect_ratio <= aspect_ratio_next;
  // Text mode 1 = text, 0 = graphics
  text_mode <= text_mode_next;
  // Resolution 0 = lowest, 3 = highest
  resolution_mode <= resolution_mode_next;

  pixel_width <= pixel_width_next;

  screen_width <= screen_width_next;
  screen_height <= screen_height_next;

  row_width <= screen_width_next >> (pixel_width_next + 3'h3);
  row_height <= ((screen_height_next >> pixel_width_next) + 3'h7) >> 3'h3;

  pixel_divider <= (3'h1 << pixel_width_next) - 1'b1;
  total_chars <= row_width * row_height;

  // All text modes except 16x9 and 4x3 high res fit in 1 slice
  x1_slices <= text_mode_next && (resolution_mode_next != 2'h3);

  // 4x3 high res text (160x120) fit in 1.5 slices
  x15_slices <= text_mode_next && ~aspect_ratio_next && resolution_mode_next == 2'h3;

  // 16x9 high res text (240x135) fit in 2 slices
  x2_slices <= text_mode_next && aspect_ratio_next && resolution_mode_next == 2'h3;

  color16_active <= color16_ok & color_mode == 3'h5;
  color8_active <= (color8_ok & color_mode == 3'h4) |
      (~color16_ok & color_mode > 3'h4);
  color4_active <= (color4_ok & color_mode == 3'h3) |
      (~color16_ok & ~color8_ok & color_mode > 3'h3);
  color2_active <= (color2_ok & color_mode == 3'h2) |
      (~color16_ok & ~color8_ok & ~color4_ok & color_mode > 3'h2);
  color1_active <= (color1_ok & color_mode == 3'h1) |
      (~color16_ok & ~color8_ok & ~color4_ok & ~color2_ok & color_mode > 3'h1);
  color0_active <= color_mode == 3'h0 |
      (~color16_ok & ~color8_ok & ~color4_ok & ~color2_ok & ~color1_ok);
end

endmodule

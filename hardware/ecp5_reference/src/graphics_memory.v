module graphics_memory(
  input pixel_clk,
  input cpu_clk,
  input text_mode_in,
  input aspect_ratio_in,
  input [1:0] resolution_in,
  input [15:0] write_data,
  input [ADDRESS_WIDTH-1:0] write_address,
  input write_enable,
  input [ADDRESS_WIDTH-1:0] read_address,
  input read_enable,
  input v_sync_en,
  output reg [15:0] data_out,
  output reg [15:0] pixel_out
);

//wire [15:0] graphics_mode = { 12'h5, aspect_ratio_in, text_mode_in, resolution_in };
wire [15:0] graphics_mode = { 12'h5, aspect_ratio_in, text_mode_in, 2'h2 };

parameter SLICE_SIZE = 20'h4000; // 16k
parameter SLICE_WIDTH = 14;
parameter ADDRESS_WIDTH = 17;
parameter MEM_INIT_FILE0 = "video0.hex";
parameter MEM_INIT_FILE1 = "video1.hex";
parameter MEM_INIT_FILE2 = "video2.hex";
parameter MEM_INIT_FILE3 = "video3.hex";
parameter MEM_INIT_FILE4 = "video4.hex";
parameter MEM_INIT_FILE5 = "video5.hex";

reg [15:0] mem0[SLICE_SIZE-1:0];
reg [15:0] mem1[SLICE_SIZE-1:0];
reg [15:0] mem2[SLICE_SIZE-1:0];
reg [15:0] mem3[SLICE_SIZE-1:0];
reg [15:0] mem4[SLICE_SIZE-1:0];
reg [15:0] mem5[SLICE_SIZE-1:0];

initial begin
  if (MEM_INIT_FILE0 != "") begin
    $readmemh(MEM_INIT_FILE0, mem0);
  end
  if (MEM_INIT_FILE1 != "") begin
    $readmemh(MEM_INIT_FILE1, mem1);
  end
  if (MEM_INIT_FILE2 != "") begin
    $readmemh(MEM_INIT_FILE2, mem2);
  end
  if (MEM_INIT_FILE3 != "") begin
    $readmemh(MEM_INIT_FILE3, mem3);
  end
  if (MEM_INIT_FILE4 != "") begin
    $readmemh(MEM_INIT_FILE4, mem4);
  end
  if (MEM_INIT_FILE5 != "") begin
    $readmemh(MEM_INIT_FILE5, mem5);
  end
end

wire x1_slices;
wire x15_slices;
wire x2_slices;
wire color16_active;
wire color8_active;
wire color4_active;
wire color2_active;
wire color1_active;

wire [7:0] row_width;
wire [7:0] row_height;
wire [2:0] pixel_divider;
wire [15:0] total_chars;
wire [10:0] screen_height;
wire text_mode;

graphics_modes graphics_modes(
  .cpu_clk(cpu_clk),
  .v_sync_en(v_sync_en),
  .graphics_mode(graphics_mode),
  .text_mode(text_mode),
  .color16_active(color16_active),
  .color8_active(color8_active),
  .color4_active(color4_active),
  .color2_active(color2_active),
  .color1_active(color1_active),
  .color0_active(color0_active),
  .screen_height(screen_height),
  .row_width(row_width),
  .row_height(row_height),
  .pixel_divider(pixel_divider),
  .total_chars(total_chars),
  .x1_slices(x1_slices),
  .x15_slices(x15_slices),
  .x2_slices(x2_slices)
);

reg [15:0] write_data0;
reg [15:0] write_data1;
reg [15:0] write_data2;
reg [15:0] write_data3;
reg [15:0] write_data4;
reg [15:0] write_data5;
reg [ADDRESS_WIDTH-1:0] write_address_b;
reg write_enable_b;
reg [ADDRESS_WIDTH-1:0] read_address_b;
reg read_enable_b;

always @(posedge cpu_clk) begin
  write_data0 <= write_data;
  write_data1 <= write_data;
  write_data2 <= write_data;
  write_data3 <= write_data;
  write_data4 <= write_data;
  write_data5 <= write_data;
  write_address_b <= write_address;
  write_enable_b <= write_enable;
  read_address_b <= read_address;
  read_enable_b <= read_enable;
end

reg [20:0] pixel_offset;
reg [2:0] bit_offset;

wire [15:0] max_row_offset = row_height - 1'b1;
wire [15:0] max_col_offset = row_width - 1'b1;
reg [15:0] char_row_offset;
reg [15:0] char_col_offset;
reg [2:0] pixel_y_count;
reg [2:0] screen_y_count;
reg [2:0] pixel_x_count;
reg [2:0] char_x_offset;
reg [2:0] char_y_offset;
reg [20:0] next_char_offset;
reg [15:0] template = 16'h467F;

//wire prefetch = ~pixel_offset[2];

always @(posedge (pixel_clk | v_sync_en)) begin
  if (v_sync_en) begin
    if (next_char_offset == 16'h0) begin
        template <= template_prefetch;
        fg_template <= x2_slices ? mem2_prefetch : mem1_prefetch;
        bg_template <= x2_slices ? mem3_prefetch : mem2_prefetch;
    end

    pixel_offset <= 21'h0;
    char_row_offset <= 16'h0;
    char_col_offset <= 16'h0;
    pixel_x_count <= 3'h0;
    pixel_y_count <= 3'h0;
    screen_y_count <= 10'h1;
    char_x_offset <= 3'h0;
    char_y_offset <= 3'h0;
    next_char_offset <= 16'h1;
  end else if (pixel_x_count == pixel_divider) begin
    pixel_x_count <= 3'h0;

    if (char_x_offset == 3'h7) begin
      char_x_offset <= 3'h0;

      // Changing characters, load the prefetch values into templates
      template <= template_prefetch;
      fg_template <=
          x1_slices ? mem1_prefetch :
          x15_slices ? (
              next_char_offset < (SLICE_SIZE >> 1) ? mem1_prefetch : mem2_prefetch
          ) :
          x2_slices ? mem2_prefetch : 16'hFFFF;

      bg_template <=
          x1_slices ? mem2_prefetch :
          x15_slices ? (
              next_char_offset < SLICE_SIZE ? mem3_prefetch : mem4_prefetch
          ) :
          x2_slices ? mem3_prefetch : 16'hFFFF;

      if (char_col_offset == max_col_offset) begin
        next_char_offset <= next_char_offset + 1'b1;
        char_col_offset <= 3'h0;
        if (pixel_y_count == pixel_divider) begin
          pixel_y_count <= 3'h0;
          screen_y_count <= screen_y_count + 1'b1;

          if (char_y_offset == 3'h7) begin
            char_y_offset <= 3'h0;
            if (char_row_offset == max_row_offset) begin
              char_row_offset <= 16'h0;
            end else begin
              char_row_offset <= char_row_offset + 1'b1;
            end
          end else begin
            char_y_offset <= char_y_offset + 1'b1;
          end
        end else begin
          pixel_y_count <= pixel_y_count + 1'b1;
          screen_y_count <= screen_y_count + 1'b1;
        end
      end else begin
        char_col_offset <= char_col_offset + 1'b1;
        if (char_col_offset + 1'b1 == max_col_offset) begin
          if (pixel_y_count == pixel_divider && char_y_offset == 3'h7) begin
            if (screen_y_count == screen_height) begin
              next_char_offset <= 16'h0;
            end else begin
              next_char_offset <= next_char_offset + 1'b1;
            end
          end else if (pixel_y_count == pixel_divider && screen_y_count == screen_height) begin
            next_char_offset <= 16'h0;
          end else begin
            next_char_offset <= next_char_offset - row_width + 1'b1;
          end
        end else begin
          next_char_offset <= next_char_offset + 1'b1;
        end
      end
    end else begin
      char_x_offset <= char_x_offset + 1'b1;
    end
  end else begin
    pixel_x_count <= pixel_x_count + 1'b1;
  end
end

reg [15:0] fg_template;
reg [15:0] bg_template;
//wire [15:0] fg_prefetch = 16'h06AA;
//wire [15:0] bg_prefetch = 16'h06AA;
wire [15:0] template_prefetch = mem5_prefetch;

wire [15:0] fg_value = fg_template;//(char_col_offset[4:0] << 11) + 5'hF;//fg_template;
wire [15:0] bg_value = bg_template;//(char_row_offset[6:1] << 5);//bg_template;
reg [15:0] template_value;

reg [2:0] pixel_y_count_cache;

reg [15:0] char_prefetch;
reg [15:0] mem0_prefetch;
reg [15:0] mem1_prefetch;
reg [15:0] mem2_prefetch;
reg [15:0] mem3_prefetch;
reg [15:0] mem4_prefetch;
reg [15:0] mem5_prefetch;
reg [15:0] graphics_prefetch;

wire [2:0] next_y_offset = char_y_offset + 1'b1;
wire [1:0] next_template_y_offset = char_col_offset == max_col_offset
    && pixel_y_count == pixel_divider ? next_y_offset[2:1] : char_y_offset[2:1];
wire [15:0] template_offset = { char_prefetch[13:0], next_template_y_offset };
wire [15:0] x15_char_offset = x15_slices ? 15'h2000 + next_char_offset : next_char_offset;

// In graphics mode, next_char_offset will be the pixel offset
wire [16:0] graphics_pixel =
    color16_active ? next_char_offset[16:0] :
    color8_active ? {1'h0,next_char_offset[17:1]} :
    color4_active ? {2'h0,next_char_offset[18:2]} :
    color2_active ? {3'h0,next_char_offset[19:3]} :
    color1_active ? {4'h0,next_char_offset[20:4]} :
    color0_active ? {5'h0,next_char_offset[20:5]} :
    16'hFFFF;
wire [SLICE_WIDTH-1:0] mem0_offset = text_mode ? next_char_offset : graphics_pixel;
wire [SLICE_WIDTH-1:0] mem1_offset = text_mode ? x15_char_offset : graphics_pixel;
wire [SLICE_WIDTH-1:0] mem2_offset = text_mode ? x15_char_offset : graphics_pixel;
wire [SLICE_WIDTH-1:0] mem3_offset = text_mode ? (x1_slices ? template_offset : next_char_offset) : graphics_pixel;
wire [SLICE_WIDTH-1:0] mem4_offset = text_mode ? (x2_slices ? template_offset : x15_char_offset) : graphics_pixel;
wire [SLICE_WIDTH-1:0] mem5_offset = text_mode ? template_offset : graphics_pixel;

reg [16:0] graphics_pixel_offset;
always @(posedge (pixel_clk | v_sync_en)) begin
  mem0_prefetch <= mem0[mem0_offset];
  mem1_prefetch <= mem1[mem1_offset];
  mem2_prefetch <= mem2[mem2_offset];
  mem3_prefetch <= mem3[mem3_offset];
  mem4_prefetch <= mem4[mem4_offset];
  mem5_prefetch <= mem5[mem5_offset];
  char_prefetch <= next_char_offset < SLICE_SIZE ? mem0_prefetch : mem1_prefetch;

  graphics_pixel_offset <= graphics_pixel;

  data_out <= 16'h0;
end

// In text mode, if x1_slices:
//   mem0: text character code
//   mem1: foreground color
//   mem2: background color
//   mem3: character set (8k-12k)
//   mem4: character set (4k-8k)
//   mem5: character set (0-4k)
// In text mode, if x15_slices:
//   mem0: text character code
//   mem1: text character code, foreground color
//   mem2: foreground color
//   mem3: background color
//   mem4: background color, second half unused
//   mem5: character set (0-4k)
// In text mode, if x2_slices:
//   mem0: text character code
//   mem1: text character code
//   mem2: foreground color
//   mem3: background color
//   mem4: character set (4-8k)
//   mem5: character set (0-4k)
always @(posedge pixel_clk) begin
  graphics_prefetch <=
      graphics_pixel_offset < SLICE_SIZE ? mem0_prefetch :
      graphics_pixel_offset < {SLICE_SIZE,1'h0} ? mem1_prefetch :
      graphics_pixel_offset < {SLICE_SIZE,2'h0} ? mem2_prefetch :
      graphics_pixel_offset < {SLICE_SIZE,4'h0} ? mem3_prefetch :
      graphics_pixel_offset < {SLICE_SIZE,5'h0} ? mem4_prefetch :
      mem5_prefetch;
end

//wire slice0_write = write_address_b[ADDRESS_WIDTH-1:14] == 3'h0;
//wire slice1_write = write_address_b[ADDRESS_WIDTH-1:14] == 3'h1;
//wire slice2_write = write_address_b[ADDRESS_WIDTH-1:14] == 3'h2;
//wire slice3_write = write_address_b[ADDRESS_WIDTH-1:14] == 3'h3;
//wire slice4_write = write_address_b[ADDRESS_WIDTH-1:14] == 3'h4;
//wire slice5_write = write_address_b[ADDRESS_WIDTH-1:14] == 3'h5;
//
//always @(posedge pixel_clk) begin
//  if (write_enable_b && slice0_write)
//      mem0[write_address_b[SLICE_WIDTH-1:0]] <= write_data0;
//  if (write_enable_b && slice1_write)
//      mem1[write_address_b[SLICE_WIDTH-1:0]] <= write_data1;
//  if (write_enable_b && slice2_write)
//      mem2[write_address_b[SLICE_WIDTH-1:0]] <= write_data2;
//  if (write_enable_b && slice3_write)
//      mem3[write_address_b[SLICE_WIDTH-1:0]] <= write_data3;
//  if (write_enable_b && slice4_write)
//      mem4[write_address_b[SLICE_WIDTH-1:0]] <= write_data4;
//  if (write_enable_b && slice5_write)
//      mem5[write_address_b[SLICE_WIDTH-1:0]] <= write_data5;
//end

wire [7:0] upper_template = template[15:8];
wire [7:0] lower_template = template[7:0];
wire [7:0] upper_char_sel = upper_template >> char_x_offset;
wire [7:0] lower_char_sel = lower_template >> char_x_offset;
wire fg_select = char_y_offset[0] ? lower_char_sel[0] : upper_char_sel[0];

wire [15:0] fg_8_value = {fg_value[15:8],fg_value[8],fg_value[8],6'h0};
wire [15:0] bg_8_value = {fg_value[7:0],fg_value[0],fg_value[0],6'h0};


always @(posedge pixel_clk) begin
  pixel_out <= (fg_select ? fg_value : bg_value);
end

//  color16_active ? (fg_select ? fg_value : bg_value) :
//  color8_active ? (fg_select ? fg_8_value : bg_8_value) :
//  color4_active ? 16'hFFFF :
//  color2_active ? 16'hF800 :
//  16'h07E0;

endmodule

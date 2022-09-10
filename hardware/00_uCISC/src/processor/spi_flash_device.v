module spi_flash_device(
    input cpu_clock,
    input write_enable,
    input is_control,
    input [7:0] short_address,
    input [15:0] cpu_data_in,
    output reg [15:0] cpu_data_out,

    /*************************************
     * Support a quad read mode:
     *  - MOSI and MISO0 will be the same pin
     *  - If MOSI_enabled = 0, MISO0 will read from
     *    device, will read MOSI otherwise
     *  - MISO1 is the same as the traditional MISO pin
     *  - MISO0-MISO3 are only valid inputs when the
     *    appropriate read mode is active.
     *  - CSLow_out is active low. Pulling high terminates I/O.
     *************************************/
    output reg MOSI,
    output MOSI_enable,
    input MISO0,
    input MISO1,
    input MISO2,
    input MISO3,
    output SCK_out,
    output CSLow_out
);
  parameter DEVICE_ID = 8'h40;
  parameter DEVICE_TYPE = 8'h3; // Block I/O device
  parameter CPU_FREQ = 10000000;

  // Block 0 is mapped to this offset on the actual SPI hardware
  // Useful if only part of the SPI flash is accessible to the CPU
  parameter ADDRESS_OFFSET = 32'h0;
  // Maximum block address. The max SPI block is
  // MAX_BLOCK_ADDRESS + ADDRESS_OFFSET
  parameter MAX_BLOCK_ADDRESS = 32'hFFFFFFFF;

  parameter RESUME_COMMAND = 8'hAB;
  parameter POWER_DOWN = 8'hB9;
  parameter READ_COMMAND = 8'h03;
  parameter FAST_READ_COMMAND = 8'h0B;
  parameter FAST_DUAL_READ_COMMAND = 8'h3B;
  parameter FAST_QUAD_READ_COMMAND = 8'h6B;

  /********************************************************************
   * Determines the supported read modes:
   *  - 0 -> normal read
   *  - 1 -> fast read, normal read
   *  - 2 -> fast dual read, fast read, normal read
   *  - 3 -> fast quad read, fast dual read, fast read, normal read
   ********************************************************************/
  parameter SUPPORTED_READ_MODES = 2'h3;
  // TODO: support write mode
  parameter WRITE_MODE_SUPPORTED = 1'h0;

  /********************************************************************
   * States:
   *  - 0 -> idle
   *  - 1 -> sending command
   *  - 2 -> sending address
   *  - 3 -> reading data
   ********************************************************************/
  reg [1:0] state = 2'h0;
  // The number of bits remaining to transfer (read or write)
  reg [12:0] remaining_bits = 13'h0;
  reg [31:0] byte_buffer = 32'h0;

  /************************
   * SPI interface logic *
   ************************/
  always @(negedge cpu_clock) begin
    // MOSI_enable when sending command or address
    MOSI <= byte_buffer[31];
  end

  assign MOSI_enable = state != 2'h3;
  assign CSLow_out = state == 2'h0;
  assign SCK_out = state == 2'h0 | initiate | cpu_clock;

  /********************************************************************
   * flags: 00MM WRDA
   *  - A: Active, 1 if data transfer is active
   *  - D: dirty bit, 1 if data block has changed since last block read
   *  - R: read enabled, updates to the address will trigger data load
   *  - W: write enabled, updates to the address will trigger data write (if dirty)
   *  - M: read mode (see SUPPORTED_READ_MODES, can only set to maximum of supported)
   *
   * Control buffer (words):
   *  - 0: DEVICE_ID
   *  - 1: flags | DEVICE_TYPE
   *  - 2: 16-bit lower block address
   *  - 3: 16-bit upper block address (if supported; MAX_BLOCK_ADDRESS > 16'hFFFF)
   *  - 4: active (MSB) | 0 | transfer_state (next 2 bits) | remaining_bits - 1 (least sig 12 bits)
   ********************************************************************/
  reg [5:0] flags_reg = 6'h0C;
  wire [7:0] flags = { flags_reg, 1'h0, initiate || state != 2'h0 };
  reg [15:0] block_address = 16'h0;
  wire [1:0] read_mode = flags_reg[3:2];
  reg enabled = 1'h0;

  wire control_write = is_control & write_enable;
  wire [3:0] control_address = short_address[3:0];

  wire [11:0] control_remaining = remaining_bits[11:0];
  wire [15:0] control_read =
      control_address == 4'h0 ? DEVICE_ID :
      control_address == 4'h1 ? { flags, DEVICE_TYPE } :
      control_address == 4'h2 ? block_address[15:0] :
      control_address == 4'h5 ? { flags[0], 1'h0, state, control_remaining } :
      16'h0;

  always @(posedge cpu_clock) begin
    cpu_data_out <= is_control ? control_read : read_mem;
  end

  reg initiate = 1'h0;
  reg [7:0] read_word_offset = 8'h0;
  wire [15:0] block_address_read = ADDRESS_OFFSET + block_address;
  wire [31:0] byte_address_read = { block_address_read[14:0], 17'h0 };

  // Mapped memory, data is copied to/from this block which is mapped at
  // the upper memory block for the device
  wire [15:0] read_mem;
  reg [15:0] read_word = 16'h0;

  reg changes = 1'h0;
  basic_memory_block #(
    .WIDTH(8),
  ) mapped_mem (
    .clock(cpu_clock),
    .write_enable(next_word),
    .read_address(short_address),
    .write_address(read_word_offset),
    .data_in(read_word),
    .data_out(read_mem),
  );

  wire needs_enable = ~enabled && (flags_reg[1] | flags_reg[0]);
  wire needs_disable = enabled && ~(flags_reg[1] | flags_reg[0]);
  wire next_word = state == 2'h3
      && (remaining_bits[12] || (remaining_bits[11:4] != 8'hFF && remaining_bits[3:0] == 4'hF));

  wire bits_empty = remaining_bits == 13'h0;
  always @(posedge cpu_clock) begin
    if(control_write && state == 2'h0) begin
      if (control_address == 4'h1) begin
        flags_reg[1:0] <= cpu_data_in[11:10];
        flags_reg[3:2] <= cpu_data_in[13:12] > SUPPORTED_READ_MODES ?
            SUPPORTED_READ_MODES : cpu_data_in[13:12];
      end else if (control_address == 4'h2) begin
        initiate <= 1'h1;
        if (cpu_data_in > MAX_BLOCK_ADDRESS)
          block_address <= MAX_BLOCK_ADDRESS;
        else
          block_address <= cpu_data_in;
      end
    end else if (initiate == 1'h0 && state == 2'h0 && (needs_enable || needs_disable)) begin
      enabled <= flags_reg[1] | flags_reg[0];
      remaining_bits <= 13'h7;
      initiate <= 1'h0;
      state <= 2'h2;
      if (needs_enable) begin
        byte_buffer[31:24] <= RESUME_COMMAND;
      end else begin
        byte_buffer[31:24] <= POWER_DOWN;
      end
    end else begin
      initiate <= 1'h0;
      if (initiate) begin
        remaining_bits <= 13'h7;
        state <= 2'h1;
//        if (read_mode == 2'h0)
//          byte_buffer[31:24] <= READ_COMMAND;
//        else if (read_mode == 2'h1)
          byte_buffer[31:24] <= FAST_READ_COMMAND;
/*        else if (read_mode == 2'h2)
          byte_buffer <= FAST_DUAL_READ_COMMAND << 24;
        else
          byte_buffer <= FAST_QUAD_READ_COMMAND << 24;*/
      end else if (state == 2'h1) begin
        if (bits_empty) begin
          // 3 byte address, last 8-bits are ignored during read setup time
          byte_buffer <= byte_address_read;
          remaining_bits <= 13'h1F;
          state <= 2'h2;
        end else begin
          byte_buffer <= byte_buffer << 1;
          remaining_bits <= remaining_bits - 1'h1;
        end
      end else if (state == 2'h2) begin
        if (bits_empty) begin
          remaining_bits <= 13'h0FFF;
          state <= 2'h3;
          read_word_offset <= 8'h0;
          read_word <= 16'h0;
        end else begin
          byte_buffer <= byte_buffer << 1;
          remaining_bits <= remaining_bits - 1'h1;
        end
      end else if (state == 2'h3) begin
        if (next_word) begin
          read_word_offset <= read_word_offset + 1'h1;
          if (read_word_offset == 8'hFF)
            state <= 2'h0;
        end
//        if (read_mode[1] == 1'h0) begin
          read_word <= { read_word[14:0], MISO1 };
          remaining_bits <= remaining_bits - 13'h1;
          changes <= changes ^ read_word[14];
/*          read_word[remaining_bits[3:0]] <= MISO1;
          remaining_bits <= remaining_bits - 13'h1;
        end else if (read_mode == 2'h2) begin
          read_word[{ remaining_bits[3:1], 1'h0 }] <= MISO0;
          read_word[{ remaining_bits[3:1], 1'h1 }] <= MISO1;
          remaining_bits <= remaining_bits - 13'h2;
        end else if (read_mode == 2'h3) begin*/
          //read_word <= 16'h4541;
 //         read_word <= { read_word[11:0], MISO3, MISO2, MISO1, MISO0 };
          /*read_word[{ remaining_bits[3:2], 2'h0 }] <= MISO0;
          read_word[{ remaining_bits[3:2], 2'h1 }] <= MISO1;
          read_word[{ remaining_bits[3:2], 2'h2 }] <= MISO2;
          read_word[{ remaining_bits[3:2], 2'h3 }] <= MISO3;*/
 //         remaining_bits <= remaining_bits - 13'h4;
        //end
      end else begin
        remaining_bits <= 13'h0;
      end
    end
  end
endmodule

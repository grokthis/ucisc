module i2c_device(
    input cpu_clock,
    input write_enable,
    input is_control,
    input [7:0] short_address,
    input [15:0] cpu_data_in,
    output reg [15:0] cpu_data_out,
    input SDA_in,
    output reg SDA_out = 1'h1,
//    output SDA_out,
    output SDA_enable,
    output SCL_out
);
    parameter DEVICE_ID = 8'h10;
    parameter DEVICE_TYPE = 8'h9;
    parameter CPU_FREQ = 10000000;

    wire [2:0] speed_mode = flags[3:1];
    // Division value gives you half the clock cycle, so equals
    // twice the desired frequency.
    wire [15:0] division_value =
      speed_mode == 3'h1 ? CPU_FREQ / 200000 :
      speed_mode == 3'h2 ? CPU_FREQ / 800000 :
      speed_mode == 3'h3 ? CPU_FREQ / 2000000 :
      16'h0;
    wire SDA_enable = SDA_write & speed_mode != 3'h0;
    reg SDA_write = 1'h1;

    /********************************
     * Clock division for I2C clock *
     ********************************/
    wire [15:0] clk_divider = division_value;
    reg [15:0] clk_count = 16'h0;
    reg i2c_clk = 1'h0;
    reg i2c_change = 1'h1;

    always@(posedge cpu_clock & speed_mode != 3'h0) begin
      if (clk_count == clk_divider) begin
        i2c_clk <= ~i2c_clk;
        clk_count <= 16'h0;
      end else begin
        clk_count <= clk_count + 1'h1;
      end

      // i2c signal will change in the middle of the stable clock value
      // so it is always stable on i2c_clk edge
      if (clk_count == { 2'h0, clk_divider[15:2] }) begin
        i2c_change <= ~i2c_change;
      end
    end

    /************************
     * I2C read/write logic *
     ************************/
    reg i2c_clk_disable = 1'h1;

    assign SCL_out = i2c_clk_disable | i2c_clk;

    reg [127:0] control_buffer = 128'h0;
    reg [7:0] i2c_address = 8'h0;
    reg [15:0] bytes_initiated = 5'h0;
    reg [7:0] flags = 8'h0;

    reg started = 1'h0;
    reg read;
    reg write_done;
    reg read_done;
    reg [3:0] byte_pos;
    reg [7:0] buffer;
    reg [3:0] remaining_bits;

    wire [7:0] selected_byte =
        byte_pos == 4'hF ? control_buffer[7:0] :
        byte_pos == 4'hE ? control_buffer[15:8] :
        byte_pos == 4'hD ? control_buffer[23:16] :
        byte_pos == 4'hC ? control_buffer[31:24] :
        byte_pos == 4'hB ? control_buffer[39:32] :
        byte_pos == 4'hA ? control_buffer[47:40] :
        byte_pos == 4'h9 ? control_buffer[55:48] :
        byte_pos == 4'h8 ? control_buffer[63:56] :
        byte_pos == 4'h7 ? control_buffer[71:64] :
        byte_pos == 4'h6 ? control_buffer[79:72] :
        byte_pos == 4'h5 ? control_buffer[87:80] :
        byte_pos == 4'h4 ? control_buffer[95:88] :
        byte_pos == 4'h3 ? control_buffer[103:96] :
        byte_pos == 4'h2 ? control_buffer[111:104] :
        byte_pos == 4'h1 ? control_buffer[119:112] :
        control_buffer[127:120];

    reg SDA_capture;
    always @(posedge i2c_change) begin
      SDA_capture <= SDA_in;

      i2c_clk_disable <= ~started & stopped;
    end

    reg [15:0] bytes_remaining = 16'h0;
    reg done = 1'h1;
    reg stopped = 1'h1;
    wire has_remaining = bytes_remaining != 16'h0;
    always @(negedge i2c_change) begin
      if (started) begin
        // Send 8 bits while started
        started <= remaining_bits != 4'h0;
        remaining_bits <= remaining_bits - 1'h1;
        SDA_out <= buffer[7];
        SDA_write <= 1'h1;
        buffer <= buffer << 1;
        done <= 1'h0;
        stopped <= 1'h0;
      end else if(~done) begin
        // This is the ACK bit.
        if (read) begin
         SDA_out <= 1'h1;
         SDA_write <= 1'h1;
        end else begin
         SDA_write <= 1'h0;
        end
        // If remaining bytes, start next byte next
        // else stop condition next
        done <= ~has_remaining;
        started <= has_remaining;
        stopped <= 1'h0;
        if (has_remaining) begin
           buffer <= selected_byte;
           remaining_bits <= 4'h7;
           bytes_remaining <= bytes_remaining - 1'h1;
           byte_pos <= byte_pos + 1'h1;
        end
      end else if(stopped & SDA_out & bytes_initiated != 4'h0) begin
        started <= 1'h1;
        done <= 1'h0;
        byte_pos <= 4'h0;
        bytes_remaining <= bytes_initiated;

        // Start condition
        SDA_out <= 1'h0;
        SDA_write <= 1'h1;

        read <= i2c_address[0];
        buffer <= i2c_address[7:0];
        remaining_bits <= 4'h7;
      end else begin
        SDA_out <= stopped;
        SDA_write <= 1'h1;
        stopped <= 1'h1;
      end
    end

    wire control_write = is_control & write_enable;
    wire [3:0] control_address = short_address[3:0];

    wire [15:0] control_read =
        control_address == 4'h0 ? DEVICE_ID :
        control_address == 4'h1 ? { flags, DEVICE_TYPE } :
        control_address == 4'h2 ? { 8'h0, i2c_address } :
        control_address == 4'h3 ? bytes_initiated :
        control_address == 4'h4 ? control_buffer[127:112] :
        control_address == 4'h5 ? control_buffer[111:96] :
        control_address == 4'h6 ? control_buffer[95:80] :
        control_address == 4'h7 ? control_buffer[79:64] :
        control_address == 4'h8 ? control_buffer[63:48] :
        control_address == 4'h9 ? control_buffer[47:32] :
        control_address == 4'hA ? control_buffer[31:16] :
        control_address == 4'hB ? control_buffer[15:0] :
        16'h0 ;

    reg initiated = 1'h0;
    always @(posedge cpu_clock) begin
      cpu_data_out <= is_control ? control_read : 16'h0;

      if(control_write) begin
        if(control_address == 4'h1)
          flags <= cpu_data_in[3:1] > 3'h3 ? (cpu_data_in[15:8] & 8'hF1) : cpu_data_in[15:8];
        else if(control_address == 4'h2)
          i2c_address <= cpu_data_in[7:0];
        else if(control_address == 4'h3)
          bytes_initiated <= cpu_data_in[4:0];
        else if(control_address == 4'h4)
          control_buffer[127:112] <= cpu_data_in;
        else if(control_address == 4'h5)
          control_buffer[111:96] <= cpu_data_in;
        else if(control_address == 4'h6)
          control_buffer[95:80] <= cpu_data_in;
        else if(control_address == 4'h7)
          control_buffer[79:64] <= cpu_data_in;
        else if(control_address == 4'h8)
          control_buffer[63:48] <= cpu_data_in;
        else if(control_address == 4'h9)
          control_buffer[47:32] <= cpu_data_in;
        else if(control_address == 4'hA)
          control_buffer[31:16] <= cpu_data_in;
        else if(control_address == 4'hB)
          control_buffer[15:0] <= cpu_data_in;
      end else if (has_remaining | ~done | ~stopped) begin
        initiated <= 1'h1;
        bytes_initiated <= bytes_remaining + 1'h1;
      end else if (initiated & SDA_out) begin
        bytes_initiated <= 16'h0;
        initiated <= 1'h0;
       end
    end
endmodule

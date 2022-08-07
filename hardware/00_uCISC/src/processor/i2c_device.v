module i2c_device(
    input cpu_clock,
    input write_enable,
    input is_control,
    input [7:0] short_address,
    input [15:0] cpu_data_in,
    output [15:0] cpu_data_out,
    input SDA_in,
    output reg SDA_out = 1'h1,
//    output SDA_out,
    output reg SDA_enable = 1'h1,
    output SCL_out
);
    parameter DEVICE_ID = 16'h0;
    parameter DEVICE_TYPE = 8'h9;
    parameter CPU_FREQ = 10000000;

//    assign SDA_out = bytes_remaining[1] | bytes_remaining[0];

    wire [15:0] clk_400k_divider = CPU_FREQ / 800000;
    reg [15:0] clk_count = 16'h0;
    reg i2c_clk = 1'h0;
    reg i2c_change = 1'h1;

    always@(posedge cpu_clock) begin
      if (clk_count == clk_400k_divider) begin
        i2c_clk <= ~i2c_clk;
        clk_count <= 16'h0;
      end else begin
        clk_count <= clk_count + 1'h1;
      end

      // i2c signal will change in the middle of the stable clock value
      // so it is always stable on i2c_clk edge
      if (clk_count == { 2'h0, clk_400k_divider[15:2] }) begin
        i2c_change <= ~i2c_change;
      end
    end

    reg i2c_clk_disable = 1'h1;
    //reg i2c_clk_disable = 1'h0;

    assign SCL_out = i2c_clk_disable | i2c_clk;
    //assign SCL_out = i2c_clk;

    reg [191:0] control_buffer;
    reg [15:0] i2c_address;
    reg [15:0] bytes_initiated;

    reg started = 1'h0;
    reg read;
    reg write_done;
    reg read_done;
    reg [9:0] byte_pos;
    reg [7:0] buffer;
    reg [3:0] remaining_bits;
    wire ack = remaining_bits[3:0] == 4'h0;

    wire [7:0] selected_byte =
        byte_pos == 10'h0 ? control_buffer[7:0] :
        byte_pos == 10'h1 ? control_buffer[15:8] :
        byte_pos == 10'h2 ? control_buffer[23:16] :
        byte_pos == 10'h3 ? control_buffer[31:24] :
        byte_pos == 10'h4 ? control_buffer[39:32] :
        byte_pos == 10'h5 ? control_buffer[47:40] :
        byte_pos == 10'h6 ? control_buffer[55:48] :
        byte_pos == 10'h7 ? control_buffer[63:56] :
        byte_pos == 10'h8 ? control_buffer[71:64] :
        byte_pos == 10'h9 ? control_buffer[79:72] :
        byte_pos == 10'hA ? control_buffer[87:80] :
        byte_pos == 10'hB ? control_buffer[95:88] :
        byte_pos == 10'hC ? control_buffer[103:96] :
        byte_pos == 10'hD ? control_buffer[111:104] :
        byte_pos == 10'hE ? control_buffer[119:112] :
        byte_pos == 10'hF ? control_buffer[127:120] :
        byte_pos == 10'h10 ? control_buffer[135:128] :
        byte_pos == 10'h12 ? control_buffer[143:136] :
        byte_pos == 10'h13 ? control_buffer[151:144] :
        byte_pos == 10'h14 ? control_buffer[159:152] :
        byte_pos == 10'h15 ? control_buffer[167:160] :
        byte_pos == 10'h16 ? control_buffer[175:168] :
        byte_pos == 10'h17 ? control_buffer[183:168] :
        control_buffer[191:168];

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
        SDA_enable <= 1'h1;
        buffer <= buffer << 1;
        done <= 1'h0;
        stopped <= 1'h0;
      end else if(~done) begin
        // This is the ACK bit.
        if (read) begin
         SDA_out <= 1'h1;
         SDA_enable <= 1'h1;
        end else begin
         SDA_enable <= 1'h0;
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
        byte_pos <= 10'h0;
        bytes_remaining <= bytes_initiated;

        // Start condition - TODO: fix, may be repeated start
        SDA_out <= 1'h0;
        SDA_enable <= 1'h1;

        //read <= i2c_address[0];
        read = 1'h0; //testing writes first

        buffer <= i2c_address[7:0];
        remaining_bits <= 4'h7;
      end else begin
        SDA_out <= stopped;
        SDA_enable <= 1'h1;
        stopped <= 1'h1;

//        if(bytes_queued != 4'h0) begin
//          started <= 1'h1;
//          done <= 1'h0;
//          byte_pos <= 10'h0;
//          bytes_remaining <= bytes_initiated;
//
//          // Start condition - TODO: fix, may be repeated start
//          SDA_out <= 1'h0;
//          SDA_enable <= 1'h1;
//
//          //read <= i2c_address[0];
//          read = 1'h0; //testing writes first
//
//          buffer <= i2c_address[7:0];
//          remaining_bits <= 4'h7;
//        end
      end
    end

    wire control_write = is_control & write_enable;
    wire [3:0] control_address = short_address[3:0];

    wire [7:0] flags = { pin_count, 2'h0 };
    wire [15:0] control_read =
        control_address == 4'h0 ? DEVICE_ID :
        control_address == 4'h1 ? { flags, DEVICE_TYPE } :
        control_address == 4'h2 ? i2c_address :
        control_address == 4'h3 ? bytes_initiated :
        control_address == 4'h4 ? control_buffer[15:0] :
        control_address == 4'h5 ? control_buffer[31:16]:
        control_address == 4'h6 ? control_buffer[47:32]:
        control_address == 4'h7 ? control_buffer[63:48]:
        control_address == 4'h8 ? control_buffer[79:64] :
        control_address == 4'h9 ? control_buffer[95:80] :
        control_address == 4'hA ? control_buffer[111:96] :
        control_address == 4'hB ? control_buffer[127:112] :
        control_address == 4'hC ? control_buffer[143:128] :
        control_address == 4'hD ? control_buffer[159:144] :
        control_address == 4'hE ? control_buffer[175:160] :
        control_address == 4'hF ? control_buffer[191:176] : 16'h0;

    assign cpu_data_out = is_control ? control_read : 16'h0;
    always @(posedge cpu_clock) begin
      if(control_write) begin
        if(control_address == 4'h2)
          i2c_address <= cpu_data_in;
        if(control_address == 4'h4)
          control_buffer[15:0] <= cpu_data_in;
        else if(control_address == 4'h5)
          control_buffer[31:16] <= cpu_data_in;
        else if(control_address == 4'h6)
          control_buffer[47:32] <= cpu_data_in;
        else if(control_address == 4'h7)
          control_buffer[63:48] <= cpu_data_in;
        else if(control_address == 4'h8)
          control_buffer[79:64] <= cpu_data_in;
        else if(control_address == 4'h9)
          control_buffer[95:80] <= cpu_data_in;
        else if(control_address == 4'hA)
          control_buffer[111:96] <= cpu_data_in;
        else if(control_address == 4'hB)
          control_buffer[127:112] <= cpu_data_in;
        else if(control_address == 4'hC)
          control_buffer[143:128] <= cpu_data_in;
        else if(control_address == 4'hD)
          control_buffer[159:144] <= cpu_data_in;
        else if(control_address == 4'hE)
          control_buffer[175:169] <= cpu_data_in;
        else if(control_address == 4'hF)
          control_buffer[191:176] <= cpu_data_in;
      end
      if (~stopped) begin
        bytes_initiated <= bytes_remaining;
      end else if(control_write & control_address == 4'h3) begin
//      if(control_write & control_address == 4'h3) begin
        bytes_initiated <= cpu_data_in;
      end
    end
endmodule

module gpio_device(
    input cpu_clock,
    input write_enable,
    input is_control,
    input [7:0] short_address,
    input [15:0] cpu_data_in,
    output reg [15:0] cpu_data_out,
    input [PINS-1:0] gpio_in,
    output reg [PINS-1:0] gpio_out,
    output reg [PINS-1:0] gpio_config,
);
    parameter PINS = 16;
    parameter DEVICE_ID = 16'h0;
    parameter DEVICE_TYPE = 8'h8;

    wire [5:0] pin_count = PINS - 1;

    wire control_write = is_control & write_enable;
    wire [3:0] control_address = short_address[3:0];

    wire [7:0] flags = { pin_count, 2'h0 };
    wire [15:0] control_read =
        control_address == 4'h0 ? DEVICE_ID :
        control_address == 4'h1 ? { flags, DEVICE_TYPE } :
        control_address == 4'h2 ? 16'h0 :
        control_address == 4'h3 ? 16'h0 :
        control_address == 4'h4 ? gpio_config[15:0] :
        control_address == 4'h5 ? (PINS > 16 ? gpio_config[31:16] : 16'h0) :
        control_address == 4'h6 ? (PINS > 32 ? gpio_config[47:32] : 16'h0) :
        control_address == 4'h7 ? (PINS > 48 ? gpio_config[63:48] : 16'h0) :
        control_address == 4'h8 ? gpio_out[15:0] :
        control_address == 4'h9 ? (PINS > 16 ? gpio_out[31:0] : 16'h0) :
        control_address == 4'hA ? (PINS > 16 ? gpio_out[47:32] : 16'h0) :
        control_address == 4'hB ? (PINS > 16 ? gpio_out[63:48] : 16'h0) :
        control_address == 4'hC ? gpio_in[15:0] :
        control_address == 4'hD ? (PINS > 16 ? gpio_in[31:16] : 16'h0) :
        control_address == 4'hE ? (PINS > 16 ? gpio_in[47:32] : 16'h0) :
        (PINS > 16 ? gpio_in[63:48] : 16'h0);

    always @(posedge cpu_clock) begin
      cpu_data_out <= is_control ? control_read : 16'h0;
    end

    always @(posedge cpu_clock) begin
      if(control_write) begin
        if(control_address == 4'h4)
          gpio_config[15:0] <= cpu_data_in;
        else if(PINS > 16 && control_address == 4'h5)
          gpio_config[31:16] <= cpu_data_in;
        else if(PINS > 32 && control_address == 4'h6)
          gpio_config[47:32] <= cpu_data_in;
        else if(PINS > 48 && control_address == 4'h7)
          gpio_config[63:48] <= cpu_data_in;
        else if(control_address == 4'h8)
          gpio_out[15:0] <= cpu_data_in;
        else if(PINS > 16 && control_address == 4'h9)
          gpio_out[31:16] <= cpu_data_in;
        else if(PINS > 32 && control_address == 4'hA)
          gpio_out[47:32] <= cpu_data_in;
        else if(PINS > 48 && control_address == 4'hB)
          gpio_out[63:48] <= cpu_data_in;
      end
    end
endmodule

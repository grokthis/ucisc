module gpio_device(
    input cpu_clock,
    input write_enable,
    input is_control,
    input [7:0] short_address,
    input [15:0] cpu_data_in,
    output [15:0] cpu_data_out,
    input [63:0] gpio_in,
    output [63:0] gpio_out,
    output [63:0] gpio_config,
);
    parameter PINS = 16;
    parameter DEVICE_ID = 16'h0;
    parameter DEVICE_TYPE = 8'h8;

    wire [5:0] pin_count = PINS - 1;

    reg [63:0] config_reg = 64'h0;
    assign gpio_config = config_reg;
    reg [63:0] outputs = 64'h0;
    assign gpio_out = (config_reg & outputs) | (~config_reg & gpio_in);

    wire control_write = is_control & write_enable;
    wire [3:0] control_address = short_address[3:0];

    wire [7:0] flags = { pin_count, 2'h0 };
    wire [15:0] control_read =
        control_address == 4'h0 ? DEVICE_ID :
        control_address == 4'h1 ? { flags, DEVICE_TYPE } :
        control_address == 4'h2 ? 16'h0 :
        control_address == 4'h3 ? 16'h0 :
        control_address == 4'h4 ? gpio_config[15:0] :
        control_address == 4'h5 ? gpio_config[31:16]:
        control_address == 4'h6 ? gpio_config[47:32]:
        control_address == 4'h7 ? gpio_config[63:48]:
        control_address == 4'h8 ? outputs[15:0] :
        control_address == 4'h9 ? outputs[31:0] :
        control_address == 4'hA ? outputs[47:32] :
        control_address == 4'hB ? outputs[63:48] :
        control_address == 4'hC ? gpio_out[15:0] :
        control_address == 4'hD ? gpio_out[31:16] :
        control_address == 4'hE ? gpio_out[47:32] :
        control_address == 4'hF ? gpio_out[63:48] : 16'h0;

    assign cpu_data_out = is_control ? control_read : 16'h0;
    always @(posedge cpu_clock) begin
      if(control_write) begin
        if(control_address == 4'h4)
          config_reg[15:0] <= cpu_data_in;
        else if(control_address == 4'h5)
          config_reg[31:16] <= cpu_data_in;
        else if(control_address == 4'h6)
          config_reg[47:32] <= cpu_data_in;
        else if(control_address == 4'h7)
          config_reg[63:48] <= cpu_data_in;
        else if(control_address == 4'h8)
          outputs[15:0] <= cpu_data_in;
        else if(control_address == 4'h9)
          outputs[31:16] <= cpu_data_in;
        else if(control_address == 4'hA)
          outputs[47:32] <= cpu_data_in;
        else if(control_address == 4'hB)
          outputs[63:48] <= cpu_data_in;
      end
    end
endmodule

module gk110 (
    input cpu_clock,
    input reset,
    output SCK,
    output MO,
    input MI,
    input RX,
    output TX,
    output SDA_out,
    input SDA_in,
    output SDA_config,
    output SCL,
    output LED_out,
    output A0_out,
    output A1_out,
    output A2_out,
    output A3_out,
    output D24_out,
    output D25_out,
    output D26_out,
    output D4_out,
    output D5_out,
    output D6_out,
    output D9_out,
    output D10_out,
    output D11_out,
    output D12_out,
    output D13_out,
    input LED_in,
    input A0_in,
    input A1_in,
    input A2_in,
    input A3_in,
    input D24_in,
    input D25_in,
    input D26_in,
    input D4_in,
    input D5_in,
    input D6_in,
    input D9_in,
    input D10_in,
    input D11_in,
    input D12_in,
    input D13_in,
    output LED_config,
    output A0_config,
    output A1_config,
    output A2_config,
    output A3_config,
    output D24_config,
    output D25_config,
    output D26_config,
    output D4_config,
    output D5_config,
    output D6_config,
    output D9_config,
    output D10_config,
    output D11_config,
    output D12_config,
    output D13_config,
    output SPI_CS,
    output SPI_SCK,
    output SPI_IO0_out,
    output SPI_IO0_config,
    input SPI_IO0_in,
    input SPI_IO1,
    input SPI_IO2,
    input SPI_IO3
);
  parameter CPU_FREQ = 10000000;
  parameter MEM_ADDRESS_WIDTH = 16;

  wire device_write_en;
  wire [15:0] device_address;
  wire [15:0] device_data_in;
  wire [15:0] device_data_out;

  wire [63:0] gpio_read = {
      48'h0,
      A3_in, A2_in, A1_in, A0_in,
      D26_in, D25_in, D24_in, D13_in, D12_in, D11_in, D10_in, D9_in, D6_in, D5_in, D4_in,
      LED_in
  };

  assign LED_out = gpio_write[0];
  assign D4_out = gpio_write[1];
  assign D5_out = gpio_write[2];
  assign D6_out = gpio_write[3];
  assign D9_out = gpio_write[4];
  assign D10_out = gpio_write[5];
  assign D11_out = gpio_write[6];
  assign D12_out = gpio_write[7];
  assign D13_out = gpio_write[8];
  assign D24_out = gpio_write[9];
  assign D25_out = gpio_write[10];
  assign D26_out = gpio_write[11];
  assign A0_out = gpio_write[12];
  assign A1_out = gpio_write[13];
  assign A2_out = gpio_write[14];
  assign A3_out = gpio_write[15];

  wire [63:0] gpio_write;
  wire [63:0] gpio_config;
  assign LED_config = gpio_config[0];
  assign D4_config = gpio_config[1];
  assign D5_config = gpio_config[2];
  assign D6_config = gpio_config[3];
  assign D9_config = gpio_config[4];
  assign D10_config = gpio_config[5];
  assign D11_config = gpio_config[6];
  assign D12_config = gpio_config[7];
  assign D13_config = gpio_config[8];
  assign D24_config = gpio_config[9];
  assign D25_config = gpio_config[10];
  assign D26_config = gpio_config[11];
  assign A0_config = gpio_config[12];
  assign A1_config = gpio_config[13];
  assign A2_config = gpio_config[14];
  assign A3_config = gpio_config[15];

  wire [7:0] control_device_id = device_address[11:4];
  wire [7:0] mem_device_id = device_address[15:8];
  wire is_control = device_address[15:12] == 4'h0;
  wire [7:0] device_id = is_control ? control_device_id : mem_device_id;

  wire [15:0] gpio_device_data_in;
  wire [15:0] i2c_device_data_in;
  wire [15:0] flash_spi_device_data_in;

  reg [7:0] captured_device_id;

  always @(posedge cpu_clock) begin
    captured_device_id <= device_id;
  end
  assign device_data_in =
      captured_device_id == 8'h4 ? gpio_device_data_in :
      captured_device_id == 8'h10 ? i2c_device_data_in :
      captured_device_id == 8'h40 ? flash_spi_device_data_in :
      16'h0;

  gpio_device #(
    .PINS(16),
    .DEVICE_ID(8'h01)
  ) gpio (
    .cpu_clock(cpu_clock),
    .write_enable(device_write_en & device_id == 8'h4),
    .is_control(is_control),
    .short_address(device_address[7:0]),
    .cpu_data_in(device_data_out),
    .cpu_data_out(gpio_device_data_in),
    .gpio_in(gpio_read),
    .gpio_out(gpio_write),
    .gpio_config(gpio_config)
  );

  //assign SDA_config = 1'b1;
  //assign SDA_out = 1'b1;
  //assign SCL = 1'b1;
  i2c_device #(
    .DEVICE_ID(8'h10),
    .CPU_FREQ(CPU_FREQ)
  ) i2c_device(
    .cpu_clock(cpu_clock),
    .write_enable(device_write_en && device_id == 8'h10),
    .is_control(is_control),
    .short_address(device_address[7:0]),
    .cpu_data_in(device_data_out),
    .cpu_data_out(i2c_device_data_in),
    .SDA_in(SDA_in),
    .SDA_out(SDA_out),
    .SDA_enable(SDA_config),
    .SCL_out(SCL)
  );

  // Uses the TinyFPGA onboard flash, starting at address 0x50000
  spi_flash_device #(
    .DEVICE_ID(8'h40),
    .CPU_FREQ(CPU_FREQ),
    .ADDRESS_OFFSET(32'h280),
    .MAX_BLOCK_ADDRESS(32'h57F) // 8Mbit flash
  ) spi_flash_device(
    .cpu_clock(cpu_clock),
    .write_enable(device_write_en && device_id == 8'h40),
    .is_control(is_control),
    .short_address(device_address[7:0]),
    .cpu_data_in(device_data_out),
    .cpu_data_out(flash_spi_device_data_in),
    .MOSI(SPI_IO0_out),
    .MOSI_enable(SPI_IO0_config),
    .MISO0(SPI_IO0_in),
    .MISO1(SPI_IO1),
    .MISO2(SPI_IO2),
    .MISO3(SPI_IO3),
    .SCK_out(SPI_SCK),
    .CSLow_out(SPI_CS),
  );

  cpu #(
    .MEM_ADDRESS_WIDTH(MEM_ADDRESS_WIDTH)
  ) cpu(
    .cpu_clock(cpu_clock),
    .reset(reset),
    .device_write_en(device_write_en),
    .device_address(device_address),
    .device_data_in(device_data_in),
    .device_data_out(device_data_out)
  );

//    wire control = address[15:12] == 4'h0;
//    wire [7:0] target_device = control ? address[11:4] : address[15:8];
//
//    always @(posedge cpu_clock) begin
//        data_out <=
//            target_device == 8'h02 ? uart_data_out :
//            target_device == 8'h10 ? 16'h0 : //graphics_data_out :
//            16'h0;
//    end
//
//    wire [15:0] uart_data_out;
//    uart_device #(.DEVICE_ID(8'h14)) uart_device(
//        .clock(cpu_clock),
//        .write_enable(write_enable & target_device == 8'h2),
//        .control(control),
//        .address(address[7:0]),
//        .data_in(data_in),
//        .rx(RX),
//        .data_out(uart_data_out),
//        .tx(TX),
//        .peek(peek)
//    );

endmodule
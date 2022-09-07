module top (
  input CLK, // ref clock
  input PIN_8, // SCK
  output PIN_9, // MO
  input PIN_10, // MI
  input PIN_11, // RX
  output PIN_12, // TX
  inout PIN_14, // sda
  output PIN_15, // scl
  //inout PIN_16, // D5
  inout LED, // D5 as LED
  inout PIN_17, // D6
  inout PIN_18, // D9
  inout PIN_19, // D10
  inout PIN_20, // D11
  inout PIN_21, // D12
  inout PIN_22, // D13

  //output PIN_2, // temp
  //output PIN_3, // temp
  //output PIN_4, // temp
  //output PIN_5, // temp
  //output PIN_6, // temp
  //output PIN_7, // temp

  // SPI Flash
  output SPI_SS,
  output SPI_SCK,
  inout SPI_IO0,
  input SPI_IO1,
  input SPI_IO2,
  input SPI_IO3
);
  // ref_clock = 16MHz
  parameter CLOCK_DIV = 7; // max 15
  parameter CLOCK_MULT = 17; // max 63
  parameter CPU_FREQ = (16 / (CLOCK_DIV + 1)) * (CLOCK_MULT + 1) * 1000000;

  //assign PIN_2 = SPI_SCK;
  //assign PIN_3 = SPI_SS;
  //assign PIN_4 = SPI_IO0_out;
  //assign PIN_5 = SPI_IO1;
  //assign PIN_6 = SPI_IO2;
  //assign PIN_7 = SPI_IO3;

  // There seems to be some sort of lag in the bootloader startup process
  // If we start the CPU too quickly, it doesn't often execute correctly
  // We leave the reset value high for the first ~1 million clock cycles
  reg [19:0] pre_count = 20'h0;
  reg reset = 1'b1;
  always @(posedge cpu_clock) begin
    if (pre_count == 20'hFFFFF)
      reset <= 1'b0;
    else
      pre_count <= pre_count + 1'h1;
  end

  // Yosys doesn't correctly handle tristate outputs yet, so
  // we have to directly call the cell library directly.

  // Tristate config for SDA
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) sda (
    .PACKAGE_PIN(PIN_14),
    .OUTPUT_ENABLE(SDA_config),
    .D_OUT_0(SDA_out),
    .D_IN_0(SDA_in)
  );

  // Tristate config for D5
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) D5_IO (
    //.PACKAGE_PIN(PIN_16),
    .PACKAGE_PIN(LED),
    .OUTPUT_ENABLE(D5_config),
    .D_OUT_0(D5_out),
    .D_IN_0(D5_in)
  );

  // Tristate config for D6
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) D6_IO (
    .PACKAGE_PIN(PIN_17),
    .OUTPUT_ENABLE(1'h1),//D6_config),
    .D_OUT_0(D6_out),
    .D_IN_0(D6_in)
  );

  // Tristate config for D9
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) D9_IO (
    .PACKAGE_PIN(PIN_18),
    .OUTPUT_ENABLE(D9_config),
    .D_OUT_0(D9_out),
    .D_IN_0(D9_in)
  );

  // Tristate config for D10
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) D10_IO (
    .PACKAGE_PIN(PIN_19),
    .OUTPUT_ENABLE(D10_config),
    .D_OUT_0(D10_out),
    .D_IN_0(D10_in)
  );

  // Tristate config for D11
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) D11_IO (
    .PACKAGE_PIN(PIN_20),
    .OUTPUT_ENABLE(D11_config),
    .D_OUT_0(D11_out),
    .D_IN_0(D11_in)
  );

  // Tristate config for D12
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) D12_IO (
    .PACKAGE_PIN(PIN_21),
    .OUTPUT_ENABLE(D12_config),
    .D_OUT_0(D12_out),
    .D_IN_0(D12_in)
  );

  // Tristate config for D13
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) D13_IO (
    .PACKAGE_PIN(PIN_22),
    .OUTPUT_ENABLE(D13_config),
    .D_OUT_0(D13_out),
    .D_IN_0(D13_in)
  );

  // Tristate config for SPI_IO0
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) SPI_IO0_IO (
    .PACKAGE_PIN(SPI_IO0),
    .OUTPUT_ENABLE(1'h1),
    .D_OUT_0(SPI_IO0_out),
    .D_IN_0(SPI_IO0_in)
  );

  wire cpu_clock = clock_buffer;//cpu_clock_div[4];
  pll #(
    .CLK_DIVR(CLOCK_DIV),
    .CLK_DIVF(CLOCK_MULT)
  )
  pll_cpu (
    .clki(CLK),
    .clko(clock_buffer)
  );
  wire clock_buffer;
  reg [4:0] cpu_clock_div = 5'h0;
  always @(posedge clock_buffer) begin
    cpu_clock_div = cpu_clock_div + 1'h1;
  end


  wire D5_out;
  wire D6_out;
  wire D9_out;
  wire D10_out;
  wire D11_out;
  wire D12_out;
  wire D13_out;
  wire D5_in;// = LED;
  wire D6_in;// = PIN_17;
  wire D9_in;// = PIN_18;
  wire D10_in;// = PIN_19;
  wire D11_in;// = PIN_20;
  wire D12_in;// = PIN_21;
  wire D13_in;// = PIN_22;
  wire D5_config;
  wire D6_config;
  wire D9_config;
  wire D10_config;
  wire D11_config;
  wire D12_config;
  wire D13_config;
  wire SDA_out;
  wire SDA_in;
  wire SDA_config;

  wire SPI_IO0_out;
  wire SPI_IO0_in;
  wire SPI_IO0_config;

  gk110 #(
    .MEM_ADDRESS_WIDTH(12),
    .CPU_FREQ(CPU_FREQ)
  )
  gk110 (
      .cpu_clock(cpu_clock),
      .reset(reset),
      .SCK(PIN_8),
      .MO(PIN_9),
      .MI(PIN_10),
      .RX(PIN_11),
      .TX(PIN_12),
      .SDA_out(SDA_out),
      .SDA_in(SDA_in),
      .SDA_config(SDA_config),
      .SCL(PIN_15),
      .D5_out(D5_out),
      .D6_out(D6_out),
      .D9_out(D9_out),
      .D10_out(D10_out),
      .D11_out(D11_out),
      .D12_out(D12_out),
      .D13_out(D13_out),
      .D5_in(D5_in),
      .D6_in(D6_in),
      .D9_in(D9_in),
      .D10_in(D10_in),
      .D11_in(D11_in),
      .D12_in(D12_in),
      .D13_in(D13_in),
      .D5_config(D5_config),
      .D6_config(D6_config),
      .D9_config(D9_config),
      .D10_config(D10_config),
      .D11_config(D11_config),
      .D12_config(D12_config),
      .D13_config(D13_config),
      .SPI_CS(SPI_SS),
      .SPI_SCK(SPI_SCK),
      .SPI_IO0_out(SPI_IO0_out),
      .SPI_IO0_in(SPI_IO0_in),
      .SPI_IO0_config(SPI_IO0_config),
      .SPI_IO1(SPI_IO1),
      .SPI_IO2(SPI_IO2),
      .SPI_IO3(SPI_IO3)
  );
endmodule

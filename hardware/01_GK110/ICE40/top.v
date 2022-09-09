module top (
  input CLK, // ref clock
  input PIN_8, // SCK
  output PIN_9, // MO
  input PIN_10, // MI
  input PIN_11, // RX
  output PIN_12, // TX
  inout PIN_14, // sda
  output PIN_15, // scl
  inout LED,
  inout PIN_2, // A0
  inout PIN_3, // A1
  inout PIN_4, // A2
  inout PIN_5, // A3
  inout PIN_6, // D24
  inout PIN_7, // D25
  inout PIN_13, // D4
  inout PIN_16, // D5
  inout PIN_17, // D6
  inout PIN_18, // D9
  inout PIN_19, // D10
  inout PIN_20, // D11
  inout PIN_21, // D12
  inout PIN_22, // D13
  inout PIN_23,

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
  parameter CLOCK_MULT = 19; // max 63
  parameter CPU_FREQ = (16 / (CLOCK_DIV + 1)) * (CLOCK_MULT + 1) * 1000000;

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

  // Tristate config for LED
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) LED_IO (
    .PACKAGE_PIN(LED),
    .OUTPUT_ENABLE(LED_config),
    .D_OUT_0(LED_out),
    .D_IN_0(LED_in)
  );

  // Tristate config for D4
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) D4_IO (
    .PACKAGE_PIN(PIN_13),
    .OUTPUT_ENABLE(D4_config),
    .D_OUT_0(D4_out),
    .D_IN_0(D4_in)
  );

  // Tristate config for D5
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) D5_IO (
    .PACKAGE_PIN(PIN_16),
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

  // Tristate config for D26
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) D26_IO (
    .PACKAGE_PIN(PIN_23),
    .OUTPUT_ENABLE(D26_config),
    .D_OUT_0(D26_out),
    .D_IN_0(D26_in)
  );

  // Tristate config for D25
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) D25_IO (
    .PACKAGE_PIN(PIN_7),
    .OUTPUT_ENABLE(D25_config),
    .D_OUT_0(D25_out),
    .D_IN_0(D25_in)
  );

  // Tristate config for D24
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) D24_IO (
    .PACKAGE_PIN(PIN_6),
    .OUTPUT_ENABLE(D24_config),
    .D_OUT_0(D24_out),
    .D_IN_0(D24_in)
  );

  // Tristate config for A3
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) A3_IO (
    .PACKAGE_PIN(PIN_5),
    .OUTPUT_ENABLE(A3_config),
    .D_OUT_0(A3_out),
    .D_IN_0(A3_in)
  );

  // Tristate config for A2
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) A2_IO (
    .PACKAGE_PIN(PIN_4),
    .OUTPUT_ENABLE(A2_config),
    .D_OUT_0(A2_out),
    .D_IN_0(A2_in)
  );

  // Tristate config for A1
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) A1_IO (
    .PACKAGE_PIN(PIN_3),
    .OUTPUT_ENABLE(A1_config),
    .D_OUT_0(A1_out),
    .D_IN_0(A1_in)
  );

  // Tristate config for A0
  SB_IO #(
    .PIN_TYPE(6'b 1010_01) // simple input, tristate output
  ) A0_IO (
    .PACKAGE_PIN(PIN_2),
    .OUTPUT_ENABLE(A0_config),
    .D_OUT_0(A0_out),
    .D_IN_0(A0_in)
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

  wire LED_out;
  wire A0_out;
  wire A1_out;
  wire A2_out;
  wire A3_out;
  wire D24_out;
  wire D25_out;
  wire D26_out;
  wire D27_out;
  wire D4_out;
  wire D5_out;
  wire D6_out;
  wire D9_out;
  wire D10_out;
  wire D11_out;
  wire D12_out;
  wire D13_out;
  wire LED_in;
  wire A0_in;
  wire A1_in;
  wire A2_in;
  wire A3_in;
  wire D24_in;
  wire D25_in;
  wire D26_in;
  wire D27_in;
  wire D4_in;
  wire D5_in;// = LED;
  wire D6_in;// = PIN_17;
  wire D9_in;// = PIN_18;
  wire D10_in;// = PIN_19;
  wire D11_in;// = PIN_20;
  wire D12_in;// = PIN_21;
  wire D13_in;// = PIN_22;
  wire LED_config;
  wire A0_config;
  wire A1_config;
  wire A2_config;
  wire A3_config;
  wire D24_config;
  wire D25_config;
  wire D26_config;
  wire D27_config;
  wire D4_config;
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
      .LED_out(LED_out),
      .A0_out(A0_out),
      .A1_out(A1_out),
      .A2_out(A2_out),
      .A3_out(A3_out),
      .D24_out(D24_out),
      .D25_out(D25_out),
      .D26_out(D26_out),
      .D4_out(D4_out),
      .D5_out(D5_out),
      .D6_out(D6_out),
      .D9_out(D9_out),
      .D10_out(D10_out),
      .D11_out(D11_out),
      .D12_out(D12_out),
      .D13_out(D13_out),
      .LED_in(LED_in),
      .A0_in(A0_in),
      .A1_in(A1_in),
      .A2_in(A2_in),
      .A3_in(A3_in),
      .D24_in(D24_in),
      .D25_in(D25_in),
      .D26_in(D26_in),
      .D4_in(D4_in),
      .D5_in(D5_in),
      .D6_in(D6_in),
      .D9_in(D9_in),
      .D10_in(D10_in),
      .D11_in(D11_in),
      .D12_in(D12_in),
      .D13_in(D13_in),
      .LED_config(LED_config),
      .A0_config(A0_config),
      .A1_config(A1_config),
      .A2_config(A2_config),
      .A3_config(A3_config),
      .D24_config(D24_config),
      .D25_config(D25_config),
      .D26_config(D26_config),
      .D4_config(D4_config),
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

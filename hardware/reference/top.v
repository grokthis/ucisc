module top(
    input CLK,
    input PIN_2,
    output LED,
    output PIN_24,
    output PIN_23,
    output PIN_22,
    output PIN_21,
    output PIN_20,
    output PIN_19,
    output PIN_18,
    output PIN_17,
    output PIN_16,
    output PIN_15,
    output PIN_14,
    output PIN_13,
    output PIN_12,
    output PIN_11,
    output PIN_10,
    output PIN_9,
    output PIN_8,
    output PIN_7,
    output PIN_6,
    output PIN_5,
    output PIN_4,
    output PIN_3,
    output PIN_2,
    output PIN_1
);

parameter address_width = 13;

assign PIN_24 = r1[15];
assign PIN_23 = r1[14];
assign PIN_22 = r1[13];
assign PIN_21 = r1[12];
assign PIN_20 = r1[11];
assign PIN_19 = r1[10];
assign PIN_18 = r1[9];
assign PIN_17 = r1[8];

assign PIN_16 = r1[7];
assign PIN_15 = r1[6];
assign PIN_14 = r1[5];
assign PIN_9 = r1[4];
assign PIN_10 = r1[3];
assign PIN_11 = r1[2];
assign PIN_12 = r1[1];
assign PIN_13 = r1[0];

assign PIN_3 = pc[3];
assign PIN_4 = pc[2];
assign PIN_5 = pc[1];
assign PIN_6 = pc[0];
assign PIN_7 = step[1];
assign PIN_8 = step[0];
assign PIN_1 = tx;

reg [24:0] slow_clock;
initial slow_clock = 25'h0;

always @(posedge CLK) begin // & PIN_1
    slow_clock <= slow_clock + 1;
end

wire clock = slow_clock[0];

//reg clock;
//initial clock = 0;
//always @* begin
//    // Debounce clock buttons
//    clock <= PIN_1 | ((~PIN_2) & clock);
//end
//assign LED = store_reg;

wire [1:0] step;
wire [15:0] r1;
wire [15:0] pc;
wire tx;
//assign LED = tx;

cpu cpu(
    .clock_input(clock),
    .reset(PIN_2),
    .step(step),
    .r1_peek(r1),
    .pc_peek(pc),
    .tx(tx)
);

endmodule
module register_block (
    input clock,
    input [1:0] step,
    input reset,
    input [3:0] desired_source,
    input [3:0] desired_destination,
    input [15:0] write_value,
    input write_enable,
    input push,
    input pop,
    input inc_enable,
    input [15:0] flags_in,
    input write_flags,
    output [15:0] source_value,
    output [15:0] destination_value,
    output [15:0] pc,
    output [15:0] flags,
    output [15:0] banking,
    output source_banked,
    output destination_banked,
    output [15:0] r1_peek
);

// r1_peek is for debug purposes
assign r1_peek = r1_out;

assign pc = pc_out;
assign flags = flags_out;
assign banking = banking_out;
assign source_banked =
    captured_desired_source == 4'h1 ? banking[1] :
    captured_desired_source == 4'h2 ? banking[2] :
    captured_desired_source == 4'h3 ? banking[3] :
    captured_desired_source == 4'h9 ? banking[5] :
    captured_desired_source == 4'hA ? banking[6] :
    captured_desired_source == 4'hB ? banking[7] :
    0;

assign destination_banked =
    captured_desired_destination == 4'h1 ? banking[1] :
    captured_desired_destination == 4'h2 ? banking[2] :
    captured_desired_destination == 4'h3 ? banking[3] :
    captured_desired_destination == 4'h9 ? banking[5] :
    captured_desired_destination == 4'hA ? banking[6] :
    captured_desired_destination == 4'hB ? banking[7] :
    0;

wire [3:0] captured_desired_source;
value_capture #(.WIDTH(4)) desired_source_capture (
    .clock(clock),
    .current_step(step),
    .capture_on(2'h1),
    .input_value(desired_source),
    .captured_out(captured_desired_source)
);

assign source_value =
    captured_desired_source == 4'h0 ? pc_out :
    captured_desired_source == 4'h1 ? r1_out :
    captured_desired_source == 4'h2 ? r2_out :
    captured_desired_source == 4'h3 ? r3_out :
    captured_desired_source == 4'h4 ? 16'h0 :
    captured_desired_source == 4'h5 ? r1_out :
    captured_desired_source == 4'h6 ? r2_out :
    captured_desired_source == 4'h7 ? r3_out :
    captured_desired_source == 4'h8 ? flags_out :
    captured_desired_source == 4'h9 ? rb1_out :
    captured_desired_source == 4'hA ? rb2_out :
    captured_desired_source == 4'hB ? rb3_out :
    captured_desired_source == 4'hC ? interrupt_out :
    captured_desired_source == 4'hD ? rb1_out :
    captured_desired_source == 4'hE ? rb2_out : rb3_out;

wire [3:0] captured_desired_destination;
value_capture #(.WIDTH(4)) desired_destination_capture (
    .clock(clock),
    .current_step(step),
    .capture_on(2'h1),
    .input_value(desired_destination),
    .captured_out(captured_desired_destination)
);

assign destination_value =
    captured_desired_destination == 4'h0 ? pc_out :
    captured_desired_destination == 4'h1 ? r1_out :
    captured_desired_destination == 4'h2 ? r2_out :
    captured_desired_destination == 4'h3 ? r3_out :
    captured_desired_destination == 4'h4 ? banking_out :
    captured_desired_destination == 4'h5 ? r1_out :
    captured_desired_destination == 4'h6 ? r2_out :
    captured_desired_destination == 4'h7 ? r3_out :
    captured_desired_destination == 4'h8 ? flags_out :
    captured_desired_destination == 4'h9 ? rb1_out :
    captured_desired_destination == 4'hA ? rb2_out :
    captured_desired_destination == 4'hB ? rb3_out :
    captured_desired_destination == 4'hC ? interrupt_out :
    captured_desired_destination == 4'hD ? rb1_out :
    captured_desired_destination == 4'hE ? rb2_out : rb3_out;

wire [15:0] increment_value = push ? 16'hFFFF : 16'h0001;
wire do_increment = inc_enable & (push | pop);
wire inc_target = pop ? captured_desired_source : captured_desired_destination;

wire do_register_capture = step == 2'h3 & ~clock;

wire [15:0] r1_out;
register #(.WIDTH(16)) r1 (
    .data_in(write_value),
    .increment(increment_value),
    .capture(do_register_capture),
    .write_enable(write_enable & captured_desired_destination == 4'h5),
    .increment_enable(do_increment & inc_target == 4'h1),
    .reset(reset),
    .value(r1_out)
);

wire [15:0] r2_out;
register #(.WIDTH(16)) r2 (
    .data_in(write_value),
    .increment(increment_value),
    .capture(do_register_capture),
    .write_enable(write_enable & captured_desired_destination == 4'h6),
    .increment_enable(do_increment & inc_target == 4'h2),
    .reset(reset),
    .value(r2_out)
);

wire [15:0] r3_out;
register #(.WIDTH(16)) r3 (
    .data_in(write_value),
    .increment(increment_value),
    .capture(do_register_capture),
    .write_enable(write_enable & captured_desired_destination == 4'h7),
    .increment_enable(do_increment & inc_target == 4'h3),
    .reset(reset),
    .value(r3_out)
);

wire [15:0] rb1_out;
register #(.WIDTH(16)) rb1 (
    .data_in(write_value),
    .increment(increment_value),
    .capture(do_register_capture),
    .write_enable(write_enable & captured_desired_destination == 4'hD),
    .increment_enable(do_increment & inc_target == 4'h9),
    .reset(reset),
    .value(rb1_out)
);

wire [15:0] rb2_out;
register #(.WIDTH(16)) rb2 (
    .data_in(write_value),
    .increment(increment_value),
    .capture(do_register_capture),
    .write_enable(write_enable & captured_desired_destination == 4'hE),
    .increment_enable(do_increment & inc_target == 4'hA),
    .reset(reset),
    .value(rb2_out)
);

wire [15:0] rb3_out;
register #(.WIDTH(16)) rb3 (
    .data_in(write_value),
    .increment(increment_value),
    .capture(do_register_capture),
    .write_enable(write_enable & captured_desired_destination == 4'hF),
    .increment_enable(do_increment & inc_target == 4'hB),
    .reset(reset),
    .value(rb3_out)
);

wire [15:0] pc_out;
register #(.WIDTH(16)) pc_register (
    .data_in(write_value),
    .increment(16'h2),
    .capture(do_register_capture),
    .write_enable(write_enable & captured_desired_destination == 4'h0),
    .increment_enable(1'h1),
    .reset(reset),
    .value(pc_out)
);

wire [15:0] flags_out;
register #(.WIDTH(16)) flags_register (
    .data_in(write_flags ? flags_in : write_value),
    .increment(16'h0),
    .capture(do_register_capture),
    .write_enable(write_flags | (write_enable & captured_desired_destination == 4'h8)),
    .increment_enable(1'h0),
    .reset(reset),
    .value(flags_out)
);

wire [15:0] banking_out;
register #(.WIDTH(16), .INIT(16'h00E0)) banking_register (
    .data_in(write_value),
    .increment(16'h0),
    .capture(do_register_capture),
    .write_enable(write_enable & captured_desired_destination == 4'h4),
    .increment_enable(1'h0),
    .reset(reset),
    .value(banking_out)
);

wire [15:0] interrupt_out;
register #(.WIDTH(16)) interrupt_register (
    .data_in(write_value),
    .increment(16'h0),
    .capture(do_register_capture),
    .write_enable(write_enable & captured_desired_destination == 4'hC),
    .increment_enable(1'h0),
    .reset(reset),
    .value(interrupt_out)
);

endmodule
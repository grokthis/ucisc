module register_block (
    input clock,
    input reset,

    // Read arguments
    input [3:0] source,
    input [15:0] immediate,
    input [3:0] dest,
    input [15:0] offset,

    // Write arguments
    input [3:0] to_write,
    input write_enable,
    input push,
    input pop,
    input [15:0] data_in,

    // Read out
    output [15:0] source_out,
    output [15:0] dest_out,

    // Control registers in
    input [15:0] flags_in,
    input write_flags,

    // Control registers out
    output [15:0] pc,
    output [15:0] banking,
    output [15:0] flags,
    output [15:0] interrupt
);

  wire [15:0][10] registers_out;

  assign pc = registers_out[0];
  assign banking = registers_out[4];
  assign flags = registers_out[8];
  assign interrupt = registers_out[9];

  wire [15:0] inc_amt =
    push ? 16'hFFFF :
    pop ? 16'h0001 :
    16'h0000;

  wire [15:0] write_value = write_enable ? data_in : to_write_out + inc_amt;
  wire [15:0][10] d_in = {
    to_write == 4'h0 && write_enable ? data_in : pc + 2,
    write_value,
    write_value,
    write_value,
    write_value,
    write_value,
    write_value,
    write_value,
    to_write == 4'h8 && write_enable ? data_in : flags_in,
    write_value
  }


  dff #(.WIDTH(16), .INIT(16'h0)) registers[10] (
      .clock(clock),
      .d(d_in),
      .async_reset(reset),
      .enable((inc_r1 | store_r1) && step == 2'h0),
      .q(registers_out)
  );

endmodule

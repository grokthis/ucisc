// program counter

module processor_core(
  input clock,
  input reset,
  input [7:0] local_page_in,
  input [3:0] page_offset,
  input [255:0] page_data,
  input write_page_in,
  input write_page_out_accepted,
  output write_page_out_ready,
  output [7:0] local_page_out,
  output [31:0] main_page_out,
  output [255:0] page_data_out
);
wire [15:0] current_pc;
wire [15:0] next_pc;
wire [15:0] instruction;
wire [15:0] alu_result;
wire [15:0] src_reg_value;
wire [15:0] dest_reg_value;
wire [15:0] src_mem_value;
wire [15:0] dest_mem_value;
wire [2:0] effect;
wire [4:0] alu_code;
wire [6:0] immediate;
wire [2:0] source_select;
wire [2:0] destination_select;
wire source_immediate;
wire source_memory;
wire pre_increment;
wire post_increment;
wire decrement;
wire destination_pc;
wire destination_mem;
wire destination_reg;
wire store_enabled;
wire set_flags;
wire [15:0] alu_a = source_memory ? src_mem_value : src_reg_value;
wire [15:0] alu_b = destination_pc ? current_pc :
  destination_mem ? dest_mem_value : dest_reg_value;

wire processor_clock = clock | neg_override_clock;
reg neg_override_clock;
always @(posedge clock) begin
  if (write_page_in)
    neg_override_clock = 1;
  else
    neg_override_clock = 0;
end

wire mem_clock = clock;
wire local_mem_write = ~processor_clock & destination_mem & store_enabled;

memory_block_64kx16 local_memory (
  .clock(mem_clock),
  .write_enable(local_mem_write),
  .portA_address(next_pc),
  .portB_address(write_page_in ? {local_page_in, page_offset, 4'h0} : dest_reg_value),
  .data_in(alu_result),
  .row_data(page_data),
  .row_write(write_page_in),
  .row_data_out(page_data_out),
  .portA_out(instruction),
  .portB_out(dest_mem_value)
);

pc pc(
  .clock(processor_clock),
  .reset(reset),
  .store_enabled(store_enabled),
  .destination_pc(destination_pc),
  .source_value(alu_a),
  .current_pc(current_pc),
  .next_pc(next_pc)
);

instruction_decoder instruction_decoder(
  .instruction(instruction),
  .source_immediate(source_immediate),
  .source_memory(source_memory),
  .immediate(immediate),
  .alu_code(alu_code),
  .pre_increment(pre_increment),
  .post_increment(post_increment),
  .decrement(decrement),
  .source_select(source_select),
  .destination_select(destination_select),
  .destination_pc(destination_pc),
  .destination_mem(destination_mem),
  .destination_reg(destination_reg),
  .effect(effect),
  .set_flags(set_flags)
);

register_block register_block(
  .clock(processor_clock),
  .pc(current_pc),
  .source_immediate(source_immediate),
  .immediate(immediate),
  .destination_write(alu_result),
  .destination_select(destination_select),
  .source_select(source_select), 
  .flags_value(flags_result),
  .set_flags(set_flags),
  .store_value(reg_destination_write),
  .pre_increment(pre_increment),
  .post_increment(post_increment),
  .decrement(decrement),
  .source_out(src_reg_value),
  .destination_out(dest_reg_value),
  .flags_out(last_flags)
);

wire [15:0] last_flags;
wire last_overflow = last_flags[0];
wire last_carry = last_flags[1];
wire last_zero = last_flags[2];
wire last_negative = last_flags[3];
wire carry_result, overflow_result, zero_result, negative_result;
wire [15:0] flags_result = {negative_result, zero_result, carry_result, overflow_result};

alu alu(
  .clock(processor_clock),
  .a(alu_a),
  .b(alu_b),
  .alu_code(alu_code),
  .carry_in(last_carry),
  .result_out(alu_result),
  .carry_out(carry_result),
  .overflow_out(overflow_result),
  .zero_out(zero_result),
  .negative_out(negative_result)
);


wire reg_destination_write = destination_reg & store_enabled;

effect_decoder effect_decoder(
  .flags(last_flags[4:0]),
  .effect(effect),
  .store(store_enabled)
);

endmodule

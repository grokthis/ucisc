// program counter

module processor_core(
  input clock,
  input reset,
  output [15:0] current_pc,
  output [15:0] instruction
);

//wire [15:0] current_pc;
wire [15:0] next_pc;
//wire [15:0] instruction;
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

reg processor_clock;
reg mem_write_clock;
always @(negedge clock) begin
  mem_write_clock = ~mem_write_clock;
end
always @(posedge clock) begin
  processor_clock = mem_write_clock;
end

wire local_mem_write = ~processor_clock & destination_mem & store_enabled;

memory_block #(.address_width(12)) local_memory (
  .clock(clock),
  .write_clock(mem_write_clock),
  .write_enable(local_mem_write),
  .portA_address(next_pc),
  .portB_address(dest_reg_value),
  .data_in(alu_result),
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

wire carry_result, overflow_result, zero_result, negative_result;
wire [15:0] flags_result = {12'b0, negative_result, zero_result, carry_result, overflow_result};

alu alu(
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

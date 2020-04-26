module instruction_decoder(
  input [15:0] instruction,
  output source_immediate,
  output source_memory,
  output [6:0] immediate,
  output [4:0] alu_code,
  output pre_increment,
  output post_increment,
  output decrement,
  output [2:0] source_select,
  output [2:0] destination_select,
  output destination_pc,
  output destination_mem,
  output destination_reg,
  output [2:0] effect,
  output set_flags
);

  wire copy_instruction = ~instruction[15];

  wire copy_source_immediate;
  wire copy_source_memory;
  wire [6:0] copy_immediate;
  wire [4:0] copy_alu_code;
  wire copy_pre_increment;
  wire copy_post_increment;
  wire copy_decrement;
  wire [2:0] copy_source_select;
  wire [2:0] copy_destination_select;
  wire copy_destination_pc;
  wire copy_destination_mem;
  wire copy_destination_reg;
  wire [2:0] copy_effect;

  copy_decoder copy_decoder (
    .instruction (instruction),
    .source_immediate (copy_source_immediate),
    .source_memory (copy_source_memory),
    .immediate (copy_immediate),
    .alu_code (copy_alu_code),
    .pre_increment (copy_pre_increment),
    .post_increment (copy_post_increment),
    .decrement (copy_decrement),
    .source_select (copy_source_select),
    .destination_select (copy_destination_select),
    .destination_pc (copy_destination_pc),
    .destination_mem (copy_destination_mem),
    .destination_reg (copy_destination_reg),
    .effect (copy_effect)
  );

  wire alu_source_immediate;
  wire alu_source_memory;
  wire [6:0] alu_immediate;
  wire [4:0] alu_alu_code;
  wire alu_pre_increment;
  wire alu_post_increment;
  wire alu_decrement;
  wire [2:0] alu_source_select;
  wire [2:0] alu_destination_select;
  wire alu_destination_pc;
  wire alu_destination_mem;
  wire alu_destination_reg;
  wire [2:0] alu_effect;

  alu_decoder alu_decoder (
    .instruction (instruction),
    .source_immediate (alu_source_immediate),
    .source_memory (alu_source_memory),
    .immediate (alu_immediate),
    .alu_code (alu_alu_code),
    .pre_increment (alu_pre_increment),
    .post_increment (alu_post_increment),
    .decrement (alu_decrement),
    .source_select (alu_source_select),
    .destination_select (alu_destination_select),
    .destination_pc (alu_destination_pc),
    .destination_mem (alu_destination_mem),
    .destination_reg (alu_destination_reg),
    .effect (alu_effect)
  );

  assign source_immediate = copy_instruction ? copy_source_immediate : alu_source_immediate;
  assign source_memory = copy_instruction ? copy_source_memory : alu_source_memory;
  assign alu_code = copy_instruction ? copy_alu_code : alu_alu_code;
  assign pre_increment = copy_instruction ? copy_pre_increment : alu_pre_increment;
  assign post_increment = copy_instruction ? copy_post_increment : alu_post_increment;
  assign decrement = copy_instruction ? copy_decrement : alu_decrement;
  assign source_select = copy_instruction ? copy_source_select : alu_source_select;
  assign destination_select = copy_instruction ? copy_destination_select : alu_destination_select;
  assign destination_pc = copy_instruction ? copy_destination_pc : alu_destination_pc;
  assign destination_mem = copy_instruction ? copy_destination_mem : alu_destination_mem;
  assign destination_reg = copy_instruction ? copy_destination_reg : alu_destination_reg;
  assign immediate = copy_instruction ? copy_immediate : alu_immedaite;
  assign effect = copy_instruction ? copy_effect : alu_effect;
  assign set_flags = alu_instruction;

endmodule

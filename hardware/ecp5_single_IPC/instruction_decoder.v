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
  wire alu_instruction = instruction[15] & ~instruction[14];
  wire page_instruction = instruction[15] & instruction[14] & ~instruction[13];


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
  wire copy_set_flags;

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
  wire alu_set_flags;

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

  wire page_source_immediate;
  wire page_source_memory;
  wire [6:0] page_immediate;
  wire [4:0] page_alu_code;
  wire page_pre_increment;
  wire page_post_increment;
  wire page_decrement;
  wire [2:0] page_source_select;
  wire [2:0] page_destination_select;
  wire page_destination_pc;
  wire page_destination_mem;
  wire page_destination_reg;
  wire [2:0] page_effect;
  wire page_set_flags;

  page_decoder page_decoder (
    .instruction (instruction),
    .source_immediate (page_source_immediate),
    .source_memory (page_source_memory),
    .immediate (page_immediate),
    .alu_code (page_alu_code),
    .pre_increment (page_pre_increment),
    .post_increment (page_post_increment),
    .decrement (page_decrement),
    .source_select (page_source_select),
    .destination_select (page_destination_select),
    .destination_pc (page_destination_pc),
    .destination_mem (page_destination_mem),
    .destination_reg (page_destination_reg),
    .effect (page_effect)
  );

  assign source_immediate =
    (copy_instruction) ? copy_source_immediate:
    (alu_instruction) ? alu_source_immediate:
    (page_instruction) ? page_source_immediate:
    1;

  assign source_memory =
    (copy_instruction) ? copy_source_memory:
    (alu_instruction) ? alu_source_memory:
    (page_instruction) ? page_source_memory:
    0;

  assign alu_code =
    (copy_instruction) ? copy_alu_code:
    (alu_instruction) ? alu_alu_code:
    (page_instruction) ? page_alu_code:
    5'h0;

  assign pre_increment =
    (copy_instruction) ? copy_pre_increment:
    (alu_instruction) ? alu_pre_increment:
    (page_instruction) ? page_pre_increment:
    0;

  assign post_increment =
    (copy_instruction) ? copy_post_increment:
    (alu_instruction) ? alu_post_increment:
    (page_instruction) ? page_post_increment:
    0;

  assign decrement =
    (copy_instruction) ? copy_decrement:
    (alu_instruction) ? alu_decrement:
    (page_instruction) ? page_decrement:
    0;

  assign source_select =
    (copy_instruction) ? copy_source_select:
    (alu_instruction) ? alu_source_select:
    (page_instruction) ? page_source_select:
    3'h0;

  assign destination_select =
    (copy_instruction) ? copy_destination_select:
    (alu_instruction) ? alu_destination_select:
    (page_instruction) ? page_destination_select:
    3'h0;

  assign destination_pc =
    (copy_instruction) ? copy_destination_pc:
    (alu_instruction) ? alu_destination_pc:
    (page_instruction) ? page_destination_pc:
    0;

  assign destination_mem =
    (copy_instruction) ? copy_destination_mem:
    (alu_instruction) ? alu_destination_mem:
    (page_instruction) ? page_destination_mem:
    0;

  assign destination_reg =
    (copy_instruction) ? copy_destination_reg:
    (alu_instruction) ? alu_destination_reg:
    (page_instruction) ? page_destination_reg:
    0;

  assign immediate =
    (copy_instruction) ? copy_immediate:
    (alu_instruction) ? alu_immediate:
    (page_instruction) ? page_immediate:
    7'h0;

  assign effect =
    (copy_instruction) ? copy_effect:
    (alu_instruction) ? alu_effect:
    (page_instruction) ? page_effect:
    3'h0;

  assign set_flags = alu_instruction;

endmodule

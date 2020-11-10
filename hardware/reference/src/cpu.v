module cpu (
  input clock_input,
  input reset,
  output [1:0] step,
  output [15:0] r1_peek,
  output [15:0] pc_peek,
  output [15:0] mem_read_address_peek,
  output [15:0] memory_result_peek,
  output store_peek,
  output destination_mem_peek,
  output [15:0] source_value_peek,
  output [15:0] destination_value_peek,
  output [15:0] result_peek,
  output [15:0] immediate_peek,
  output push_peek,
  output pop_peek,
  output tx
);

    parameter MEM_INIT_FILE = "prog.hex";

    // peek vars are for debugging
    assign pc_peek = pc;
    assign mem_read_address_peek = mem_read_address;
    assign memory_result_peek = memory_result;
    assign source_value_peek = source_value;
    assign destination_value_peek = destination_value;
    assign result_peek = alu_result;
    assign immediate_peek = immediate;
    assign push_peek = push;
    assign pop_peek = pop;

    //wire clock = ~(halted | clock_input);
    wire clock = clock_input;

    wire [1:0] step;
    dff #(.WIDTH(2), .INIT(2'h3)) step_ff (
        .clock(clock),
        .d(step + 2'h1),
        .async_reset(reset),
        .enable(1'h1),
        .q(step)
    );

    wire [15:0] mem_read_address =
        step == 2'h0 ? immediate_read_address :
        step == 2'h1 ? source_read_address :
        step == 2'h2 ? destination_read_address :
        instruction_read_address;

    wire destination_mem;
    wire store;
    wire [15:0] memory_result;
    wire [15:0] alu_result;
    wire [15:0] destination_value;

    memory_block #(.WIDTH(13), .MEM_INIT_FILE(MEM_INIT_FILE)) memory_block (
        .clock(clock),
        .step(step),
        .write_enable(store & destination_mem),
        .read_address(mem_read_address),
        .write_address(destination_value),
        .data_in(alu_result),
        .data_out(memory_result)
    );

    wire [3:0] desired_source;
    wire [3:0] desired_destination;
    wire [15:0] alu_flags_out;
    wire write_flags;
    wire [15:0] source_value;
    wire [15:0] pc;
    wire [15:0] flags;
    wire [15:0] banking;
    wire push;
    wire pop;
    wire source_banked;
    wire destination_banked;

    register_block register_block (
        .clock(clock),
        .step(step),
        .reset(reset),
        .desired_source(desired_source),
        .desired_destination(desired_destination),
        .write_value(alu_result),
        .write_enable(store & ~destination_mem),
        .push(push),
        .pop(pop),
        .inc_enable(store),
        .flags_in(alu_flags_out),
        .write_flags(write_flags),
        .source_value(source_value),
        .destination_value(destination_value),
        .pc(pc),
        .flags(flags),
        .banking(banking),
        .source_banked(source_banked),
        .destination_banked(destination_banked),
        .r1_peek(r1_peek)
    );

    wire [3:0] op;
    wire [2:0] effect;
    wire source_mem;
    wire instruction_could_halt;
    wire [15:0] instruction_read_address;

    instruction_loader instruction_loader (
        .clock(clock),
        .step(step),
        .memory_in(memory_result),
        .pc(pc),
        .op(op),
        .effect(effect),
        .destination(desired_destination),
        .source(desired_source),
        .push(push),
        .pop(pop),
        .destination_mem(destination_mem),
        .source_mem(source_mem),
        .could_halt(instruction_could_halt),
        .read_address(instruction_read_address)
    );

    wire [15:0] immediate_read_address;
    wire immediate_could_halt;
    wire [15:0] immediate;
    wire [3:0] offset;

    immediate_loader immediate_loader (
        .clock(clock),
        .step(step),
        .memory_in(memory_result),
        .pc(pc),
        .destination_mem(destination_mem),
        .read_address(immediate_read_address),
        .immediate(immediate),
        .offset(offset),
        .can_halt(immediate_could_halt)
    );

    wire [15:0] source_read_address;
    wire [15:0] captured_source_value;
    wire [15:0] banked_result;

    arg_loader source_loader (
        .clock(clock),
        .step(step),
        .capture_on(2'h2),
        .immediate(immediate),
        .register_value(source_value),
        .memory_in(source_banked ? banked_result : memory_result),
        .is_mem(source_mem),
        .read_address(source_read_address),
        .source_value(captured_source_value)
    );

    wire [15:0] destination_read_address;
    wire [15:0] captured_destination_value;

    arg_loader destination_loader (
        .clock(clock),
        .step(step),
        .capture_on(2'h3),
        .immediate({12'h0, offset}),
        .register_value(destination_value),
        .memory_in(destination_banked ? banked_result : memory_result),
        .is_mem(destination_mem),
        .read_address(destination_read_address),
        .source_value(captured_destination_value)
    );

    alu alu (
        .source(captured_source_value),
        .destination(captured_destination_value),
        .op_code(op),
        .flags(flags),
        .result_out(alu_result),
        .flags_out(alu_flags_out),
        .write_flags(write_flags)
    );

    effect_decoder effect_decoder (
        .flags(flags),
        .effect(effect),
        .store(store)
    );

    wire halted;

    dff #(.WIDTH(1)) halt_dff (
        .clock(clock_input),
        .d(step == 2'h3 & instruction_could_halt & immediate_could_halt & op == 4'h0 & store),
        .async_reset(1'h0),
        .enable(1'h1),
        .q(halted)
    );

    wire is_uart_control_address = (mem_read_address[15:12] == 4'h0 & mem_read_address[11:4] == 8'h2);
    assign banked_result = is_uart_control_address ? uart_read : 16'hA;

    wire [15:0] uart_out;
    wire tx;
    wire [15:0] uart_read =
        mem_read_address[3:0] == 4'h0 ? 16'h4 : uart_out;

    uart_device uart_device(
        .clock(clock),
        .control_address(mem_read_address[3:0]),
        .control_write(store & destination_banked & is_uart_control_address),
        .data_in(alu_result),
        .control_read(uart_out),
        .Tx(tx)
    );

endmodule
module cpu (
  input clock_input,
  input reset,
  output [1:0] step,
  output [15:0] r1_peek,
  output [15:0] pc_peek,
  output tx
);

    parameter MEM_INIT_FILE = "prog.hex";

    // peek vars are for debugging
    assign pc_peek = pc;
    assign r1_peek = r1;

    //wire clock = ~(halted | clock_input);
    wire clock = clock_input;

    wire [1:0] step;
    dff #(.WIDTH(2), .INIT(2'h0)) step_ff (
        .clock(clock),
        .d(step + 2'h1),
        .async_reset(reset),
        .enable(1'h1),
        .q(step)
    );

    wire st_en;

    wire store_pc = st_en && dst == 4'h0;
    wire [15:0] pc_out;
    wire [15:0] pc_result = store_pc ? result : pc_out + 2'h2;
    wire [15:0] pc = step == 2'h0 ? pc_result : pc_out;
    dff #(.WIDTH(16), .INIT(16'hFFFE)) pc_register (
        .clock(clock),
        .d(pc_result),
        .async_reset(reset),
        .enable(step == 2'h0),
        .q(pc_out)
    );

    wire [15:0] result;

    wire [15:0] inc_amt = push ? 16'hFFFF : 16'h0001;

    wire [15:0] r1;
    wire store_r1 = st_en && dst == 4'h5;
    wire inc_r1 = (push & dst == 4'h1) | (pop & src == 4'h1);
    dff #(.WIDTH(16), .INIT(16'h0)) r1_register (
        .clock(clock),
        .d(store_r1 ? result : r1 + inc_amt),
        .async_reset(reset),
        .enable((inc_r1 | store_r1) && step == 2'h0),
        .q(r1)
    );

    wire [15:0] r2;
    wire store_r2 = st_en && dst == 4'h6;
    wire inc_r2 = (push & dst == 4'h2) | (pop & src == 4'h2);
    dff #(.WIDTH(16), .INIT(16'h0)) r2_register (
        .clock(clock),
        .d(store_r2 ? result : r2 + inc_amt),
        .async_reset(reset),
        .enable((inc_r2 | store_r2) && step == 2'h0),
        .q(r2)
    );

    wire [15:0] r3;
    wire store_r3 = st_en && dst == 4'h7;
    wire inc_r3 = (push & dst == 4'h3) | (pop & src == 4'h3);
    dff #(.WIDTH(16), .INIT(16'h0)) r3_register (
        .clock(clock),
        .d(store_r3 ? result : r3 + inc_amt),
        .async_reset(reset),
        .enable((inc_r3 | store_r3) && step == 2'h0),
        .q(r3)
    );

    wire [15:0] r4;
    wire store_r4 = st_en && dst == 4'hD;
    wire inc_r4 = (push & dst == 4'h9) | (pop & src == 4'h9);
    dff #(.WIDTH(16), .INIT(16'h0)) r4_register (
        .clock(clock),
        .d(store_r4 ? result : r4 + inc_amt),
        .async_reset(reset),
        .enable((inc_r4 | store_r4) && step == 2'h0),
        .q(r4)
    );

    wire [15:0] r5;
    wire store_r5 = st_en && dst == 4'hE;
    wire inc_r5 = (push & dst == 4'hA) | (pop & src == 4'hA);
    dff #(.WIDTH(16), .INIT(16'h0)) r5_register (
        .clock(clock),
        .d(store_r5 ? result : r5 + inc_amt),
        .async_reset(reset),
        .enable((inc_r5 | store_r5) && step == 2'h0),
        .q(r5)
    );

    wire [15:0] r6;
    wire store_r6 = st_en && dst == 4'hF;
    wire inc_r6 = (push & dst == 4'hB) | (pop & src == 4'hB);
    dff #(.WIDTH(16), .INIT(16'h0)) r6_register (
        .clock(clock),
        .d(store_r6 ? result : r6 + inc_amt),
        .async_reset(reset),
        .enable((inc_r6 | store_r6) && step == 2'h0),
        .q(r6)
    );

    wire [15:0] memory_result;

    // Memory source address needs to be setup on step 2, so we short
    // cut the ir value on step 2 since it wouldn't otherwise be ready
    // until step 3
    wire [31:0] ir = step == 2'h2 ? {_ir[31:16], memory_result} : _ir;
    wire [31:0] _ir;
    dff #(.WIDTH(32), .INIT(32'h0)) instruction_register (
        .clock(clock),
        .d(step == 2'h1 ? {memory_result, ir[15:0]} : {ir[31:16], memory_result}),
        .async_reset(reset),
        .enable(step == 2'h1 | step == 2'h2),
        .q(_ir)
    );
    wire [3:0] src = ir[31:28];
    wire src_mem = (src[0] | src[1]) & ~src[2];
    wire [3:0] dst = ir[27:24];
    wire dst_mem = (dst[0] | dst[1]) & ~dst[2];
    wire inc = ir[23];
    wire push = dst_mem & inc;
    wire pop = src_mem & inc & ~dst_mem;
    wire [2:0] eff = ir[22:20];
    wire [3:0] op = ir[19:16];
    wire [15:0] imm = dst_mem ? { {4{ir[11]}}, ir[11:0] } : ir[15:0];
    wire [3:0] off = dst_mem ? ir[15:12] : 4'h0;


    wire [15:0] src_addr =
        src == 4'h0 ? pc + imm :
        src == 4'h1 ? r1 + imm :
        src == 4'h2 ? r2 + imm :
        src == 4'h3 ? r3 + imm :
        src == 4'h4 ? imm :
        src == 4'h5 ? r1 + imm :
        src == 4'h6 ? r2 + imm :
        src == 4'h7 ? r3 + imm :
        src == 4'h8 ? 16'h0 :
        src == 4'h9 ? r4 + imm :
        src == 4'hA ? r5 + imm :
        src == 4'hB ? r6 + imm :
        src == 4'hC ? 16'h0 :
        src == 4'hD ? r4 + imm :
        src == 4'hE ? r5 + imm :
        r6 + imm;

    wire [15:0] _dst_addr =
        dst == 4'h0 ? pc :
        dst == 4'h1 ? r1 + off :
        dst == 4'h2 ? r2 + off :
        dst == 4'h3 ? r3 + off :
        dst == 4'h4 ? 16'h0 :
        dst == 4'h5 ? r1 :
        dst == 4'h6 ? r2 :
        dst == 4'h7 ? r3 :
        dst == 4'h8 ? 16'h0 :
        dst == 4'h9 ? r4 + off :
        dst == 4'hA ? r5 + off :
        dst == 4'hB ? r6 + off :
        dst == 4'hC ? 16'h0 :
        dst == 4'hD ? r4 :
        dst == 4'hE ? r5 :
        r6;

    // Without register, this creates a wire loop with the PC
    // We don't capture the dest value before handing it to store
    // which means pc could be the source and dest and create a
    // wire loop.
    wire [15:0] dst_addr;
    dff #(.WIDTH(16), .INIT(16'h0)) dst_addr_register (
        .clock(clock),
        .d(_dst_addr),
        .async_reset(reset),
        .enable(step == 2'h2),
        .q(dst_addr)
    );

    wire [15:0] mem_read_address =
        step == 2'h0 ? pc :
        step == 2'h1 ? pc + 1 :
        step == 2'h2 ? src_addr :
        dst_addr;

    memory_block #(.WIDTH(13), .MEM_INIT_FILE(MEM_INIT_FILE)) memory_block (
        .clock(clock),
        .read_address(mem_read_address),
        .write_enable(step == 2'h0 & st_en & dst_mem),
        .write_address(push ? dst_addr - 1'b1 : dst_addr),
        .data_in(result),
        .data_out(memory_result)
    );

     wire [15:0] src_val;
     dff #(.WIDTH(16), .INIT(16'h0)) src_val_register (
         .clock(clock),
         .d(src_mem ? memory_result : src_addr),
         .async_reset(reset),
         .enable(step == 2'h3),
         .q(src_val)
     );

     wire [15:0] dst_val = step == 2'h0 & dst_mem ? memory_result : dst_addr;
     wire [15:0] flags;
     wire [15:0] alu_flags;
     wire alu_write_flags;

     alu alu (
         .source(src_val),
         .destination(dst_val),
         .op_code(op),
         .flags(flags),
         .result_out(result),
         .flags_out(alu_flags),
         .write_flags(alu_write_flags)
     );

     effect_decoder effect_decoder (
         .flags(flags),
         .effect(eff),
         .store(st_en)
     );

     wire store_flags = alu_write_flags | (st_en && dst == 4'hF);
     dff #(.WIDTH(16), .INIT(16'h0100)) flags_register (
         .clock(clock),
         .d(alu_write_flags ? alu_flags : result),
         .async_reset(reset),
         .enable(store_flags && step == 2'h0),
         .q(flags)
     );
endmodule
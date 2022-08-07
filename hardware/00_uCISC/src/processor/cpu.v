module cpu (
  input cpu_clock,
  input reset,
  output device_write_en,
  output [15:0] device_address,
  input [15:0] device_data_in,
  output [15:0] device_data_out,
  output [2:0] step,
);
    parameter MEM_ADDRESS_WIDTH = 16;

    wire [2:0] step;
    dff #(.WIDTH(3), .INIT(3'h0)) step_ff (
        .clock(cpu_clock),
        .d(step == 3'h5 ? 3'h0 : step + 3'h1),
        .async_reset(reset),
        .enable(1'h1),
        .q(step)
    );

    wire step0;
    dff #(.WIDTH(1), .INIT(1'h1)) step0_register (
        .clock(cpu_clock),
        .d(step == 3'h5),
        .async_reset(reset),
        .enable(1'h1),
        .q(step0)
    );

    wire step1;
    dff #(.WIDTH(1), .INIT(1'h0)) step1_register (
        .clock(cpu_clock),
        .d(step == 3'h0),
        .async_reset(reset),
        .enable(1'h1),
        .q(step1)
    );

    wire step2;
    dff #(.WIDTH(1), .INIT(1'h0)) step2_register (
        .clock(cpu_clock),
        .d(step == 3'h1),
        .async_reset(reset),
        .enable(1'h1),
        .q(step2)
    );

    wire step3;
    dff #(.WIDTH(1), .INIT(1'h0)) step3_register (
        .clock(cpu_clock),
        .d(step == 3'h2),
        .async_reset(reset),
        .enable(1'h1),
        .q(step3)
    );

    wire step4;
    dff #(.WIDTH(1), .INIT(1'h0)) step4_register (
        .clock(cpu_clock),
        .d(step == 3'h3),
        .async_reset(reset),
        .enable(1'h1),
        .q(step4)
    );

    wire step5;
    dff #(.WIDTH(1), .INIT(1'h0)) step5_register (
        .clock(cpu_clock),
        .d(step == 3'h4),
        .async_reset(reset),
        .enable(1'h1),
        .q(step5)
    );

//=================================================
//STAGE 0 - LOAD INSTRUCTION
//=================================================

    wire [15:0] result;

    wire store_pc;
    dff #(.WIDTH(1), .INIT(1'h0)) store_pc_register (
        .clock(cpu_clock),
        .d(st_en && dst == 4'h0),
        .async_reset(reset),
        .enable(step2),
        .q(store_pc)
    );

    wire [15:0] pc;
    dff #(.WIDTH(16), .INIT(16'hFFFE)) pc_register (
        .clock(cpu_clock),
        .d(store_pc ? result : pc + 2'h2),
        .async_reset(reset),
        .enable(step0),
        .q(pc)
    );

    wire [15:0] r1;
    wire store_r1 = st_en && dst == 4'h5;
    wire inc_r1 = (push & dst == 4'h1) | (pop & src == 4'h1);
    dff #(.WIDTH(16), .INIT(16'h0)) r1_register (
        .clock(cpu_clock),
        .d(store_r1 ? result : r1 + inc_amt),
        .async_reset(reset),
        .enable((inc_r1 | store_r1) && step0),
        .q(r1)
    );

    wire [15:0] r2;
    wire store_r2 = st_en && dst == 4'h6;
    wire inc_r2 = (push & dst == 4'h2) | (pop & src == 4'h2);
    dff #(.WIDTH(16), .INIT(16'h0)) r2_register (
        .clock(cpu_clock),
        .d(store_r2 ? result : r2 + inc_amt),
        .async_reset(reset),
        .enable((inc_r2 | store_r2) && step0),
        .q(r2)
    );

    wire [15:0] r3;
    wire store_r3 = st_en && dst == 4'h7;
    wire inc_r3 = (push & dst == 4'h3) | (pop & src == 4'h3);
    dff #(.WIDTH(16), .INIT(16'h0)) r3_register (
        .clock(cpu_clock),
        .d(store_r3 ? result : r3 + inc_amt),
        .async_reset(reset),
        .enable((inc_r3 | store_r3) && step0),
        .q(r3)
    );

    wire [15:0] r4;
    wire store_r4 = st_en && dst == 4'hD;
    wire inc_r4 = (push & dst == 4'h9) | (pop & src == 4'h9);
    dff #(.WIDTH(16), .INIT(16'h0)) r4_register (
        .clock(cpu_clock),
        .d(store_r4 ? result : r4 + inc_amt),
        .async_reset(reset),
        .enable((inc_r4 | store_r4) && step0),
        .q(r4)
    );

    wire [15:0] r5;
    wire store_r5 = st_en && dst == 4'hE;
    wire inc_r5 = (push & dst == 4'hA) | (pop & src == 4'hA);
    dff #(.WIDTH(16), .INIT(16'h0)) r5_register (
        .clock(cpu_clock),
        .d(store_r5 ? result : r5 + inc_amt),
        .async_reset(reset),
        .enable((inc_r5 | store_r5) && step0),
        .q(r5)
    );

    wire [15:0] r6;
    wire store_r6 = st_en && dst == 4'hF;
    wire inc_r6 = (push & dst == 4'hB) | (pop & src == 4'hB);
    dff #(.WIDTH(16), .INIT(16'h0)) r6_register (
        .clock(cpu_clock),
        .d(store_r6 ? result : r6 + inc_amt),
        .async_reset(reset),
        .enable((inc_r6 | store_r6) && step0),
        .q(r6)
    );

    wire [15:0] banking;
    wire store_banking = st_en && dst == 4'h4;
    dff #(.WIDTH(16), .INIT(16'h0070)) banking_register (
        .clock(cpu_clock),
        .d(result),
        .async_reset(reset),
        .enable(store_banking && step0),
        .q(banking)
    );

    wire [15:0] flags;
    wire store_flags = alu_write_flags | (st_en && dst == 4'hF);
    dff #(.WIDTH(16), .INIT(16'h0100)) flags_register (
        .clock(cpu_clock),
        .d(alu_write_flags ? alu_flags : result),
        .async_reset(reset),
        .enable(store_flags && step0),
        .q(flags)
    );

    wire [15:0] pc_preload;
    dff #(.WIDTH(16), .INIT(16'h0)) pc_preload_register (
        .clock(cpu_clock),
        .d(pc + 2'h2),
        .async_reset(reset),
        .enable(step1),
        .q(pc_preload)
    );


//=================================================
//STAGE 1 - DECODE
//=================================================

    // src_in == instruction[31:28]
    wire [3:0] src_in = mem_out_select ? mem_result_b[15:12] : mem_result_a[15:12];
    wire src_mem_in = (src_in[0] | src_in[1]) & ~src_in[2];

    // dst_in == instruction[27:24]
    wire [3:0] dst_in = mem_out_select ? mem_result_b[11:8] : mem_result_a[11:8];
    wire dst_mem_in = (dst_in[0] | dst_in[1]) & ~dst_in[2];

    // inc_in == instruction[23]
    wire inc_in = mem_out_select ? mem_result_b[7] : mem_result_a[7];
    wire push_in = dst_mem_in & inc_in;
    wire pop_in = src_mem_in & inc_in & ~dst_mem_in;

    // eff_in == instruction[22:20]
    wire [2:0] eff_in = mem_out_select ? mem_result_b[6:4] : mem_result_a[6:4];

    // op_in == instruction[19:16]
    wire [3:0] op_in = mem_out_select ? mem_result_b[3:0] : mem_result_a[3:0];

    // imm_in == instruction[15:0] or instruction[11:0]
    wire [15:0] imm_in = dst_mem_in ?
        (mem_out_select ?
            { {4{mem_result_a[11]}}, mem_result_a[11:0] } :
            { {4{mem_result_b[11]}}, mem_result_b[11:0] }
        ) :
        (mem_out_select ? mem_result_a : mem_result_b);

    // off_in == instruction[15:12]
    wire [3:0] off_in = mem_out_select ? mem_result_a[15:12] : mem_result_b[15:12];

    wire push;
    dff #(.WIDTH(1), .INIT(1'h0)) push_register (
        .clock(cpu_clock),
        .d(push_in),
        .async_reset(reset),
        .enable(step1),
        .q(push)
    );

    wire pop;
    dff #(.WIDTH(1), .INIT(1'h0)) pop_register (
        .clock(cpu_clock),
        .d(pop_in),
        .async_reset(reset),
        .enable(step1),
        .q(pop)
    );

    wire [15:0] imm;
    dff #(.WIDTH(16), .INIT(16'h0)) imm_register (
        .clock(cpu_clock),
        .d(imm_in),
        .async_reset(reset),
        .enable(step1),
        .q(imm)
    );

    wire [3:0] op;
    dff #(.WIDTH(4), .INIT(4'h0)) op_register (
        .clock(cpu_clock),
        .d(op_in),
        .async_reset(reset),
        .enable(step1),
        .q(op)
    );

    wire [2:0] eff;
    dff #(.WIDTH(3), .INIT(3'h0)) eff_register (
        .clock(cpu_clock),
        .d(eff_in),
        .async_reset(reset),
        .enable(step1),
        .q(eff)
    );

    wire [3:0] src;
    dff #(.WIDTH(4), .INIT(4'h0)) src_register (
        .clock(cpu_clock),
        .d(src_in),
        .async_reset(reset),
        .enable(step1),
        .q(src)
    );

    wire [3:0] off;
    dff #(.WIDTH(4), .INIT(4'h0)) off_register (
        .clock(cpu_clock),
        .d(off_in),
        .async_reset(reset),
        .enable(step1),
        .q(off)
    );

    wire [3:0] dst;
    dff #(.WIDTH(4), .INIT(4'h0)) dst_register (
        .clock(cpu_clock),
        .d(dst_in),
        .async_reset(reset),
        .enable(step1),
        .q(dst)
    );

    wire [2:0] dst_mem_reg = mem_out_select ?
        { mem_result_b[11], mem_result_b[9:8] } :
        { mem_result_a[11], mem_result_a[9:8] };
    wire [15:0] dst_reg_fast;
    dff #(.WIDTH(16), .INIT(16'h0)) dst_fast_register (
        .clock(cpu_clock),
        .d(
            (dst_mem_reg == 3'h1 ? r1 :
            dst_mem_reg == 3'h2 ? r2 :
            dst_mem_reg == 3'h3 ? r3 :
            dst_mem_reg == 3'h5 ? r4 :
            dst_mem_reg == 3'h6 ? r5 :
            r6) + off_in
        ),
        .async_reset(reset),
        .enable(step1),
        .q(dst_reg_fast)
    );

    wire dst_mem;
    dff #(.WIDTH(1), .INIT(1'h0)) dst_mem_register (
        .clock(cpu_clock),
        .d(dst_mem_in),
        .async_reset(reset),
        .enable(step1),
        .q(dst_mem)
    );

    wire st_en_in;
    effect_decoder effect_decoder (
        .flags(flags),
        .effect(eff_in),
        .store(st_en_in)
    );

    wire st_en;
    dff #(.WIDTH(1), .INIT(1'h0)) st_en_register (
        .clock(cpu_clock),
        .d(st_en_in),
        .async_reset(reset),
        .enable(step1),
        .q(st_en)
    );

    wire [15:0] pc_stage2;
    dff #(.WIDTH(16), .INIT(16'h0)) pc_stage2_register (
        .clock(cpu_clock),
        .d(pc),
        .async_reset(reset),
        .enable(step1),
        .q(pc_stage2)
    );

//=================================================
//STAGE 2 - SRC | DEST | ALU
//=================================================

    wire [15:0] inc_amt_in = push ? 16'hFFFF : imm + 1;
    wire [15:0] inc_amt;
    dff #(.WIDTH(16), .INIT(16'h0)) inc_amt_register (
        .clock(cpu_clock),
        .d(inc_amt_in),
        .async_reset(reset),
        .enable(step2),
        .q(inc_amt)
    );

    wire [15:0] dst_addr_in =
        dst == 4'h0 ? pc_stage2 :
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

    wire [15:0] dst_addr;
    dff #(.WIDTH(16), .INIT(16'h0)) dst_addr_register (
        .clock(cpu_clock),
        .d(dst_addr_in),
        .async_reset(reset),
        .enable(step2),
        .q(dst_addr)
    );

    wire [15:0] src_addr_in =
        src == 4'h0 ? pc_stage2 + imm :
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

    wire [15:0] src_addr;
    dff #(.WIDTH(16), .INIT(16'h0)) src_addr_register (
        .clock(cpu_clock),
        .d(src_addr_in),
        .async_reset(reset),
        .enable(step2),
        .q(src_addr)
    );

    wire [15:0] src_addr_fast;
    dff #(.WIDTH(16), .INIT(16'h0)) src_addr_fast_register (
        .clock(cpu_clock),
        .d(src_addr_in),
        .async_reset(reset),
        .enable(step2),
        .q(src_addr_fast)
    );

    wire src_mem = (src[0] | src[1]) & ~src[2];
    wire src_dev_in = src_mem & (
        (src == 4'h1 & banking[0]) |
        (src == 4'h2 & banking[1]) |
        (src == 4'h3 & banking[2]) |
        (src == 4'h9 & banking[4]) |
        (src == 4'hA & banking[5]) |
        (src == 4'hB & banking[6])
    );

    wire src_dev;
    dff #(.WIDTH(1), .INIT(1'h0)) src_dev_register (
        .clock(cpu_clock),
        .d(src_dev_in),
        .async_reset(reset),
        .enable(step2),
        .q(src_dev)
    );

    wire dst_dev_in = dst_mem & (
        (dst == 4'h1 & banking[0]) |
        (dst == 4'h2 & banking[1]) |
        (dst == 4'h3 & banking[2]) |
        (dst == 4'h9 & banking[4]) |
        (dst == 4'hA & banking[5]) |
        (dst == 4'hB & banking[6])
    );

    wire dst_dev;
    dff #(.WIDTH(1), .INIT(1'h0)) dst_dev_register (
        .clock(cpu_clock),
        .d(dst_dev_in),
        .async_reset(reset),
        .enable(step2),
        .q(dst_dev)
    );

    wire [15:0] mem_read_address =
        step0 ? pc_next :
        step1 ? pc_preload :
        step2 ? dst_reg_fast :
        step3 ? src_addr_fast :
        src_addr; // step4

    wire [15:0] mem_result_a;
    wire [15:0] mem_result_b;
    wire mem_out_select;
    memory_block #(
        .WIDTH(MEM_ADDRESS_WIDTH)
    ) memory_block (
      .clock(cpu_clock),
      .read_address(mem_read_address),
      .write_enable(step0 & st_en & dst_mem & ~dst_dev),
      .write_address(push ? dst_addr - 1'b1 : dst_addr),
      .data_in(result),
      .data_out_a(mem_result_a),
      .data_out_b(mem_result_b),
      .primary_out_select(mem_out_select)
    );

    assign device_write_en = step0 & st_en & dst_dev;
    assign device_address = step0 ? dst_addr : mem_read_address;
    assign device_data_out = result;
    wire [15:0] device_result;
    assign device_result = device_data_in;

//    wire [15:0] device_peek;
//    devices devices(
//        .ref_clock(clk),
//        .cpu_clock(cpu_clock),
//        .write_enable(step0 & st_en & dst_dev),
//        .address(step0 ? dst_addr : mem_read_address),
//        .data_in(result),
//        .data_out(device_result),
//        .tx(tx),
//        .rx(rx),
//        .pixel_out(pixel_out),
//        .h_sync_signal(h_sync_signal),
//        .v_sync_signal(v_sync_signal)//,
//        //.peek(device_peek)
//    );

    wire [15:0] src_val;
    dff #(.WIDTH(16), .INIT(16'h0)) src_val_register (
        .clock(cpu_clock),
        .d(
            src_dev ? device_result :
            src_mem & mem_out_select ? mem_result_b :
            src_mem ? mem_result_a :
            src_addr
         ),
        .async_reset(reset),
        .enable(step4),
        .q(src_val)
    );

    wire [15:0] pc_stage3;
    dff #(.WIDTH(16), .INIT(16'h0)) pc_stage3_register (
        .clock(cpu_clock),
        .d(pc_stage2),
        .async_reset(reset),
        .enable(step2),
        .q(pc_stage3)
    );

//=================================================
//STAGE 3 - DEST | ALU
//=================================================

    wire [15:0] dst_val;
    dff #(.WIDTH(16), .INIT(16'h0)) dst_val_register (
        .clock(cpu_clock),
        .d(
            dst_dev ? device_result :
            dst_mem & mem_out_select ? mem_result_b :
            dst_mem ? mem_result_a :
            dst_addr
        ),
        .async_reset(reset),
        .enable(step3),
        .q(dst_val)
    );

//=================================================
//STAGE 4 - ALU
//=================================================

    wire [15:0] result_in;
    wire [1:0] alu_oc_in;
    wire overflow;
    wire carry;
    alu alu (
        .source(src_val),
        .destination(dst_val),
        .op_code(op),
        .flags(flags),
        .result_out(result_in),
        .overflow(overflow),
        .carry(carry)
    );

    wire [1:0] oc_out;
    dff #(.WIDTH(2), .INIT(2'h0)) alu_oc_register (
        .clock(cpu_clock),
        .d({ overflow, carry }),
        .async_reset(reset),
        .enable(step5),
        .q(oc_out)
    );
    wire [15:0] alu_flags = { flags[15:5], 1'b0, oc_out[1], oc_out[0], result[15], result == 16'h0 };

    wire alu_write_flags;
    dff #(.WIDTH(1), .INIT(1'h0)) write_flags_register (
        .clock(cpu_clock),
        .d(op != 4'h0),
        .async_reset(reset),
        .enable(step5),
        .q(alu_write_flags)
    );

    dff #(.WIDTH(16), .INIT(16'h0)) result_register (
        .clock(cpu_clock),
        .d(result_in),
        .async_reset(reset),
        .enable(step5),
        .q(result)
    );

    wire [15:0] pc_next;
    dff #(.WIDTH(16), .INIT(16'h0)) pc_next_register (
        .clock(cpu_clock),
        .d(store_pc ? result_in : pc_stage3 + 2'h2),
        .async_reset(reset),
        .enable(step5),
        .q(pc_next)
    );

endmodule
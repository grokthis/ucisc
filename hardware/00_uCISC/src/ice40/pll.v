module pll(input clki, output clko);
    parameter CLK_DIVR = 0;
    parameter CLK_DIVF = 0;

    SB_PLL40_CORE #(
        .FEEDBACK_PATH("PHASE_AND_DELAY"),
        .DELAY_ADJUSTMENT_MODE_FEEDBACK("FIXED"),
        .DELAY_ADJUSTMENT_MODE_RELATIVE("FIXED"),
        .PLLOUT_SELECT("SHIFTREG_0deg"),
        .SHIFTREG_DIV_MODE(1'b0),
        .FDA_FEEDBACK(4'h0),
        .FDA_RELATIVE(4'h0),
        .DIVR(CLK_DIVR),
        .DIVF(CLK_DIVF),
        .DIVQ(3'h0),
        .FILTER_RANGE(3'b001),
    ) uut (
        .REFERENCECLK   (clki),
        .PLLOUTGLOBAL   (clko),
        .BYPASS         (1'b0),
        .RESETB         (1'b1)
    );

endmodule

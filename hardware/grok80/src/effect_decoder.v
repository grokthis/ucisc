module effect_decoder (
    input [15:0] flags,
    input [2:0] effect,
    output store
);
    assign store =
        effect == 3'h0 ? flags[0] :
        effect == 3'h1 ? ~flags[0] :
        effect == 3'h2 ? flags[1] :
        effect == 3'h3 ? ~flags[1] & ~flags[0] :
        effect == 3'h4 ? 1 :
        effect == 3'h5 ? flags[3] :
        effect == 3'h6 ? flags[4] :
        0;

endmodule

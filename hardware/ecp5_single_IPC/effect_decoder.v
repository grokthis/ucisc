module effect_decoder(
  input [4:0] flags,
  input [2:0] effect,
  output reg store
);

wire overflow = flags[0];
wire zero = flags[2];
wire negative = flags[3];

always @* begin
  case (effect)
    // store if zero
    3'h0: store = zero;
    // store if not zero
    3'h1: store = ~zero;
    // store if positive
    3'h2: store = zero ~| negative;
    // store always
    3'h3: store = 1;
    // store if not negative (>= 0)
    3'h4: store = ~negative;
    // store if negative
    3'h5: store = negative;
    // store if not overflow
    3'h6: store = ~overflow;
    // never store
    3'h7: store = 0;
  endcase
end

endmodule

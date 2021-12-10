`define assert(signal, value) if (signal !== value) begin $display("ASSERTION FAILED in %m: signal != value"); $finish; end

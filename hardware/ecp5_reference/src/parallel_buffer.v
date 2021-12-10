module parallel_buffer(
    input clock,
    input [WIDTH-1:0] data_in,
    input write_in,
    input next_ready,
    output write_ready,
    output reg [WIDTH-1:0] data_out,
    output reg write_out
);
    parameter WIDTH = 8;
    initial write_out = 0;

    assign write_ready = ~write_out | next_ready;

    always @(posedge clock) begin
        if (write_ready & write_in) begin
            data_out <= data_in;
            write_out <= 1'b1;
        end else if (write_out & next_ready) begin
            data_out <= 0;
            write_out <= 1'b0;
        end
    end
endmodule
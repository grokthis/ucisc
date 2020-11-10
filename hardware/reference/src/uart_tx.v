module uart_tx(
    input clock,
    input [15:0] clock_divider,
    input [WIDTH-1:0] data_in,
    input write_en,
    output data_ready,
    output tx
);
    parameter WIDTH = 8;

    reg [15:0] clock_count;
    initial clock_count = 0;

    reg [WIDTH+1:0] data;
    reg [3:0] shift_remaining;
    initial shift_remaining = 0;
    assign tx = data_ready | data[0];
    assign data_ready = shift_remaining == 0;

    always @(posedge clock) begin
        if (data_ready & write_en) begin
            data <= { 1'b1, data_in, 1'b0 };
            shift_remaining = WIDTH + 2;
            clock_count <= 0;
        end else if (~data_ready) begin
            if (clock_count == clock_divider) begin
                shift_remaining = shift_remaining - 1;
                data <= { 1'b1, data[WIDTH:1] };
                clock_count <= 0;
            end else begin
                clock_count = clock_count + 1;
            end
        end
    end
endmodule
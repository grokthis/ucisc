TRELLIS?=/usr/share/trellis
PROJ = ecp5_reference
VERILOG := top.v $(wildcard src/*.v)
TESTS := $(basename $(wildcard **/*_tb.v))

PIN_DEF = ecp5_5g_eval.lpf
DEVICE = um5g-85k
FREQ_TARGET = 50

all: prog

$(PROJ).json: $(VERILOG) #firmware.hex
	yosys -p "synth_ecp5 -json $@ -top top" $(VERILOG)

$(PROJ)_out.config: $(PROJ).json $(PIN_DEF)
	nextpnr-ecp5 --json $(PROJ).json --lpf $(PIN_DEF) --textcfg $@ --$(DEVICE) --freq $(FREQ_TARGET) --package CABGA381

$(PROJ).svf: $(PROJ)_out.config
	ecppack --svf-rowsize 100000 --svf $(PROJ).svf $< $@

prog: $(PROJ).svf
	openocd -f ${TRELLIS}/misc/openocd/ecp5-evn.cfg -c "transport select jtag; init; svf $<; exit"

clean:
	rm -f $(PROJ).json $(PROJ)_out.config $(PROJ).svf **/*.v.iverilog*

verify: $(TESTS)
	@echo "Tests complete."

$(TESTS): %: %.v.iverilog

./test/%.v.iverilog: test/%.v
	@iverilog -o $@.out test/defines.v $(VERILOG) $^
	@vvp $@.out > $@.results
	@if grep "ASSERTION FAILED" $@.results; \
    then echo "Tests failed - $^"; echo "Run 'vvp $@.out' for results"; exit 1; \
    else echo "Tests passed - $^"; \
    fi

.SECONDARY:
.PHONY: all prog clean

PROJ=ucisc
TRELLIS?=/usr/share/trellis
VERILOG = processor_core.v instruction_decoder.v copy_decoder.v alu_decoder.v page_decoder.v effect_decoder.v alu.v register_block.v pc.v memory_block.v
#VERILOG = register_block.v

all: ${PROJ}.bit

${PROJ}.json: $(VERILOG)
	yosys -p "synth_ecp5 -json $@" $(VERILOG)

${PROJ}_out.config: ${PROJ}.json
	nextpnr-ecp5 --json $< --textcfg $@ --um5g-85k --package CABGA381 --lpf ecp5evn.lpf

${PROJ}.bit: ${PROJ}_out.config
	ecppack --svf ${PROJ}.svf $< $@

${PROJ}.svf : ${PROJ}.bit

ucisc: ${PROJ}.svf
	openocd -f ${TRELLIS}/misc/openocd/ecp5-evn.cfg -c "transport select jtag; init; svf $<; exit"

clean:
	rm -f *.svf *.bit *.config *.json

.PHONY: ucisc clean

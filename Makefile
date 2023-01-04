install_plugin:
	cd cmd/thrift-gen-kitex-reflect && go install ./

gen_idl:
	kitex -v -thrift-plugin=kitex-reflect:auto_setup idl.thrift

sinclude ./testdir/*.mk

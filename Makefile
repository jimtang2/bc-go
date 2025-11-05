# Makefile
.PHONY: proto clean

proto: clean
	buf generate --path proto/v1 --template buf.gen.yaml

clean:
	rm -rf pkg/pb/ client/src/gen/
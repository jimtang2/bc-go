# Makefile
.PHONY: proto clean docker_build_recreate

proto: clean
	buf generate --path proto/v1 --template buf.gen.yaml

clean:
	rm -rf pkg/pb/ client/src/gen/

docker_build_recreate:
	docker compose up -d --build --force-recreate dashboard stream

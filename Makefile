clean:
	make -C parser clean
	rm grammar/*.png

test:
	go test ./...

generate:
	go generate ./...

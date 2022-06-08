build:
	go build pacman.go

test: build
	@echo Test 1
	./pacman --enemies 4

clean:
	rm -rf pacman

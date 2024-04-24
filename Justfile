build:
	go build -o ortfols cmd/main.go

install:
	just build
	cp ortfols ~/.local/bin/

run:
	just build
	./ortfols ~/projects/portfolio/ortfodb.yaml

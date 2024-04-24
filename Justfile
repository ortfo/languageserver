vscode:
	cd vscode; bun vsce package

publish:
	cd vscode; bun vsce publish

build:
	go build -o ortfols cmd/main.go

install:
	just build
	cp ortfols ~/.local/bin/

run:
	just build
	./ortfols ~/projects/portfolio/ortfodb.yaml

now := `date --iso-8601=seconds`

vscode:
	cd vscode; bun vsce package

publish:
	cd vscode; bun vsce publish

build:
	go build -ldflags "-X github.com/ortfo/languageserver.BuiltAt={{ now }}" -o ortfols cmd/main.go

install:
	just build
	cp ortfols ~/.local/bin/

run:
	just build
	./ortfols ~/projects/portfolio/ortfodb.yaml

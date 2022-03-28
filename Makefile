CGO_ENABLED=0

all: marc2xml marc2json marc2mrk

marc2mrk: cmd/marc2mrk.go
	CGO_ENABLED=$(CGO_ENABLED) go build -o dist/$@ $<

marc2xml: cmd/marc2xml.go
	CGO_ENABLED=$(CGO_ENABLED) go build -o dist/$@ $<

marc2json: cmd/marc2json.go
	CGO_ENABLED=$(CGO_ENABLED) go build -o dist/$@ $<

clean:
	rm -f dist/*

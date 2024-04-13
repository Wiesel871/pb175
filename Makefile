SOURCES := $(shell templ generate > /dev/null; find cmd -type f -name '*.go')

.PHONY: all clean

all: $(SOURCES)
	go build -o ./tmp/main $(SOURCES)

clean:
	rm tmp/main

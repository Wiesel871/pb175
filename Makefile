.PHONY: all clean

all: $(SOURCES)
	templ generate
	go build -o ./tmp/main

clean:
	rm tmp/main

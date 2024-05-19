.PHONY: all clean

all: $(SOURCES)
	templ generate
	go build -o ./install/bazos.in

clean:
	rm tmp/bazos.in

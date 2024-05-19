.PHONY: all clean

all: $(SOURCES)
	templ generate
	go build -o ./install/bazos

clean:
	rm install/bazos

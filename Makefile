.PHONY: all clean

all: $(SOURCES)
	templ generate
	go build -o ./install/main.in

clean:
	rm tmp/install.in

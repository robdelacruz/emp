SHAREDIR=/usr/local/share/emp
BINDIR=/usr/local/bin

all: empd

dep:
	go get -u github.com/mattn/go-sqlite3

empd: empd.go
	go build -o empd empd.go

clean:
	rm -rf empd

install: empd
	mkdir -p $(SHAREDIR)
	cp empd $(BINDIR)

uninstall:
	rm -rf $(SHAREDIR)
	rm -rf $(BINDIR)/empd


GOPATH	= $(CURDIR)
BINDIR	= $(CURDIR)/bin

PROGRAMS = mosquitto-stats-provider

depend:
	env GOPATH=$(GOPATH) go get -u github.com/sirupsen/logrus
	env GOPATH=$(GOPATH) go get -u github.com/eclipse/paho.mqtt.golang
	env GOPATH=$(GOPATH) go get -u github.com/gorilla/websocket
	env GOPATH=$(GOPATH) go get -u github.com/gorilla/mux
	env GOPATH=$(GOPATH) go get -u golang.org/x/net/proxy
	env GOPATH=$(GOPATH) go get -u gopkg.in/ini.v1
	env GOPATH=$(GOPATH) go get -u github.com/nu7hatch/gouuid

build:
	env GOPATH=$(GOPATH) go install $(PROGRAMS)

destdirs:
	mkdir -p -m 0755 $(DESTDIR)/usr/bin

strip: build
	strip --strip-all $(BINDIR)/mosquitto-stats-provider

install: strip destdirs install-bin

install-bin:
	install -m 0755 $(BINDIR)/mosquitto-stats-provider $(DESTDIR)/usr/bin

clean:
	/bin/rm -f bin/mosquitto-stats-provider

distclean: clean
	rm -rf src/github.com/
	rm -rf src/git.ypbind.de/
	rm -rf src/gopkg.in/
	rm -rf src/golang.org/

uninstall:
	/bin/rm -f $(DESTDIR)/usr/bin

all: depend build strip install


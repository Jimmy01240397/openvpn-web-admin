builder := go
builddir := bin
exe := $(builddir)/openvpn-web-admin-backend
#config := $(builddir)/config.yaml
#install := $(builddir)/install.sh
#systemd := $(builddir)/piproxyclient.service

all: $(exe) #$(config) $(install) $(systemd)

#$(config): config.yaml
#		cp config.yaml $(config)

#$(install): install.sh
#		cp install.sh $(install)

#$(systemd): piproxyclient.service
#		cp piproxyclient.service $(systemd)

$(exe): main.go go.mod go.sum internal pkg
		$(builder) build -o $(exe) $<

.PHONY = clean

clean: 
		rm -r $(builddir)

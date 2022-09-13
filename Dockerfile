FROM golang:1.18

MAINTAINER Evgeniy Belousov <evgeny.belousov@orbitsoft.com>

# fix permissions
RUN chmod -R 775 /go \
	&& mkdir -p /.cache && chmod -R 775 /.cache \
	&& mkdir -p /go/pkg/mod && chmod -R 775 /go/pkg/mod

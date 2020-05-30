FROM golang:1.14-alpine
RUN apk add --update git
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd
RUN printf "machine github.com\nlogin jenkins\npassword PUTYOURGITHUBTOKENHERE" > ~/.netrc
ENV GOPROXY "https://proxy.golang.org,direct"
ENV CGO_ENABLED 0
ENV GOPRIVATE "github.com/RingierIMU/rsb-go-lib"
COPY . src/github.com/RingierIMU/rsb-service-portscan
RUN cd src/github.com/RingierIMU/rsb-service-portscan && go build .
RUN chown nobody: /go/src/github.com/RingierIMU/rsb-service-portscan/rsb-service-portscan
RUN chmod 500 /go/src/github.com/RingierIMU/rsb-service-portscan/rsb-service-portscan

FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/github.com/RingierIMU/rsb-service-portscan/rsb-service-portscan /rsb-service-portscan
COPY --from=0 /etc_passwd /etc/passwd
USER nobody
ENTRYPOINT ["/rsb-service-portscan"]

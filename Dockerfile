#FIRST PHASE
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

RUN adduser -D -g '' svruser

WORKDIR /svr
COPY . /svr


#DEPENDENCIES
RUN go get -d -v github.com/gin-gonic/gin
RUN go get -d -v github.com/gin-gonic/gin/render

#COMPILING
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o COMP319Server


#SECOND PHASE
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /svr/COMP319Server /

USER svruser

EXPOSE 8080

ENTRYPOINT ["/COMP319Server"]

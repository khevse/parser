FROM golang:1.8.1

ENV GOPATH /go
ENV GOBIN /go/bin
ENV PROJECT_DIR ${GOPATH}/src/github.com/khevse/invitro

COPY . ${PROJECT_DIR}

# External libs
RUN go get github.com/ncw/swift && \
    go get gopkg.in/mgo.v2 && \
    go get github.com/lib/pq && \
    go get github.com/NikSmith/cache

# Parser libs
RUN go get github.com/khevse/parser/db && \
    go get github.com/khevse/parser/engine && \
    go get github.com/khevse/parser/workers && \
    go get github.com/khevse/parser/html && \
    go get github.com/khevse/parser/page

WORKDIR ${PROJECT_DIR}

CMD go run main.go
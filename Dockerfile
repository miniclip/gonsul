ARG GONSUL=/go/src/github.com/miniclip/gonsul

FROM golang:1-alpine as build
ARG GONSUL

RUN apk --no-cache add build-base dep git
RUN mkdir -p $GONSUL
WORKDIR $GONSUL
COPY . .
RUN env
RUN make

FROM alpine
ARG GONSUL

COPY --from=build $GONSUL/bin/gonsul /usr/bin/gonsul

ENTRYPOINT [ "/usr/bin/gonsul" ]

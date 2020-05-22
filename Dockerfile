ARG GONSUL=/go/src/github.com/miniclip/gonsul

FROM golang:1.14.3-alpine3.11 as build
ARG GONSUL

RUN apk --no-cache add build-base dep git
RUN mkdir -p $GONSUL
WORKDIR $GONSUL
COPY . .
RUN make

FROM alpine
ARG GONSUL

COPY --from=build $GONSUL/bin/gonsul /usr/bin/gonsul
RUN adduser -D gonsul
USER gonsul

ENTRYPOINT [ "/usr/bin/gonsul" ]

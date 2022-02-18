FROM golang:1.17 AS build

  WORKDIR /go/src/slack-new-channel
  COPY . .

  RUN go get -d -v ./...
  RUN CGO_ENABLE=0 go install -ldflags '-w -s' -v ./...

FROM gcr.io/distroless/base

  COPY --from=build /go/bin/slack-new-channel /
  ENTRYPOINT ["/slack-new-channel"]
  CMD ["--daemon"]

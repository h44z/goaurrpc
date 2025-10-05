FROM golang:1.25-alpine AS build
WORKDIR /app

# Restore modules - Start
COPY go.mod ./
COPY go.sum ./
RUN go mod download
# Restore modules - End

COPY cmd/ cmd/
COPY internal/ internal/
ENV GOEXPERIMENT=jsonv2
RUN go build -ldflags="-s -w" -o /goaurrpc ./cmd/aur_rpc_service/main.go

FROM alpine:3.22
WORKDIR /
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "1000" \
    "nonroot"
USER nonroot:nonroot

COPY --from=build /goaurrpc /goaurrpc

ENTRYPOINT ["/goaurrpc"]

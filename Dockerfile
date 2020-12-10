ARG BUILDER_IMAGE=golang:1.15
ARG RUNTIME_IMAGE=gcr.io/distroless/static:nonroot

FROM $BUILDER_IMAGE as builder

	WORKDIR /workspace

	COPY go.mod   go.mod
	COPY go.sum   go.sum
	COPY pkg/     pkg/
	COPY cmd/     cmd/

	RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on \
		go build -a -v -o controller ./cmd/controller


FROM $RUNTIME_IMAGE

	WORKDIR /
	COPY --from=builder /workspace/controller .
	USER 65532:65532

	ENTRYPOINT ["/controller"]

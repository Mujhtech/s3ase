# ---------------------------------------------------------#
#                   Build image                            #
# ---------------------------------------------------------#
FROM --platform=$BUILDPLATFORM golang:1.26.5-alpine AS builder
RUN apk update && apk add --no-cache build-base musl-dev

WORKDIR /app
# RUN git config --global --add safe.directory '/app'

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download


# ---------------------------------------------------------#
#            COPY the source code as the last step         #
# ---------------------------------------------------------#
COPY . .

# ---------------------------------------------------------#
#                           Build the app                  #
# ---------------------------------------------------------#
ARG S3ASE_COMMIT
ARG S3ASE_VERSION_MAJOR
ARG S3ASE_VERSION_MINOR
ARG S3ASE_VERSION_PATCH
ARG TARGETOS TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build \
	--mount=type=cache,target=/go/pkg/mod \
    LDFLAGS="-X github.com/mujhtech/s3ase/version.GitCommit=${S3ASE_COMMIT} -X github.com/mujhtech/s3ase/version.major=${S3ASE_VERSION_MAJOR} -X github.com/mujhtech/s3ase/version.minor=${S3ASE_VERSION_MINOR} -X github.com/mujhtech/s3ase/version.patch=${S3ASE_VERSION_PATCH} -extldflags '-static'" && \
    CGO_ENABLED=1 \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="$LDFLAGS" -o ./s3ase ./cmd


FROM alpine:3.23 AS final

# ---------------------------------------------------------#
#                  setup app dir and its content           #
# ---------------------------------------------------------#
WORKDIR /app
VOLUME /data

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/s3ase /app/s3ase

EXPOSE 5555

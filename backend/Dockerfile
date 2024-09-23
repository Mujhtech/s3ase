# ---------------------------------------------------------#
#                   Build image                            #
# ---------------------------------------------------------#
FROM --platform=$BUILDPLATFORM golang:1.23-alpine3.20 AS builder

WORKDIR /app
# RUN git config --global --add safe.directory '/app'

COPY go.mod .
COPY go.sum .


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

RUN if [ "$TARGETARCH" = "arm64" ]; then \
    wget -P ~ https://musl.cc/aarch64-linux-musl-cross.tgz && \
    tar -xvf ~/aarch64-linux-musl-cross.tgz -C ~ ; \
    fi

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    if [ "$TARGETARCH" = "arm64" ]; then CC=~/aarch64-linux-musl-cross/bin/aarch64-linux-musl-gcc; fi && \
    LDFLAGS="-X github.com/mujhtech/s3ase/version.GitCommit=${S3ASE_COMMIT} -X github.com/mujhtech/s3ase/version.major=${S3ASE_VERSION_MAJOR} -X github.com/mujhtech/s3ase/version.minor=${S3ASE_VERSION_MINOR} -X github.com/mujhtech/s3ase/version.patch=${S3ASE_VERSION_PATCH} -extldflags '-static'" && \
    CGO_ENABLED=1 \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    CC=$CC go build -ldflags="$LDFLAGS" -o ./s3ase ./cmd


FROM --platform=$TARGETPLATFORM alpine/git:2.43.0 AS final

# ---------------------------------------------------------#
#                  setup app dir and its content           #
# ---------------------------------------------------------#
WORKDIR /app
VOLUME /data

COPY --from=builder /app/s3ase /app/s3ase

EXPOSE 5555
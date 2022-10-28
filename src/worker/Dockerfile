FROM alpine

ARG FCVersion=1.1.1

RUN apk --no-cache add \
    bash \
    curl \
    iproute2 \
    libc6-compat \
    make \
    go \
    sudo

WORKDIR /fc

RUN curl -Lo firecracker.tgz https://github.com/firecracker-microvm/firecracker/releases/download/v${FCVersion}/firecracker-v${FCVersion}-x86_64.tgz \
    && mkdir firecracker \
    && tar -xf firecracker.tgz -C firecracker \
    && chmod +x firecracker/release-v${FCVersion}-x86_64/firecracker-v${FCVersion}-x86_64 \
    && mv firecracker/release-v${FCVersion}-x86_64/firecracker-v${FCVersion}-x86_64 /usr/local/bin/firecracker \
    && chmod +x firecracker/release-v${FCVersion}-x86_64/jailer-v${FCVersion}-x86_64 \
    && mv firecracker/release-v${FCVersion}-x86_64/jailer-v${FCVersion}-x86_64 /usr/local/bin/jailer \
    && rm -rf /fc/firecracker*

COPY go.mod go.sum /fc/
RUN go mod download

WORKDIR /embr
COPY . .

RUN make cni

ENTRYPOINT ["/embr/worker"]

#[=======================================================================[
# Description : Docker image containing the godyl binary
#]=======================================================================]

ARG GO_VERSION=1.23.2
ARG DISTRO=bookworm

#### ---- Build ---- ####
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-${DISTRO} AS build

LABEL maintainer=arash.idelchi

USER root

RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

ARG TASK_VERSION=v3.39.2
ARG TARGETARCH
ARG TASK_ARCH=${TARGETARCH}
RUN wget -qO- https://github.com/go-task/task/releases/download/${TASK_VERSION}/task_linux_${TASK_ARCH}.tar.gz | tar -xz -C /usr/local/bin


WORKDIR /work

# Create User (Debian/Ubuntu)
ARG USER=user
RUN groupadd -r -g 1001 ${USER} && \
    useradd -r -u 1001 -g 1001 -m -c "${USER} account" -d /home/${USER} -s /bin/bash ${USER}

USER ${USER}
WORKDIR /tmp/go

ENV GOMODCACHE=/home/${USER}/.cache/.go-mod
ENV GOCACHE=/home/${USER}/.cache/.go

COPY go.mod go.sum ./
RUN --mount=type=cache,target=${GOMODCACHE},uid=1001,gid=1001 \
    --mount=type=cache,target=${GOCACHE},uid=1001,gid=1001 \
    go mod download

ARG TARGETOS TARGETARCH

COPY . .
ARG GODYL_VERSION="unofficial & built by unknown"
RUN --mount=type=cache,target=${GOMODCACHE},uid=1001,gid=1001 \
    --mount=type=cache,target=${GOCACHE},uid=1001,gid=1001 \
    # go install golang.org/x/tools/cmd/stringer && \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -ldflags="-s -w -X 'main.version=${GODYL_VERSION}'" -o bin/ ./cmd/...

# RUN go mod download

COPY --chown=${USER}:{USER} .bashrc /home/${USER}/.bashrc

ENV PATH=$PATH:/home/${USER}/.local/bin
ENV PATH=$PATH:/root/.local/bin

RUN mkdir -p /home/${USER}/.local/bin
RUN cp bin/godyl /home/${USER}/.local/bin


WORKDIR /home/${USER}

# Timezone
ENV TZ=Europe/Zurich

FROM debian:12 AS final

RUN apt-get update && apt-get install -y \
    curl \
    ca-certificates \
    git \
    && rm -rf /var/lib/apt/lists/*

# Create User (Debian/Ubuntu)
ARG USER=user
RUN groupadd -r -g 1001 ${USER} && \
    useradd -r -u 1001 -g 1001 -m -c "${USER} account" -d /home/${USER} -s /bin/bash ${USER}

# Install Go
# ARG GO_VERSION=go1.23.2
# ARG GO_ARCH=${TARGETARCH}
# ARG TARGETARCH
# ARG GO_ARCH=${TARGETARCH}
# RUN apt-get update && apt-get install -y \
#     wget \
#     && rm -rf /var/lib/apt/lists/*

# RUN mkdir -p /go
# RUN wget -qO- https://go.dev/dl/${GO_VERSION}.linux-${GO_ARCH}.tar.gz | tar -xz -C /go
# ENV GOPATH=/opt/go
# RUN mkdir ${GOPATH} && chown -R ${USER}:${USER} ${GOPATH}

USER ${USER}
WORKDIR /home/${USER}

COPY --from=build --chown=${USER}:{USER} /tmp/go/bin/godyl /home/${USER}/.local/bin/godyl

COPY --chown=${USER}:{USER} .bashrc /home/${USER}/.bashrc
# RUN echo "alias goo=/go/go/bin/go" >> /home/${USER}/.bashrc

ENV PATH=$PATH:/go/go/bin
ENV PATH=$PATH:/home/${USER}/.local/bin
ENV PATH=$PATH:/root/.local/bin

# Timezone
ENV TZ=Europe/Zurich

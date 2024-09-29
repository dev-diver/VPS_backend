# 단계 1: 빌드 단계
FROM golang:1.22.5-bookworm AS builder

RUN apt-get update && apt-get install -y \
  ca-certificates \
  curl \
  gnupg && \
  install -m 0755 -d /etc/apt/keyrings && \
  curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc && \
  chmod a+r /etc/apt/keyrings/docker.asc && \
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
  apt-get update && \
  apt-get install -y docker-ce docker-ce-cli && \
  update-ca-certificates && \
  curl -L https://github.com/regclient/regclient/releases/latest/download/regctl-linux-amd64 > /usr/local/bin/regctl && \
  chmod +x /usr/local/bin/regctl && \
  rm -rf /var/lib/apt/lists/*

RUN groupadd -f docker && usermod -aG docker root

# 애플리케이션 빌드
COPY /backend /backend

WORKDIR /backend
RUN go build -o server

EXPOSE 3000

# 애플리케이션 실행
CMD ["/backend/server"]
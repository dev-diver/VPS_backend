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
  curl -L "https://github.com/docker/compose/releases/download/v2.12.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && \
  chmod +x /usr/local/bin/docker-compose && \
  update-ca-certificates && \
  rm -rf /var/lib/apt/lists/*

RUN groupadd -f docker && usermod -aG docker root

COPY /docker-compose.yml /app/backend/docker-compose.yml
COPY /backend /app/backend

WORKDIR /app/backend
# Go 모듈 정리
RUN go mod tidy

# 애플리케이션 빌드
RUN go build -o server

# 단계 2: 실행 단계
FROM ubuntu

# 빌드 단계에서 빌드된 애플리케이션 복사
COPY --from=builder /app/backend/ /app/backend/
COPY --from=builder /usr/bin/docker /usr/bin/docker
COPY --from=builder /usr/local/bin/docker-compose /usr/local/bin/docker-compose
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app/backend

# 포트 설정
EXPOSE 3000

# 애플리케이션 실행
CMD ["/app/backend/server"]
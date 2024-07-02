# 단계 1: 빌드 단계
FROM golang:1.22.3 AS builder

# 필요한 패키지 설치
RUN apt-get update && apt-get install -y git

# Github 캐싱 방지
ADD https://api.github.com/repos/dev-diver/VPS_backend/git/refs/heads/main version.json

# 작업 디렉토리 설정
WORKDIR /app

# GitHub에서 최신 소스를 가져옴
RUN git clone https://github.com/dev-diver/VPS_backend.git .

WORKDIR /app/backend

# Go 모듈 정리
RUN go mod tidy

# 애플리케이션 빌드
RUN go build -o server

# 단계 2: 실행 단계
FROM ubuntu

# 빌드 단계에서 빌드된 애플리케이션 복사
COPY --from=builder /app/backend/server /app/backend/server

# 애플리케이션 실행
CMD ["/app/backend/server"]

# 포트 설정
EXPOSE 3000
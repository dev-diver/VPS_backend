# 단계 1: 빌드 단계
FROM golang:1.22.3 AS builder

COPY /backend/app /app/backend

WORKDIR /app/backend
# Go 모듈 정리
RUN go mod tidy

# 애플리케이션 빌드
RUN go build -o server

# 단계 2: 실행 단계
FROM ubuntu

# 빌드 단계에서 빌드된 애플리케이션 복사
COPY --from=builder /app/backend/server /app/backend/server

WORKDIR /app/backend

# 포트 설정
EXPOSE 3000

# 애플리케이션 실행
CMD ["/app/backend/server"]
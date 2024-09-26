# 단계 1: 빌드 단계
FROM golang:1.22.5-bookworm AS builder

# 애플리케이션 빌드
COPY /backend /backend
RUN go build -o server

FROM ubuntu

COPY --from=builder /backend/server /backend/server

WORKDIR /backend
EXPOSE 3000

# 애플리케이션 실행
CMD ["/backend/server"]
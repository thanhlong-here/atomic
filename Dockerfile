# Sử dụng Go 1.24.1 chính xác theo go.mod
FROM golang:1.24.1-alpine

# Cài thêm git nếu go get cần
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod & go.sum để cache layer go mod download
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ source code
COPY . .

# Build service từ thư mục cmd/api
RUN go build -o atomic ./cmd/api

# Mở port 8081
EXPOSE 8081

# Khởi động service
CMD ["./atomic"]

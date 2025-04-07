# 🧠 Atomic WebSocket Manager

Hệ thống WebSocket sử dụng Golang để giao tiếp real-time giữa client và server, hỗ trợ thao tác **MongoDB qua command**.

---

## 🚀 Tính năng

- Quản lý WebSocket kết nối đa client
- Gửi command dưới dạng JSON để thao tác:
  - `HandleCreate`
  - `HandleUpdate`
  - `HandleDelete`
  - `HandleFind`
  - `HandleFindMany`
- Auto dispatch command theo tên hàm
- Tích hợp BigCache & MongoDB
- Đọc cấu hình từ `.env`

---

## 🧱 Cấu trúc thư mục

```
internal/
│
├── db/              // Kết nối và truy vấn MongoDB
├── ws/              // Xử lý WebSocket + registry command
│   └── command/     // Các command thao tác Mongo
└── cache/           // BigCache xử lý cache (tùy chọn)
```

---

## ⚙️ Cài đặt

```bash
git clone https://github.com/tenban/atomic.git
cd atomic
go mod tidy
```

---

## 📦 Cấu hình `.env`

Tạo file `.env` ở thư mục gốc:

```env
MONGO_URI=mongodb://localhost:27017
```

---

## ▶️ Chạy server

```bash
go run main.go
```

WebSocket lắng nghe tại:

```
ws://localhost:8080/ws
```

---

## 🧪 Test bằng Postman

1. Mở **Postman > WebSocket Request**
2. Kết nối tới:
   ```
   ws://localhost:8080/ws
   ```
3. Gửi ví dụ:

```json
{
  "command": "HandleCreate",
  "payload": {
    "model": "posts",
    "data": {
      "title": "Hello WebSocket",
      "content": "Demo bài viết"
    }
  }
}
```

---

## ✅ Danh sách command hỗ trợ

| Command         | Mô tả                     |
|----------------|---------------------------|
| `HandleCreate` | Thêm mới vào collection   |
| `HandleUpdate` | Cập nhật theo filter      |
| `HandleDelete` | Xoá theo filter           |
| `HandleFind`   | Tìm 1 bản ghi             |
| `HandleFindMany` | Phân trang bản ghi      |

---

## 📌 Gợi ý mở rộng

- `SoftDelete` – đánh dấu thay vì xoá thật
- `Broadcast` – gửi dữ liệu tới nhiều client
- `Auth` – xác thực bằng token trước khi xử lý command
- `RateLimit` – giới hạn số lệnh/giây theo IP hoặc user

---

## 📜 Giấy phép

MIT License

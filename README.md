
เปิด terminal แล้ว รันคำสั่งนี้ใน root โปรเจกต์ (ที่มี go.mod) เพื่อให้ Go ดึง dependency ทั้งหมด
go mod tidy

รัน seed
go run cmd/seed_admin/main.go 
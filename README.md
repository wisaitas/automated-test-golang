# Automated Golang User API (Fiber + GORM + Postgres)

โปรเจกต์ตัวอย่าง API `Create User` ด้วย **Go + Fiber + GORM + Postgres** ครบพร้อมเทส 3 ระดับ (Unit / Integration / E2E) และ Docker Compose สำหรับฐานข้อมูล ฝั่งนี้เน้น **อธิบายละเอียด** ไล่ตั้งแต่ไฟล์ `docker-compose.yml`, `main.go` จนถึงโครงสร้างภายใน (`internal/user`) ว่าแต่ละบรรทัดทำอะไร ทำไมต้องเขียนแบบนั้น และจุดที่ควรระวัง

---

## สารบัญ

- [Stack และ Requirements](#stack-และ-requirements)
- [โครงสร้างโปรเจกต์](#โครงสร้างโปรเจกต์)
- [Docker Compose (Postgres)](#docker-compose-postgres)

  - [อธิบายบรรทัดต่อบรรทัด](#อธิบายบรรทัดต่อบรรทัด-compose)

- [รันแอปบนเครื่อง (เชื่อม Postgres ใน Docker)](#รันแอปบนเครื่อง-เชื่อม-postgres-ใน-docker)
- [ไฟล์ `main.go` อธิบายบรรทัดต่อบรรทัด](#ไฟล์-maingo-อธิบายบรรทัดต่อบรรทัด)
- [เลเยอร์ภายใน `internal/user`](#เลเยอร์ภายใน-internaluser)

  - `model.go`
  - `repository.go`
  - `service.go`
  - `validator.go`
  - `handler.go`

- [สัญญา API (Contract)](#สัญญา-api-contract)
- [การทดสอบ](#การทดสอบ)

  - Unit (Service/Handler)
  - Integration (Repository + Postgres via Testcontainers)
  - E2E (HTTP flow จริง)
  - คำสั่งรันและเคล็ดลับ

- [Troubleshooting & Notes](#troubleshooting--notes)
- [Next Steps](#next-steps)

---

## Stack และ Requirements

- **Go** 1.22+
- **Fiber** v2 (เว็บเฟรมเวิร์ก)
- **GORM** (ORM) + ไดรเวอร์ **Postgres**
- **Postgres 17** (รันด้วย Docker Compose)
- **testcontainers-go** (สำหรับ Integration/E2E ที่ใช้ Postgres จริงในคอนเทนเนอร์)

> หมายเหตุ: ถ้ารันเทสต์ที่อาศัย Testcontainers ต้องมี Docker Desktop/Engine เปิดอยู่

---

## โครงสร้างโปรเจกต์

```
automated-golang/
  docker-compose.yml
  go.mod
  main.go
  internal/user/
    model.go
    repository.go
    service.go
    validator.go
    handler.go
    # tests
    service_test.go                  # Unit (Service)
    handler_test.go                  # Unit (Handler)
    repo_integration_postgres_test.go# Integration (Postgres via Testcontainers)
    e2e_test.go                      # E2E (ยิง HTTP จริง)
```

---

## Docker Compose (Postgres)

ตัวอย่าง (ขั้นต่ำ) ที่คุณเพิ่มไว้:

```yaml
version: "3.8"

services:
  postgres:
    image: postgres:17
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

> แนะนำเพิ่มเติม: ใส่ **healthcheck** เพื่อให้ service อื่น (เช่นแอป) รอจน Postgres พร้อม
>
> ```yaml
> healthcheck:
>   test: ["CMD-SHELL", "pg_isready -U postgres"]
>   interval: 5s
>   timeout: 3s
>   retries: 10
> ```

### อธิบายบรรทัดต่อบรรทัด (Compose) {#อธิบายบรรทัดต่อบรรทัด-compose}

- **บรรทัด 1** `version: "3.8"` – ระบุสเปคของ compose file
- **บรรทัด 3** `services:` – ส่วนประกาศบริการที่เราจะรัน
- **บรรทัด 4** `postgres:` – ชื่อ service = `postgres` (จะกลายเป็น hostname ภายใน network ของ compose)
- **บรรทัด 5** `image: postgres:17` – ใช้ image official ของ Postgres เวอร์ชัน 17
- **บรรทัด 6-7** `ports: "5432:5432"` – map พอร์ต 5432 ของคอนเทนเนอร์ออกสู่เครื่อง host → แอปรันบนเครื่องเชื่อมต่อผ่าน `localhost:5432`
- **บรรทัด 8-11** `environment:` – เซ็ตค่าเริ่มต้นของ DB (user, password, db name)
- **บรรทัด 12-13** `volumes:` – ผูกโฟลเดอร์ข้อมูลในคอนเทนเนอร์กับ volume ชื่อ `postgres_data` เพื่อให้ data คงอยู่ระหว่างการ restart/upgrade
- **บรรทัด 15-16** `volumes:` – ประกาศ volume ชื่อ `postgres_data`

**คำสั่งใช้งานเร็ว ๆ**

```bash
docker compose up -d          # สตาร์ต Postgres
docker compose logs -f postgres
docker compose ps
```

---

## รันแอปบนเครื่อง (เชื่อม Postgres ใน Docker)

1. ติดตั้งไดรเวอร์ Postgres สำหรับ GORM

```bash
go get gorm.io/driver/postgres
```

2. สตาร์ตฐานข้อมูล

```bash
docker compose up -d
```

3. รันแอป (บนเครื่อง)

```bash
go run .     # หรือ go run main.go
```

> เมื่อรันบนเครื่อง **host** ให้ใช้ `host=localhost` ใน DSN (ดูใน `main.go` ด้านล่าง)

---

## ไฟล์ `main.go` อธิบายบรรทัดต่อบรรทัด

> หมายเหตุ: อ้างอิงตามไฟล์ที่คุณส่งมา (จำนวนบรรทัดอาจต่างถ้าเพิ่ม/ลบ import) — ข้างล่างคือ mapping โดยประมาณ **\[1–38]**

```go
1  package main

3  import (
4      "automated-golang/internal/user"
5      "fmt"
6      "log"

8      "github.com/gofiber/fiber/v2"
9      "gorm.io/driver/postgres"
10     "gorm.io/gorm"
11 )

13 func main() {
14     dsn := fmt.Sprintf(
15         "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=Asia/Bangkok",
16         "localhost",
17         "5432",
18         "postgres",
19         "postgres",
20         "postgres",
21     )
22     db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
23     if err != nil {
24         log.Fatal(err)
25     }
26     if err := db.AutoMigrate(&user.User{}); err != nil {
27         log.Fatal(err)
28     }

30     repo := user.NewGormRepository(db)
31     svc := user.NewService(repo)

33     app := fiber.New()
34     user.RegisterRoutes(app, svc)

36     log.Println("listening on :8080")
37     log.Fatal(app.Listen(":8080"))
38 }
```

### อธิบายทีละช่วง

- **บรรทัด 1**: กำหนดแพ็กเกจ `main` – ไฟล์ที่โปรแกรมเริ่มรัน

- **บรรทัด 3–11 (imports)**:

  - `internal/user` – โมดูลโดเมน user ของเรา (model/repo/service/handler)
  - `fmt`, `log` – ใช้จัดรูปประโยค DSN และพิมพ์ log/fatal
  - `github.com/gofiber/fiber/v2` – เว็บเฟรมเวิร์กสำหรับ HTTP server
  - `gorm.io/driver/postgres` + `gorm.io/gorm` – ORM และไดรเวอร์ Postgres

- **บรรทัด 13**: จุดเริ่ม `func main()`

- **บรรทัด 14–21 (ประกอบ DSN)**: สร้าง connection string ของ Postgres ด้วย `fmt.Sprintf`

  - `host=localhost` – เพราะแอปรันบนเครื่อง host และ Postgres ถูก map พอร์ตไว้ที่เครื่อง (จาก compose)
  - `port=5432` – พอร์ตเริ่มต้นของ Postgres
  - `user/password/dbname` – ต้องตรงกับค่าใน compose
  - `sslmode=disable` – ปิด SSL (สะดวกสำหรับ local dev)
  - `timezone=Asia/Bangkok` – เซ็ตโซนเวลา (ข้อแนะนำ: ใช้รูปแบบ `TimeZone=Asia/Bangkok` ตัวใหญ่-เล็กตรงตามไดรเวอร์; หลายเวอร์ชันรองรับทั้งคู่ แต่เพื่อความชัดเจนควรใช้ **`TimeZone`**)

- **บรรทัด 22**: เปิดการเชื่อมต่อ DB ด้วย GORM โดยส่ง DSN ผ่านไดรเวอร์ `postgres`

- **บรรทัด 23–25**: ถ้าเปิดไม่สำเร็จ → `log.Fatal` ทำให้โปรแกรมจบพร้อมพิมพ์ error

- **บรรทัด 26–28 (AutoMigrate)**: ให้ GORM สร้าง/อัปสคีมา table ของ `user.User` ให้ตรงกับ struct (field, index, unique constraint)

- **บรรทัด 30–31 (Wire dependencies)**:

  - สร้าง `repo := user.NewGormRepository(db)` – repository layer ที่คุยกับ DB จริง
  - สร้าง `svc := user.NewService(repo)` – business logic layer (validate, hash password, map error)

- **บรรทัด 33–34 (เว็บเซิร์ฟเวอร์)**:

  - `fiber.New()` – สร้างแอป HTTP
  - `user.RegisterRoutes(app, svc)` – ลงทะเบียนเส้นทาง `/users` (POST) ผูกกับ service

- **บรรทัด 36–37 (Start server)**:

  - log แจ้งว่าฟังบนพอร์ต `:8080`
  - `app.Listen(":8080")` – เริ่มฟัง HTTP; ถ้ามี error ใด ๆ ก็ `log.Fatal` ปิดโปรแกรม

> ถ้าคุณย้ายแอปรัน “ใน Compose เดียวกัน” กับ Postgres ให้เปลี่ยน DSN เป็น `host=postgres` (ชื่อ service) แทน `localhost`

---

## เลเยอร์ภายใน `internal/user`

แนวคิดแยกเลเยอร์เพื่อให้ทดสอบได้ง่าย และโค้ดชัดเจน

### `model.go`

- Struct `User` กำหนดรูปแบบตาราง: `ID`, `Email` (unique), `Name`, `PasswordHash`, `CreatedAt`
- DTO: `CreateUserRequest`, `CreateUserResponse` สำหรับรับ/ส่งผ่าน API

### `repository.go`

- อินเตอร์เฟซ `Repository` กำหนด method หลัก: `Create`, `FindByEmail`
- `NewGormRepository(db)` – ติดตั้ง repository บน GORM
- `Create(u *User)` – เรียก `db.Create(u)` แล้ว **map error** ของ Postgres

  - Duplicate email: Postgres error code **`23505`** → คืน `ErrDuplicateEmail`

- **ข้อดี**: ดึงการผูกพันกับ GORM/SQL ไว้ที่ชั้นนี้ ทำให้ Service สามารถ mock ได้ง่ายใน Unit test

### `service.go`

- อินเตอร์เฟซ `Service` มี `Create(ctx, req)`
- Flow ของ `Create`:

  1. `ValidateCreate(req)` – ตรวจ name/email/password
  2. `bcrypt.GenerateFromPassword` – แฮชรหัสผ่านก่อนบันทึก
  3. ประกอบ `User{...}` แล้วเรียก `repo.Create`
  4. คืน `CreateUserResponse` (ไม่คืน hash ออกไป)

### `validator.go`

- กฎขั้นต่ำ: name ต้องไม่ว่าง, email ต้องมี `@`, password ≥ 8 ตัวอักษร
- แยกไฟล์เพื่อให้ทดสอบแยกและ reuse ง่าย

### `handler.go`

- ฟังก์ชัน `RegisterRoutes(app, svc)` ผูก **POST `/users`**
- ขั้นตอนใน handler:

  1. `c.BodyParser(&req)` – แปลง JSON → struct
  2. `svc.Create(...)`
  3. แมป error → HTTP status:

     - `ErrBadEmail|ErrBadPassword|ErrBadName` → **400**
     - `ErrDuplicateEmail` → **409**
     - อื่น ๆ → **500**

  4. สำเร็จ → **201 Created** พร้อม JSON response

---

## สัญญา API (Contract)

- **POST** `/users`
- **Request (JSON)**

```json
{
  "email": "a@b.com",
  "name": "Alice",
  "password": "supersecret"
}
```

- **Success 201**

```json
{
  "id": 1,
  "email": "a@b.com",
  "name": "Alice",
  "createdAt": "2025-09-06T02:18:27Z"
}
```

- **Duplicate 409**

```json
{ "error": "email exists" }
```

- **Bad Request 400** (เช่น email ไม่ถูกต้อง)

```json
{ "error": "invalid email" }
```

**ทดสอบเร็ว ๆ**

```bash
curl -i -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"email":"a@b.com","name":"Alice","password":"supersecret"}'

curl -i -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"email":"a@b.com","name":"Dup","password":"supersecret"}'
```

---

## การทดสอบ

> เทสต์ถูกออกแบบให้สะท้อนเลเยอร์แต่ละส่วน จับบั๊กได้เร็วและแม่น

### 1) Unit Tests

- `service_test.go` – ทดสอบ business logic ด้วย mock repo (ไม่แตะ DB)
- `handler_test.go` – ทดสอบ HTTP handler ด้วย fake service (ไม่แตะ DB)

### 2) Integration Test (Repository + Postgres จริง)

- `repo_integration_postgres_test.go` – ใช้ **testcontainers-go** สปิน Postgres ชั่วคราว, `AutoMigrate`, แล้วทดสอบ `Create` + unique constraint (`23505`)
- ข้อดี: พฤติกรรมตรงกับ DB จริงของเรา 100%

### 3) E2E Test (HTTP flow จริง)

- `e2e_test.go` – สร้าง Fiber app + Service + Repository บน **Postgres ในคอนเทนเนอร์** → ยิง `POST /users` 2 ครั้ง (ครั้งสองต้องได้ 409)

### คำสั่งรันและทิปส์

```bash
go test ./... -v
```

- ต้องเปิด Docker ให้พร้อม (สำหรับเทสต์ที่ใช้ Testcontainers)
- อยากรันเฉพาะบางเทสต์:

```bash
go test ./internal/user -run TestE2E_ -v
```

---

## Troubleshooting & Notes

- **CGO / SQLite error**: ถ้าใครใช้ SQLite (`go-sqlite3`) กับ `CGO_ENABLED=0` จะเจอ `requires cgo` → ในโปรเจกต์นี้เราแก้โดย **ย้ายเทสต์ทั้งหมดไปใช้ Postgres ผ่าน Testcontainers** (ข้อแนะนำอันดับ 1)

  - ถ้ายังอยากใช้ SQLite แบบ cgo-free ใช้ไดรเวอร์ `github.com/glebarez/sqlite` (แต่พฤติกรรมบางอย่างต่างจาก Postgres)

- **TimeZone vs timezone**: ใน DSN ของ Postgres แนะนำใช้ `TimeZone=Asia/Bangkok` (ตัวพิมพ์ใหญ่ตรงตาม convention) เพื่อลดความสับสนระหว่างเวอร์ชันไดรเวอร์
- **Host ใน DSN**:

  - แอปรันบน **host** → `host=localhost`
  - แอปรันใน **Compose เดียวกัน** กับ Postgres → `host=postgres` (ชื่อ service)

- **Unique constraint**: ตรวจ `pgconn.PgError{Code: "23505"}` ใน repository เพื่อ map เป็น `ErrDuplicateEmail` ให้ handler ตอบ 409 ได้ถูกต้อง

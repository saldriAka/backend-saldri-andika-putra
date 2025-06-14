# 📂 Simple E-commerce API

RESTful API sederhana untuk platform e-commerce yang mendukung dua jenis pengguna: **Customer** dan **Merchant**. Proyek ini dibangun menggunakan Go (Fiber), GORM, MySQL, dan Docker Compose.

---

## 🚀 Cara Menjalankan Proyek

### Prasyarat

* [Docker](https://www.docker.com/)
* [Docker Compose](https://docs.docker.com/compose/)

### Langkah Menjalankan

1. **Clone repository ini**:

   ```bash
   git clone <repo-url>
   cd <repo-folder>
   ```

2. **Jalankan service menggunakan Docker Compose**:

   ```bash
   rename .env.example tp .env
   docker compose up -d
   ```

## 🧱 Struktur Proyek

```
.
├── api/             # Handler HTTP/API
├── domain/          # Interface & domain models
├── dto/             # DTOs (Data Transfer Objects)
├── repository/      # Akses database (GORM)
├── service/         # Logika bisnis
├── extra/           # dump.sql & Postman collection
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

---

## 📌 Fitur

### Customer

* ✅ Register & Login
* ✅ Melihat produk
* ✅ Checkout transaksi
* ✅ Melihat histori transaksi pribadi

### Merchant

* ✅ Register & Login
* ✅ Melihat produk
* ✅ Membuat, mengubah, dan menghapus produk sendiri
* ✅ Melihat pembeli yang membeli produknya

---

## 🔐 Autentikasi & Otorisasi

* Menggunakan session middleware (`fiber/session`)
* Role disimpan di session:

  * `"customer"`
  * `"merchant"`
* Middleware mengecek role pengguna untuk akses tertentu

---


## 🧰 Pengujian

* Gunakan file `extra/saldri-postman.json` untuk mengakses seluruh endpoint menggunakan Postman
* Semua role, header, dan payload sudah disiapkan

---

## 📝 Catatan

* Akses API: `http://localhost:8080`
* Database default: `marketplace` (MySQL)
* Session disimpan di memori

---

## 📂 Folder `extras/`

| File                  | Keterangan                                |
| --------------------- | ----------------------------------------- |
| `marketplace.sql`            | Struktur dan data awal MySQL              |
| `PTAmar.postman_collection.json` | Postman collection lengkap untuk API test |

---

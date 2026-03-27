# 📦 Smart Inventory System

A simple inventory management system built with:

* **Backend**: Golang (Gin + pgx)
* **Frontend**: React (Vite + Redux Toolkit)
* **Database**: PostgreSQL

---

# 🚀 FEATURES

* Inventory management
* Stock In (with new product support)
* Stock Out (with reservation system)
* Status workflow (in_progress → done / canceled)
* Reporting (Stock In & Out)
* Pagination & search

---

# 🗄️ DATABASE SETUP

Run the DDL to your postgresql query in DDLDB.sql
add some inventory to inventory_items table
```

---

## 1. Clone Project

```bash
git clone https://github.com/titis-elizer/smart-inventory.git
cd smart-inventory
```


---

## 2. Install Dependencies

```bash
go mod tidy
```

---

## 3. Setup `.env`

Create file:

```bash
.env
```

Example:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=smart_inventory
PORT=9090
```

👉 **IMPORTANT:**
Adjust sesuai konfigurasi PostgreSQL kamu.

---

## 4. Run Backend

```bash
go run ./cmd/server/main.go
```

Server will run at:

```
http://localhost:9090
```

---

# 🌐 FRONTEND SETUP (REACT + VITE)

## 1. Masuk ke folder frontend

```bash
cd frontend\smart-inventory-ui
```

---

## 2. Install dependencies

```bash
npm install
```

---

## 3. Run Frontend

```bash
npm run dev
```

App akan berjalan di:

```
http://localhost:5173
```

---


# 📊 WORKFLOW SYSTEM

## 📥 Stock In

* Create → `in_progress`
* Done → tambah `physical_stock`
* Cancel → tidak ada perubahan

👉 Bisa:

* pilih item existing
* atau create product baru langsung

---

## 📤 Stock Out

* Create → `allocated`
* In Progress → proses pengeluaran
* Done → kurangi stock
* Cancel → rollback reserved stock

---

## 📈 Report

Menampilkan:

* Stock In (done)
* Stock Out (done)

Dengan:

* Product
* Quantity
* Created date
* Done date

---

# 🧠 NOTES

* Semua operasi menggunakan transaction (`pgx.Tx`)
* Inventory update menggunakan row locking (`FOR UPDATE`)
* Status mengikuti state machine

---

# 🚀 RUN SUMMARY

## Backend

```bash
cd backend
go run main.go
```

## Frontend

```bash
cd frontend
npm install
npm run dev
```

---

# 💡 FUTURE IMPROVEMENTS

* Dashboard analytics
* Export CSV / Excel
* Multi-item transaction
* Product search & autocomplete
* Role-based access

---

# smart-inventory

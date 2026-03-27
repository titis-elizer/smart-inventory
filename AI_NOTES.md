AI Tools : ChatGPT

1 prompt yang paling rumit adalah salah satunya di Stock Out module
Saya ingin mengimplementasikan logic Stock Out (reservation system) pada inventory system menggunakan Golang + PostgreSQL, dengan requirement berikut:

* Saat user membuat Stock Out:

  * Sistem harus langsung melakukan reservasi (allocated)
  * reserved_stock bertambah
  * Tidak boleh melebihi available stock (physical - reserved)

* Saat status berubah:

  * in_progress → tidak mengubah stock
  * done → physical_stock dikurangi dan reserved_stock dikurangi
  * canceled → reserved_stock dikembalikan

* Sistem harus aman dari race condition:

  * Jika ada 2 request bersamaan untuk item yang sama
  * Tidak boleh terjadi over-reservation

* Semua operasi harus:

  * menggunakan database transaction
  * menggunakan row-level locking (SELECT ... FOR UPDATE)

* Tolong implementasikan:

  * service layer (transaction-safe)
  * repository layer (dengan locking)
  * validasi stock (available stock)
  * error handling jika stok tidak cukup

Pastikan:

* tidak ada kemungkinan reserved_stock > physical_stock
* tidak ada negative stock
* sistem tetap konsisten dalam kondisi high concurrency


dimana permasalahan logic yang digenerate AI, ada masalah yaitu Over-reservation, sehingga data stock yang sudah ada tidak terambil, namun data inputan yang terambil.
disini akhirnya saya revisi dengan menggunakan query FOR UPDATE pada select untuk prodak yang dikeluarkan stocknya, sehingga ada distribution pattern system

# Fleet Management Service

Backend untuk sistem manajemen armada Transjakarta (fleet management). Service menerima lokasi kendaraan lewat MQTT, menyimpannya ke PostgreSQL, menyediakan REST API untuk lokasi terakhir dan riwayat, serta mengirim event ke RabbitMQ ketika kendaraan masuk ke area geofence.

# Cara Menjalankan
- Jalankan perinta docker compose up -d pada root direktori repo ini
- Import collection postman
- Jalankan endpoint trigger publisher mock

# Arsitektur Codingan
- Menggunakan Hexagonal Architecture -> memisahkan antara business logic (core) dengan dependensi-nya supaya loosly-couple untuk memudahkan testing dan juga lebib maintainable + readable.
- Menggunakan Domain Driven Design -> Menerjemahkan dan memecah business entity menjadi domain-domain tersendiri. Lebih mudah untuk dikembangkan lagi jika akan ada perubahan ke depannya.
- Implementasi ini murni karena preferensi pribadi setelah sekian lama ngoding backend Golang.
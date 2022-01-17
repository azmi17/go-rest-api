Sebelum running project, jalankan command line berikut:
1. go mod init
2. go mod tidy
3. Buat database bernama go_rest_api dan import table categories.sql ke dalam database yang sudah dibuat
4. Untuk melihat HTTP Method dan Path Restful API cek pada file apispec.json (bisa pake swagger atau insomnia untuk check-nya)
5. Jalankan perintah go run main.go
6. Testing pakai HTTP client seperti POSTMAN dll..
7. Setiap request divalidasi dengan Authorization, maka pada header tambahkan key x-api-key dan value-nya adalah p0l1moRphY5m

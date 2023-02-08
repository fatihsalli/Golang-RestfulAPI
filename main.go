package main

import (
	"RestfulWithEcho/app"
	"RestfulWithEcho/configs"
	"RestfulWithEcho/repository"
	"RestfulWithEcho/service"
)

//TODO: #1: models klasörü yerine "models.go" ismi "models.go" olarak değiştirildi.
// ====> Soru klasör içinde olmayınca package ismi mainde kalıyor ve referans verirken hata alıyoruz???
// devamı ===> program ayağa kalkerken hata aldık

//TODO: #2: object id yerine uuid(v4) veya guid olarak değiştirildi.
//TODO: #3: createddate ve updateddate değerleri eklendi
//TODO: #4: app tarafında da nesne üretme işi metot ile tanımlanacak (NewBookController)
//TODO: #5: servis katmanı oluşturuldu. Dto nesneleri tanımlandı. Controller yerine handler ismi verildi.
//TODO: #6: repository içerisinde ReplaceOne yerine UpdateOne metotu kullanıldı.
//TODO: #7: ConnectDB => configs setup.go dosyası içine taşındı
//TODO: #8: .env üzerinden bağlantıyı okuyacak şekilde yeniden düzenlendi
//TODO: #9: RouteHandler fonksiyonu tanımlanarak route metotları fonksiyon içinde tanımlandı

func main() {
	// it has moved into app-book_handler.go
	//e := echo.New()

	dbClient := configs.ConnectDB().Database("booksDB").Collection("books")

	//TODO:Value yerine referans olarak verilecek aşağıdaki değerler???

	// to create new repository
	BookRepository := repository.NewBookRepository(dbClient)

	// to create new service
	BookService := service.NewBookService(BookRepository)

	// to create new app
	BookHandler := app.NewBookHandler(BookService)

	e := app.RouteHandler(BookHandler)
	e.Start(":8080")
}

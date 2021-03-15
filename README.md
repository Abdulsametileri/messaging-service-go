## TO DO LIST

- [x] Kullanıcılar sistemde hesap oluşturabilir ve login olabilir.
- [x] Kullanıcılar birbirlerinin kullanıcı adını bildiği sürece mesajlaşabilirler.
- [x] Kullanıcılar geçmişe dönük mesajlaşmalarına erişebilirler.
- [x] Bir kullanıcı mesaj almak istemediği diğer kullanıcıyı bloklayabilir.
- [x] Kullanıcıların aktivite (login, invalid login, vb.) logları tutulmalıdır.
- [x] Tüm hatalar kayıt altına alınmalı ve kritik detaylar kullanıcılara iletilmemelidir.
- [x] Dockerize and scalability
- [x] Unit Test Coverage

# API Endpoints

`POST api/v1/register` <br/>

`POST api/v1/login` <br/>

### JWT Middleware 
İstek yapabilmek için token gerekli. <br/>
Authorization: Bearer eydsad.....

`GET api/v1/mutateUser/:mutateUserId` <br/>
  
`GET api/v1/users` <br/>   

`GET api/v1/messagesWith/:userName` <br/>

`POST api/v1/sendMessage/:userName` <br/>

Compose ile ayağa kaldırmak için
```
$ docker-compose up --build
```

Docker kullanmadan development modunda ayağa kaldırmak için 
tags'ı vermek gerekir. 
```
$ go run --tags dev main.go
```

Örnek İstekler
http://localhost/api/v1/register <br/>
http://localhost/api/v1/sendMessage/abdulsamet <br/>
http://localhost/api/v1/mutateUser/3 <br/>

NOT: nginx ile 80. porttan serve edildiği için portu url'e yazmıyoruz.
Round robin algoritmasına göre oluşan yükü 3 containere dağıtıyor.
 
``` 
api:
    ....
    deploy:
      replicas: 3
    ...
```

Docker compose dosyasındaki deploy replicas kısmından scale edilebilir. 
Eğer docker compose kullanılmadan ayağa kaldırılacaksa 8080 portunu istek yaparken 
belirtmek gerekir.

http://localhost:8080/api/v1/register <br/>
http://localhost:8080/api/v1/sendMessage/abdulsamet <br/>
http://localhost:8080/api/v1/mutateUser/3 <br/>

### Godaki unit testleri kaldırmak ve coverage i görmek için

```
$ go test --tags dev --cover ./...
```


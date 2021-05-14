## TO DO LIST

- [x] Users can create account and login
- [x] The only way the user can message have to know the receiver user name.
- [x] Users can access their chat history.
- [x] Users can block each other
- [x] All errors must be stored and dont send critical error messages on client side
- [x] Dockerize and scalability
- [x] Unit Test Coverage 

# API Endpoints

`POST api/v1/register` <br/>

`POST api/v1/login` <br/>

### JWT Middleware

In order to request as you predict, you have to put token.
 
`Authorization: Bearer eydsad.....`

`GET api/v1/mutateUser/:mutateUserId` 
  
`GET api/v1/users`
   
`GET api/v1/messagesWith/:userName` 

`POST api/v1/sendMessage/:userName`

Build and up with docker-compose

```
$ docker-compose up --build
```

Build development mode
```
$ go run --tags dev main.go
```

Example Request URL's

**You must add :8080 if you are development mode.** 

http://localhost/api/v1/register 

http://localhost/api/v1/sendMessage/abdulsamet

http://localhost/api/v1/mutateUser/3 


As default application runs 3 instance. You can easly edit `docker-compose.yml` 

``` dockerfile
api:
    ....
    deploy:
      replicas: 3
    ...
```


### To see unit tests and their coverage, run

```shell script
$ go test --tags dev --cover ./...
```


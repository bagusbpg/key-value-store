# About
This repository demonstrates the implementation of simple key-value store in Golang. It uses a package provided by go-redis.

# How to
First, we need to install redis. The easiest way is by running the following script in your terminal.
```bash
docker run -d -p 6379:6379 --name redis redis
```
You may also verivy that redis has been running with success using this command
```bash
docker exec -it redis bash
```
and inside the bash shell, you may try some redis commands
```bash
redis-cli
INFO
SET key value
GET key
DEL key
exit
exit
```

As usual, initiate your Golang project.
```bash
mkdir <your-project-directory>
cd <your-project-directory>
go mod init <your-project-module-name>
```
And then install go-redis package.
```bash
go get github.com/go-redis/redis/v9
```
Go-redis v9 is to be used when your redis is of version 7.x.x, while go-redis v8 is for redis version 6.x.x. Note: you can check your installed redis' version from INFO command executed before.

In your service app, you need to create a redis client instance, specifying the address where redis is serving, password, and database. Sometime, you may also want to specify some other options as well, like what is the maximum number of idle connections should be allowed. In the following example, you can set such limit using `MaxIdleCons` option.
```go
rdb := redis.NewClient(&redis.Options{
	Addr:         "localhost:6379",
	Password:     "",
	DB:           0,
    MaxIdleConns: 5,
})
```
In this minimum setup, you are already able to implement key-value store using go-redis. In the following example, you can set some value to redis with expiration time of one minute, meaning that your stored value will be kept available to fetch for the next one minute since its creation.
```go
expired := time.Now().Add(time.Minute)
if err := rdb.Set(
    someContext,
    key,
    value,
    time.Until(expired)
).Err(); err != nil {
	// error handling
}
```
Fetching the stored value is also simple.
```go
value, err := rdb.Get(someContext, key).Result()
if err != nil {
	if strings.Contains(err.Error(), "nil") {
		// error handling attributed to key not found (probably client error)
	}
    // error handling attributed to internal server error
}
```
# References
- Visit the repository of [go-redis](https://github.com/go-redis/redis) if you want to report new or search for already-answered issue, make pull request, etc.
- Visit official website of [redis](https://redis.io/) if you want to explore their compendium of documentations, etc.

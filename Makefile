USER_DATABASE_NAME=connect_b_users
FRIEND_DATABASE_NAME=connect_b_friends

bin/bconnect: main.go
	go build -o bin/heroku-go-db-example main.go
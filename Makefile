USER_DATABASE_NAME=connect_b_users
EXEC=BConnect-backend

.PHONY: test

run: bin/bconnect test
	./bin/$(EXEC)

bin/bconnect:
	go build -o bin/$(EXEC) *.go

clean:
	rm -rf bin/BConnect-backend

test:
	go test -v ./...

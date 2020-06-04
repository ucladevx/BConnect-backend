USER_DATABASE_NAME=connect_b_users
EXEC=BConnect-backend

.PHONY: test

run: clean bin/bconnect test
	./bin/$(EXEC)

bin/bconnect: clean
	go build -o bin/$(EXEC) *.go

setup:
	./scripts/removal.sh
	./scripts/init.sh

clean:
	rm -rf bin/BConnect-backend

test:
	go test -v ./...

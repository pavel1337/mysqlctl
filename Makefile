# test
.PHONY: test
test:
	docker run --name mysqlctl-test --rm -e MYSQL_ROOT_PASSWORD=password -p 127.0.0.1:6603:3306 -d mysql:8.0 || true
	go test ./...

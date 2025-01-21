push:
	git add .
	git commit -m"update"
	git push

test:
	go test -cover -v ./...
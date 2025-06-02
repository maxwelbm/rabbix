all:
	go build -o rabbix main.go

run_mock:
	./rabbix add --name test --routeKey test --file mocks/exemplo01.json
all:
	go build -o rabbix main.go

payment_correct:
	./rabbix batch onbording transaction1 transaction2 closed payment1 payment2 --delay 600

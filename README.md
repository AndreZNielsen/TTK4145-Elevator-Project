# How to run the elevators


hvordan kjøre:
start to simulatorer med port 12345 og 12346 (./SimElevatorServer --port ______ i simulator mappen)
kjør go run -ldflags="-X root/config.Elevator_id=A" main.go
og så go run -ldflags="-X root/config.Elevator_id=B" main2.go
på samme maskin
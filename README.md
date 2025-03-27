# How to run the elevators

Start elevatorserver

Adjust Config:

  Set LocalElevatorServerPort to match elevatorservers port (e.g., "localhost:15657")
  
  Define the wanted number of floors
  
  Assign the elevator ID in Elevator_id (e.g., "A" or "B")
  
  List IDs of other elevators in possibleIDs (e.g., "A","B","C" for 3 elevators)
  
  Set IP addresses for the other elevators in Elevators_ip

Run the program:
  run with "go run main.go"

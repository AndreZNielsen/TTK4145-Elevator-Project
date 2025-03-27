# How to run the elevators

- Start elevatorserver

- Adjust Config: <br>
  Set LocalElevatorServerPort to match elevatorservers port (e.g., "localhost:15657") <br>
  Define the wanted number of floors <br>
  Assign the elevator ID in Elevator_id (e.g., "A" or "B") <br>
  List IDs of other elevators in possibleIDs (e.g., "A","B","C" for 3 elevators) <br>
  Set IP addresses for the other elevators in Elevators_ip <br>

- Run the program:
  run with "go run main.go"

- Due to permission issues, the program may not boot properly if launched from the integrated terminal in VS Code on Linux. In this case, try running it directly from the GNOME terminal instead.

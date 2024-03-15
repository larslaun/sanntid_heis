Overview of elevator project in TTK4145
=======================================

Code execution
--------------

The code is runnable by first starting the elevator server with
```
elevatorserver 
```
followed by running the program by writing the following in a separate terminal:
```
go run main.go <ElevatorID> <CommunicationPort> <ElevatorPort>
```
The IDs for the elevators have to begin at 0. By leaving the "Elevatorport" field empty the program will use the standard 15657 port.

Pushing the Stop button will trigger a panic for the program running that specific elevator.



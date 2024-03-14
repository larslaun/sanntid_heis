Overview of elevator project in TTK4145
=======================================

Code execution
--------------

The code is runnable by first starting the elevator server:
```
elevatorserver --port="ElevatorPort"
```
followed by running the program by writing the following:
```
go run main.go "ElevatorID" "CommunicationPort" "ElevatorPort"
```

Features
--------

- Pushing the Stop button will trigger a panic for the program running that specific elevator.



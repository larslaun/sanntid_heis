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


Modules
-------

### communicationHandler
Responsible for both distributing and collecting states, as well as distributing orders on the network. Also includes functionality for package loss handling. 
  
### elevatorDriver
This module includes handed-out code customized for our use. It manages elevator behavior, including state initialization, button handling, motor control, and behaviour decision-making for each elevator. 

### hallAssigner
Module responsible for finding the most efficient distribution of hall orders. This is done by adding an incoming order to each active elevator and assigning that order to the one with the minimal cost.

### networkDriver
Module that includes handed-out code used for broadcasting orders and states using UDP. It also has functionality for keeping control of peers on the network, used for recovery and redistribution of hall- and cab orders.
  
### watchdog
Includes two watchdog timers for handling both internal failures and the loss of elevators on the network. In both cases, it will redistribute the existing orders at the time of failure. 

### settings
Contains system settings for all modules, including constants for timing values and elevator parameters. 
  


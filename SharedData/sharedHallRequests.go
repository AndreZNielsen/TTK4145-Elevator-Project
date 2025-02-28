package sharedData

import (
	elevalgo "root/elevator"
)




 var sharedHallRequests [4][2]bool // om buttonEvent eller update skal sharedHallRequests oppdateres
 //sharedHallRequests er input i assigner.go og setAllLights() (skal skru på/av hallLys)
 //output fra assigner.go skal oppdatere elevator.requests 

//Vi må lage en update kanal-variabel/kanal som varsler hver gang det sendes eller mottas en TCP-melding 

//er det bedre/mulig å gjøre buttonEvent->sendTCPmessage til en update sånn at sharedHallRequests bare har ett input?
//Funksjonen som henter data til sharedHallRequests må hente data fra både når buttonEvent skjer(lokalt)
//og fra TCP-meldings-datastrukturen får en oppdatering. 


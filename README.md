# flightbooking

1.Get Flight Details
 
http://localhost:8080/GetFlightDetails
METHOD :GetFlightDetails

2.Add Flight Details with Existing one 

http://localhost:8080/PostFlightDetails
METHOD : POST
PARAM :[{"flightID":5,"Destination":"LAA","DepartFrom":"HOU","DepartsAt":"2020-11-02T15:30:00Z","Seats":[{"Number":1,"IsLocked":false},{"Number":2,"IsLocked":false},{"Number":3,"IsLocked":false},{"Number":4,"IsLocked":false}]}]

3.UpdateFlightDetails - able to update one data per request
http://localhost:8080/UpdateFlightDetails
METHOD : POST
PARAM: {"flightID":3,"Destination":"dgsskdhk","DepartFrom":"HOU","DepartsAt":"2020-11-02T15:30:00Z","Seats":[{"Number":1,"IsLocked":false},{"Number":2,"IsLocked":false},{"Number":3,"IsLocked":false},{"Number":4,"IsLocked":false}]}

4.BookFlight
http://localhost:8080/BookFlight
METHOD : POST
PARAM:
 {
	"FlightID" :  3,
	"SeatNumber" :  1
}


5.DeleteFlight
METHOD : DELETE
http://localhost:8080/DeleteFlight
{
	"ID" :2	
}
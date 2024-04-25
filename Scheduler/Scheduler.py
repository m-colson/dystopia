import requests
from flask import Flask, request, make_response, jsonify
app = Flask(__name__)

PFaaS_endpoint = "http://127.0.0.1:9080"
aidoom_endpoint = "http://127.0.0.1:8080" 
WorldState_endpoint = "http://127.0.0.1:9090"

'''
This endpoint gets called when someone requests a ride from the frontend.
Ends when a suitable car is found and sent to the requesters location
This is a post request and gets a json objet in this format:
    {
        "current_location":"your mom's house", 
        "destination": "your dad's house",
        "passengers": "4",
        "flags": "0,0"
    }

'''
@app.post("/ride")
def request_ride():
    if request.method == 'POST':
        json = request.get_json()
        capacity = json["passengers"]
        current_location = json['current_location']
        destination = json['destination']
        flags = json['flags']

        # get cars from inventory
        cars = find_cars(capacity, flags)

        # get closest car from PFaaS
        closest_car = find_closest_car(cars, current_location)

        # tell inventory car is being used
        mark_in_use(closest_car)

        # send the car to the requester
        send_to_requester(closest_car, current_location, destination)

        return f"sending car with id {closest_car}.", 200
    else:
        return "error", 400

'''
when the ride arrives at the requesters location, this endpoint is 
called from the worldsate. the user is sent a notification and the
cars new path to the requesters final desination is gotten from the
PFaaS
'''
@app.put("/ride/at/start")
def ride_arrived():
    print("ride arrived!")

'''
This enpoint gets called when a ride is done. The worldstate sends the 
id of the car whose ride is done as a parameter 'id'
'''
@app.put("/ride/at/dest")
def ride_done():
    endpoint = aidoom_endpoint + f"mark/ride/over?id={requests.args.get("id")}"
    request = requests(endpoint)
    print("ride done!")

'''
makes a request to the inventory to get all cars that match the requesters
needs. ie number of passengers and disabilities
'''
def find_cars(capacity, flags):
    endpoint = f"http://127.0.0.1:8080/find/cars?capacity={capacity}&flags:{flags}"
    print(endpoint)
    request = requests.get(endpoint)
    cars = request.json()
    return cars

'''
makes a request to the PFaaS to find the closest car to 
a destination from a list of cars.
'''
def find_closest_car(cars, destination):
    endpoint = PFaaS_endpoint + "/api/closest?to="
    endpoint += f"{destination}&to="

    for i, car in enumerate(cars):
        if i != len(cars):
            endpoint += car + ","
        else:
            endpoint += car
    request = requests(endpoint)

    return request.json()["id"]
    
'''
sends a put request to inventory to make a car as in use
/mark/in/use?id={id}
'''
def mark_in_use(id):
    pass

'''
makes a request to the PFaaS to find the best path from a 
cars current location to the location of the requester
/api/path?from={id}&to={destination}&
'''
def send_to_requester(id, rider_location, destination):
    endpoint = f"{WorldState_endpoint}/{id}/path?from={rider_location}&to={destination}"
    request = requests(endpoint)
    request = request.json()["pos"]
    return 


if __name__ == '__main__':
    app.run(debug=True)

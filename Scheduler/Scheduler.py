'''
By: Skyler DeVaughn
Course: ECE 479

Schdeuler.py is the file that orchistrates communication and actions made by the 
PFaaS (Path Finding as a Service) and aidoom (Inventory Management System) and the
World State (tracking all active cars)

When a user makes a request for a ride, a front end sends a POST request with /ride.
This sets a seriers of events into action, 1) asking aidoom for cars that fit the request
2) asking PFaaS for the car that is the closest to the requester 3) telling aidoom 
which car is going to be used (i.e. the cosest car) 4) telling the PFaaS to send the 
closest car to the reqester

Once these four steps are complete, the World State will keep track of the cars and 
where they are. Once the ride is over the worldstate will send a request with the 
/ride/at/dest endpoint. Then aidoom will be told to put the car that was used for that
trip that the car can be put in the inactive pool again.
'''

import requests
from flask import Flask, request
app = Flask(__name__)

PFaaS_endpoint = "http://pfaas:9080"
aidoom_endpoint = "http://inventory:8080" 
WorldState_endpoint = "http://frontend:9081"

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
    json = request.get_json()
    capacity = json["capacity"]
    current_location = json['from']
    destination = json['to']
    flags = json['flags'] if 'flags' in json else ""

    # get cars from inventory that match user requirements
    cars = find_cars(capacity, flags)
    if len(cars) == 0:
        return "no cars found", 406

    # get closest car from PFaaS
    closest_car = find_closest_car(cars, current_location)

    if closest_car == -1:
        return "no car could make a path", 406

    # tell inventory car is being used
    mark_in_use(closest_car, current_location, destination)

    # send the car to the requester
    send_to_requester(closest_car, current_location, destination)
    
    return f"sending car with id {closest_car}.", 200
'''
when the ride arrives at the requesters location, this endpoint is 
called from the worldsate. the user is sent a notification and the
cars new path to the requesters final desination is gotten from the
PFaaS
'''
@app.put("/ride/at/start")
def ride_arrived():
    print("ride arrived!")
    return ""

'''
This enpoint gets called when a ride is done. The worldstate sends the 
id of the car whose ride is done as a parameter 'id'
'''
@app.put("/ride/at/dest")
def ride_done():
    endpoint = aidoom_endpoint + f"/mark/ride/over?id={request.args.get('id')}"
    requests.post(endpoint)
    print("ride done!")
    return ""

'''
makes a request to the inventory to get all cars that match the requesters
needs. ie number of passengers and disabilities
'''
def find_cars(capacity, flags=""):
    endpoint = f"{aidoom_endpoint}/find/cars?capacity={capacity}"
    if flags != "":
        endpoint += f"&flags={flags}"
    print(endpoint)
    request = requests.get(endpoint)
    cars = request.json()
    return cars

'''
makes a request to the PFaaS to find the closest car to 
a destination from a list of cars.
'''
def find_closest_car(cars, destination):
    endpoint = f"{PFaaS_endpoint}/api/closest?to={destination}&options={','.join(str(car['id']) for car in cars)}"

    request = requests.get(endpoint)

    if "id" in request.json():
        return request.json()["id"]

    return -1
    
'''
sends a put request to inventory to make a car as in use
/mark/in/use?id={id}
/api/car/{car_id}/trip?from={node_id}&to={node_id}
'''
def mark_in_use(id, frm, to):
    print(requests.post(f"{aidoom_endpoint}/mark/in/use?id={id}"))
    print(requests.post(f"{WorldState_endpoint}/api/car/{id}/trip?from={frm}&to={to}"))

'''
makes a request to the PFaaS to find the best path from a 
cars current location to the location of the requester
/api/path?from={id}&to={destination}&
'''
def send_to_requester(id, rider_location, destination):
    endpoint = f"{WorldState_endpoint}/{id}/trip?from={rider_location}&to={destination}"
    request = requests.post(endpoint)
    return 

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=5000, debug=True)

import requests
from flask import Flask, request, make_response, jsonify
app = Flask(__name__)

PFaaS_endpoint = "http://127.0.0.1:9000"
aidoom_endpoint = "http://127.0.0.1:8080" 

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

        current_location = json['current_location']
        destination = json['destination']
        passengers = json['passengers']
        flags = json['flags']

        # get cars from inventory
        cars = find_cars(request.get_json())

        # get closest car from PFaaS
        closest_car = find_closest_car(cars, request.args['current_loaction'])

        # tell inventory car is being used
        mark_in_use(closest_car)

        # send the car to the requester
        response = send_to_requester(closest_car)
        response = response.json()
        time = response["totalTime"]

        return f"sending car with id {closest_car}, arriving in {time} time.", 200
    else:
        return "error", 400

'''
when the ride arrives at the requesters location, this endpoint is 
called from the worldsate. the user is sent a notification and the
cars new path to the requesters final desination is gotten from the
PFaaS
'''
@app.put("/ride/arrived")
def ride_arrived():
    pass

'''
This enpoint gets called when a ride is done. The worldstate sends the 
id of the car whose ride is done as a parameter 'id'
'''
@app.put("/ride/done")
def ride_done():
    if request.method == 'GET':
        print(request.args)

'''
makes a request to the inventory to get all cars that match the requesters
needs. ie number of passengers and disabilities
'''
def find_cars(args):
    endpoint = "http://127.0.0.1:8080/find/cars?"
    endpoint += "capacity=" + args['passengers'] + "&"
    endpoint += "wheelchairADA=" + args['vision'] + "&"
    endpoint += "visionImpairedADA=" + args['movement']
    print(endpoint)
    request = requests.get(endpoint)
    cars = request.json()


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
/api/path?from={id}&to={destination}
'''
def send_to_requester(id, destination):
    pass


if __name__ == '__main__':
    app.run(debug=True)

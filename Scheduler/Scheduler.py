import requests
from flask import Flask, request, make_response, jsonify
app = Flask(__name__)

PFaaS_endpoint = "http://127.0.0.1:9000"
aidoom_endpoint = "http://127.0.0.1:8080" 

'''
This endpoint gets called when someone requests a ride from the frontend
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


        # get cars from aidoom
        cars = find_cars(request.get_json())

        # get closest car from PFaaS
        closest_car = find_closest_car(cars, request.args['current_loaction'])

        return "sending ride", 200
    else:
        return "error", 400

'''
This enpoint gets called when a ride is done. The worldstate sends the 
id of the car whose ride is done as a parameter 'id'
'''
@app.put("/ride/done")
def ride_done():
    if request.method == 'GET':
        print(request.args)


def find_cars(args):
    endpoint = "http://127.0.0.1:8080/find/cars?"
    endpoint += "passengers=" + args['passengers'] + "&"
    endpoint += "vision=" + args['vision'] + "&"
    endpoint += "movement=" + args['movement']
    
    request = requests.get(endpoint)
    cars = request.json()

def find_closest_car(cars, ):
    endpoint = PFaaS_endpoint + "/api/closest?to={}"
    for car in cars:
        endpoint += 

if __name__ == '__main__':
    app.run(debug=True)

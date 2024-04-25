import requests
from flask import Flask, request, make_response, jsonify
app = Flask(__name__)


'''
This endpoint gets called when someone requests a ride from the frontend
This is a post request and gets a json objet in this format:
    {
        "current_location":"your mom's house", 
        "destination": "your dad's house",
        "passengers": "4",
        "vision": "false",
        "movement": "false"
    }
'''
@app.post("/request/ride")
def request_ride():
    if request.method == 'POST':
        print(request.get_json())
        cars = find_cars(request.get_json())

        return "sending ride", 200
    else:
        return "error", 400
'''
This enpoint gets called when a ride is done. The worldstate sends the 
id of the car whose ride is done as a parameter 'id'
'''
@app.get("/ride/done")
def ride_done():
    if request.method == 'GET':
        print(request.args)



def find_cars(args):
    endpoint = "http://127.0.0.1:8080/find/cars?"
    endpoint += "passengers=" + args['passengers'] + "&"
    endpoint += "vision=" + args['vision'] + "&"
    endpoint += "movement=" + args['movement']
    print(endpoint)
    
    request = requests.get(endpoint)


if __name__ == '__main__':
    app.run(debug=True)

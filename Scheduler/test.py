import requests

data = {
        "current_location":"your mom's house", 
        "destination": "your dad's house",
        "passengers": "4",
        "vision": "false",
        "movement": "false"
       }

response = requests.post("http://127.0.0.1:5000/ride", json=data, timeout=2.5)

if response.status_code == 200:
    data = response.text
    print(data)
else:
    print("Failed to retrieve data, status code:", response.status_code)

response = requests.post("http://127.0.0.1:5000/ride/done?id=1")
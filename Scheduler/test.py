import requests

data = {'name':'daddy', 'message':'FUCK'}
response = requests.post("http://127.0.0.1:5000/request/ride", data=data, timeout=2.5)

if response.status_code == 200:
    data = response.json()
    print("Received JSON data:")
    print(data)
else:
    print("Failed to retrieve data, status code:", response.status_code)
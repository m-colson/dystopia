import requests
from flask import Flask, request, jsonify
app = Flask(__name__)

@app.post("/request/ride")
def hello_world():
    if request.method == 'POST':
        return jsonify({'text':'FUCK', 'user':'ur mom'})


import requests
from flask import Flask, request
from flask_restful import reqparse, abort, Api, Resource

import resource.vm as VMres


app = Flask(__name__)
api = Api(app)

api.add_resource(VMres.VMInfo, "/vminfo")


if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=25252)
    
    

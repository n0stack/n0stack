import requests
from flask import Flask, request
from flask_restful import reqparse, abort, Api, Resource


app = Flask(__name__)
api = Api(app)

class Test(Resource):
    def get(self, name):
        print(name)

        return "ok", 200


api.add_resource(Test, '/<name>')


if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=25252)
    
    

import requests
from flask import Flask, request
from flask_restful import reqparse, abort, Api, Resource


app = Flask(__name__)
api = Api(app)


class VMInfo(Resource):
    """
    Get vm informations.
    """
    def get(self):
        agents = [
            {hostname: 'localhost', port: 5000},
        ]

        res = []
        for agent in agents:
            res.append(requests.get(
                'http://'+agent.hostname+str(agent.port)).json())
        
        return res.json()


api.add_resource(VMInfo, '/vminfo')


if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=25252)
    
    

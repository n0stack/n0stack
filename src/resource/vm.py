import requests
from flask_restful import reqparse, abort, Api, Resource


class VMInfo(Resource):
    """
    Get vm informations.
    """
    def get(self):
        agents = [
            {"hostname": '10.8.0.6', "port": 5000}
        ]

        res = {}
        for agent in agents:
            uri = agent["hostname"] + ':' + str(agent["port"]) + "/vm"
            response = requests.get("http://" + uri).json()
            res.update({agent["hostname"]: response})
            
        return res

        

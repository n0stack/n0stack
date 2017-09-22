import requests
from flask_restful import reqparse, abort, Api, Resource


class VM(Resource):
    """
    Get all vm informations.
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

        
class VMname(Resource):
    def get(self, name):
        agents = [
            {"hostname": '10.8.0.6', "port": 5000}
        ]

        res = {}
        for agent in agents:

            # send the get request
            uri = agent["hostname"] + ':' + str(agent["port"]) + "/vm/" + name
            response = requests.get("http://" + uri)

            # if exists on the agent
            if response.status_code == 200:
                res_json = response.json()
                res.update({"host": agent["hostname"],
                            "info": res_json})

                return res
                
        abort(404, message="{} does not exist".format(name))
    
        

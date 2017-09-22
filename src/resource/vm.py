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
    

    def post(self, name):
        """
        create vm
        {
            "host": "hogehoge.com (automatically decided if not specified)",
            "cpu": {
                "arch": "x86_64, ...",
                "nvcpu": "number of vcpus",
            },
            "memory": "memory size of VM",
            "disk": {
                "pool": "pool name where disk is stored",
                "size": "volume size"
            },
            "cdrom": "iso image path",
            "mac_addr": "mac address (automatically generated if not specified)",
            "vnc_password": "vnc password (no password if not specified)"
        }
        """

        ##################################
        #ここでvmが存在するかチェックする#
        ##################################

        parser = reqparse.RequestParser()
        parser.add_argument('host', type=str, location='json', required=False, default=None)
        parser.add_argument('cpu', type=dict, location='json', required=True)
        parser.add_argument('memory', type=dict, location='json', required=True)
        parser.add_argument('disk', type=dict, location='json', required=True)
        parser.add_argument('cdrom', type=str, location='json', required=True)
        parser.add_argument('mac_addr', type=str, location='json', required=False, default=None)
        parser.add_argument('vnc_password', type=str, location='json', required=False, default="")
        


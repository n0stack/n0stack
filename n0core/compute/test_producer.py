import sys
sys.path.append('../../')
import pulsar

from n0core.lib.proto import UpdateVMPowerStateRequest, VMPowerState
from n0core.lib.messenger import Messenger

def main():
    client = pulsar.Client('pulsar://localhost:6650')
    producer = client.create_producer(
        'persistent://sample/standalone/compute/114514')

    producer.send('Hello1')
    client.close()


if __name__ == '__main__':
    main()

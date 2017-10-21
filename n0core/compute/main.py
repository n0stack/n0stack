import sys
sys.path.append('../../')  # NOQA
import pulsar

from n0core.lib.n0mq import N0MQ
from n0core.lib.proto import CreateVMRequest


client = N0MQ('pulsar://localhost:6650')
consumer = client.subscribe('persistent://sample/standalone/compute/handle')


@consumer.on('CreateVMRequest')
def create_VM_request(message, auto_ack=False):
    print('create vm request')
    print(message)
    consumer.ack(message)

    
if __name__ == '__main__':
    client.listen()

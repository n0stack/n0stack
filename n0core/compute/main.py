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
    print(message.data)
    consumer.ack(message)

    
if __name__ == '__main__':
    client.listen()
=======
from mqhandler import MQHandler

mqhandler = MQHandler('pulsar://localhost:6550',
                      'persistent://sample/standalone/volumer/114514',
                      subscription_name='compute')


@mqhandler.handle('CreateVMRequest')
def create_VM_handler(inner_msg, messenger):
    print('create vm')
    print(inner_msg)


@mqhandler.handle('DeleteVMRequest')
def delete_vm_handler(inner_msg, messenger):
    print('delete vm')
    print(inner_msg)


if __name__ == '__main__':
    mqhandler.listen()

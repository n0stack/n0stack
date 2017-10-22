import sys
sys.path.append('../../')  # NOQA

from n0core.lib.n0mq import N0MQ
from n0core.lib.proto import (CreateVMRequest,
                              DeleteVMRequest,
                              UpdateVMRequest,
                              UpdateVMPowerStateRequest)  # NOQA

from n0core.compute.kvm import vm


client = N0MQ('pulsar://localhost:6650')
consumer = client.subscribe('persistent://sample/standalone/compute/handle')


@consumer.on('CreateVMRequest')
def create_VM_request(message, auto_ack=False):
    print('create vm request')
    print(message.data)
    consumer.ack(message)


@consumer.on('DeleteVMRequest')
def delete_VM_request(message, auto_ack=False):
    print('delete vm request')
    print(message.data)
    consumer.ack(message)


@consumer.on('UpdateVMRequest')
def update_VM_request(message, auto_ack=False):
    print('update vm request')
    print(message.data)
    consumer.ack(message)


@consumer.on('UpdateVMPowerStateRequest')
def update_VM_power_state_request(message, auto_ack=False):
    print('update vm power state request')
    print(message.data)
    consumer.ack(message)


if __name__ == '__main__':
    client.listen()

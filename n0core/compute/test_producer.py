import sys
sys.path.append('../../')
import pulsar

from n0core.lib.proto import UpdateVMPowerStateRequest, VMPowerState


def main():
    client = pulsar.Client('pulsar://localhost:6650')
    producer = client.create_producer(
        'persistent://sample/standalone/compute/114514')

    req = UpdateVMPowerStateRequest(id="some_vm",
                                    status=VMPowerState.Value('POWEROFF'))

    serialized = req.SerializeToString()

    producer.send(serialized)
    client.close()


if __name__ == '__main__':
    main()

import sys
sys.path.append('../../')
import pulsar

from n0core.lib.proto import UpdateVMPowerStateRequest


def main():
    client = pulsar.Client('pulsar://localhost:6650')
    consumer = client.subscribe('persistent://sample/standalone/compute/114514',
                                subscription_name='compute')
    while True:
        msg = consumer.receive()
        data = msg.data().encode('utf-8')
        deserialized = UpdateVMPowerStateRequest()
        deserialized.ParseFromString(data)
        print(deserialized)
        consumer.acknowledge(msg)
        
    client.close()


if __name__ == '__main__':
    main()

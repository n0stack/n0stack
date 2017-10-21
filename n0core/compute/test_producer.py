import sys
sys.path.append('../../')
import pulsar

from n0core.lib.n0mq import N0MQ
from n0core.lib.proto import CreateVMRequest


def main():
    client = pulsar.Client('pulsar://localhost:6650')
    producer = client.create_producer('persistent://sample/standalone/compute/handle')

    req = CreateVMRequest(id='1',
                          host='test',
                          arch='x86_64',
                          vcpus=1,
                          memory_mb=1024,
                          vnc_password='test')

    producer.send(req)
    client.close()


if __name__ == '__main__':
    main()

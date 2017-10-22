import sys
sys.path.append('../../')  # NoQA
import argparse

from n0core.lib.n0mq import N0MQ
from n0core.lib.proto import (CreateVMRequest,
                              DeleteVMRequest,
                              UpdateVMRequest,
                              UpdateVMPowerStateRequest,
                              VMPowerState)


def main(req):
    # type: (str) -> bool
    if req == 'CreateVM':
        req = CreateVMRequest(id='1',
                              host='test',
                              arch='x86_64',
                              vcpus=1,
                              memory_mb=1024,
                              vnc_password='test')
    elif req == 'DeleteVM':
        req = DeleteVMRequest(id='1')
    elif req == 'UpdateVM':
        req = UpdateVMRequest(id='1',
                              host='test',
                              vcpus=2,
                              memory_mb=2048,
                              vnc_password='n0stack')
    elif req == 'UpdateVMPowerState':
        req = UpdateVMPowerStateRequest(id='vyos1',
                                        status=VMPowerState.Value('BOOT'))
    else:
        print("wrong argument")
        return False

    client = N0MQ('pulsar://localhost:6650')
    producer = client.create_producer('persistent://sample/standalone/compute/handle')
    producer.send(req)

    client.close()

    return True


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Producer test')
    parser.add_argument('--request', '-r', default='CreateVM', type=str)
    args = parser.parse_args()

    result = main(args.request)

    if result:
        print("success")
    else:
        print("failed")

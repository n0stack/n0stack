import sys
sys.path.append('../../') # NoQA
from n0core.lib.proto import CreateVolumeRequest, DeleteVolumeRequest
from n0core.lib.messenger import Messenger

import pulsar


def main():
    client = pulsar.Client('pulsar://localhost:6650')
    producer = client.create_producer('persistent://sample/standalone/volumer/114514')

    req = CreateVolumeRequest(id='1', host='hoge', size_mb=1)
    Messenger.send_new_message(producer, req)

    req = DeleteVolumeRequest(id='1', host='hoge')
    Messenger.send_new_message(producer, req)

    client.close()


if __name__ == '__main__':
    main()

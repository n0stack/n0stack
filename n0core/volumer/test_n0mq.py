try:
    from n0core.lib.n0mq import N0MQ
except:
    import sys
    sys.path.append('../../')
    from n0core.lib.n0mq import N0MQ

from n0core.lib.proto import CreateVolumeRequest, DeleteVolumeRequest


def main():
    client = N0MQ('pulsar://localhost:6650')
    producer = client.create_producer('persistent://sample/standalone/volumer/test')

    req = CreateVolumeRequest(id='1', host='test', size_mb=1)
    producer.send(req)

    req = DeleteVolumeRequest(id='1', host='test')
    producer.send(req)


if __name__ == '__main__':
    main()

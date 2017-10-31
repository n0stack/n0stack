try:
    from n0core.lib.n0mq import N0MQ
except:
    import sys
    sys.path.append('../../')
    from n0core.lib.n0mq import N0MQ

from uuid import uuid4 as uuid

from n0core.lib.proto import CreateVolumeRequest, DeleteVolumeRequest


client = N0MQ('pulsar://localhost:6650')
producer = client.create_producer('persistent://sample/standalone/volumer/test')


def test_create(id, host, size_mb):
    req = CreateVolumeRequest(id=id, host=host, size_mb=size_mb)
    producer.send(req)


def test_delete(id, host):
    req = DeleteVolumeRequest(id=id, host=host)
    producer.send(req)


def main():
    # test_create(str(uuid()), 'test', 1024)
    test_delete(str(uuid()), 'test')


if __name__ == '__main__':
    main()

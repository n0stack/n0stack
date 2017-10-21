try:
    from n0core.lib.n0mq import N0MQ
except:
    import sys
    sys.path.append('../../')
    from n0core.lib.n0mq import N0MQ


client = N0MQ('pulsar://localhost:6650')
consumer = client.subscribe('persistent://sample/standalone/volumer/test')


@consumer.on('CreateVolumeRequest')
def on_create_volume_req(message):
    print('CreateVolumeRequest')
    req = message.data
    print(req.id)


@consumer.on('DeleteVolumeRequest')
def on_delete_volume_req(message):
    print('DeleteVolumeRequest')
    req = message.data
    print(req.id)


if __name__ == '__main__':
    client.listen()

try:
    from n0core.lib.n0mq import N0MQ
except:
    import sys
    sys.path.append('../../') # NoQA
    from n0core.lib.n0mq import N0MQ


client = N0MQ('pulsar://localhost:6650/')
consumer = client.subscribe('persistent://sample/standalone/volumer/114514')

@consumer.on('CreateVolumeRequest')
def create_volume_request(message):
    print('CreateVolumeRequest')
    print(message.data)
    consumer.ack(message)

@consumer.on('DeleteVolumeRequest')
def delete_volume_request(message):
    print('DeleteVolumeRequest')
    print(message.data)
    consumer.ack(message)


if __name__ == '__main__':
    client.listen()

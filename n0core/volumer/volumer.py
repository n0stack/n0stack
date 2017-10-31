try:
    from n0core.lib.n0mq import N0MQ
except:
    import sys
    sys.path.append('../../')
    from n0core.lib.n0mq import N0MQ

from n0library.logger import Logger
from operation import Volume


client = N0MQ('pulsar://localhost:6650')
consumer = client.subscribe('persistent://sample/standalone/volumer/test')
volume = Volume()
logger = Logger()


@consumer.on('CreateVolumeRequest')
def on_create_volume_req(message):
    logger.info('CreateVolumeRequest')

    data = message.data
    host = data.host
    size_mb = data.size_mb

    if not volume.create(host, size_mb):
        logger.error('unable to create volume: {}'.format(str(data)))
        return False


@consumer.on('DeleteVolumeRequest')
def on_delete_volume_req(message):
    logger.info('DeleteVolumeRequest')

    data = message.data
    host = data.host

    if not volume.delete(host):
        logger.error('unable to delete volume: {}'.format(str(data)))
        return False


if __name__ == '__main__':
    client.listen()

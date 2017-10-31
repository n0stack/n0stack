import pulsar  # NOQA
try:
    from n0core.lib import proto
except:  # NOQA
    import sys
    sys.path.append('../../')
    from n0core.lib import proto  # NOQA
from initialize import consumer, logger, send, volumer_producer  # NOQA


@consumer.on('CreateVolumeRequest')
def create_volume_req(message, auto_ack=False):
    # type: (pulsar.Message, bool) -> bool
    logger.info('CreateVolumeRequest')
    data = message.data
    req = proto.CreateVolumeRequest(id=data.id,
                                    host=data.host,
                                    size_mb=data.size_mb)
    send(volumer_producer + data.host, req)
    consumer.ack(message)


@consumer.on('DeleteVolumeRequest')
def delete_volume_req(message, auto_ack=False):
    # type: (pulsar.Message, bool) -> bool
    logger.info('DeleteVolumeRequest')
    data = message.data
    req = proto.DetachVolumeRequest(id=data.id,
                                    host=data.host)
    send(volumer_producer + data.host, req)
    consumer.ack(message)

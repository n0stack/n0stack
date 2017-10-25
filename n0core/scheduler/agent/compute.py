from initialize import consumer, logger, send
try:
    from n0core.lib import proto
except:  # NOQA
    import sys
    sys.path.append('../../')
    from n0core.lib import proto


compute_producer = 'persistent://sample/standalone/compute/'


@consumer.on('CreateVMRequest')  # type: ignore
def create_VM_request(message, auto_ack=False):
    # type: (pulsar.Message, bool) -> bool
    logger.info('Received CreateVMRequest')
    hostid = message.data.host
    print(hostid)
    req = proto.CreateVMRequest(id='1', host='test')
    send(compute_producer + hostid, req)


@consumer.on('DeleteVMRequest')  # type: ignore
def delete_VM_request(message, auto_ack=False):
    # type: (pulsar.Message, bool) -> bool
    logger.info('Received DeleteVMRequest')
    hostid = message.data.host
    print(hostid)
    req = proto.DeleteVMRequest(id='1', host='test')
    send(compute_producer + hostid, req)

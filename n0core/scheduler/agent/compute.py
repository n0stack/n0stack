try:
    from n0core.lib import proto
except:  # NOQA
    import sys
    sys.path.append('../../')
    from n0core.lib import proto
from initialize import consumer, logger, send, compute_producer


@consumer.on('CreateVMRequest')  # type: ignore
def create_VM_request(message, auto_ack=False):
    # type: (pulsar.Message, bool) -> bool
    logger.info('Received CreateVMRequest')
    data = message.data
    print(data.host)
    req = proto.CreateVMRequest(id=data.id,
                                host=data.host,
                                arch=data.arch,
                                vcpus=data.vcpus,
                                memory_mb=data.memory_mb,
                                vnc_password=data.vnc_password)
    send(compute_producer + data.host, req)
    consumer.ack(message)


@consumer.on('DeleteVMRequest')  # type: ignore
def delete_VM_request(message, auto_ack=False):
    # type: (pulsar.Message, bool) -> bool
    logger.info('Received DeleteVMRequest')
    data = message.data
    print(data.host)
    req = proto.DeleteVMRequest(id=data.id,
                                host=data.host)
    send(compute_producer + data.host, req)
    consumer.ack(message)


@consumer.on('UpdateVMRequest')  # type: ignore
def update_VM_request(message, auto_ack=False):
    # type: (pulsar.Message, bool) -> bool
    logger.info('Received UpdateVMRequest')
    data = message.data
    req = proto.UpdateVMRequest(id=data.id,
                                host=data.host,
                                arch=data.arch,
                                vcpus=data.vcpus,
                                memory_mb=data.memory_mb,
                                vnc_password=data.vnc_password)
    send(compute_producer + data.host, req)
    consumer.ack(message)


@consumer.on('UpdateVMPowerStateRequest')  # type: ignore
def update_VM_power_state_request(message, auto_ack=False):
    # type: (pulsar.Message, bool) -> bool
    logger.info('Received UpdateVMPowerStateRequest')
    data = message.data
    req = proto.UpdateVMPowerStateRequest(id=data.id,
                                          host=data.host,
                                          VMPowerState=data.VMPowerState)
    send(compute_producer + data.host, req)
    consumer.ack(message)

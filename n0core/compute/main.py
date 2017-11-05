import pulsar  # NOQA

try:
    from n0core.lib.n0mq import N0MQ
except:  # NOQA
    import sys
    sys.path.append('../../')
    from n0core.lib.n0mq import N0MQ

from n0library.logger import Logger
from n0core.compute.kvm import VM


client = N0MQ('pulsar://localhost:6650')
consumer = client.subscribe('persistent://sample/standalone/compute/handle')
vm = VM()
logger = Logger()


@consumer.on('CreateVMRequest')  # type: ignore
def create_VM_request(message, auto_ack=False):
    # type: (pulsar.Message, bool) -> bool

    logger.info('Received CreateVMRequest')

    data = message.data
    vm_name = data.id
    host = data.host  # NOQA
    vcpus = data.vcpus
    memory_mb = data.memory_mb
    vnc_password = data.vnc_password
    # TODO:reveive disk_path, cdrom, device, mac_addr, nic_type
    disk_path = 'hoge'
    cdrom = 'hoge'
    device = 'hoge'
    mac_addr = 'hoge'
    nic_type = 'hoge'

    if not vm.create(vm_name,
                     vcpus,
                     memory_mb,
                     disk_path,
                     cdrom,
                     device,
                     mac_addr,
                     vnc_password,
                     nic_type):
        logger.error('failed to create vm: {}'.format(vm_name))
        return False

    consumer.ack(message)
    return True


@consumer.on('DeleteVMRequest')  # type: ignore
def delete_VM_request(message, auto_ack=False):
    # type: (pulsar.Message, bool) -> bool

    logger.info('Received DeleteVMRequest')

    data = message.data
    vm_name = data.id

    if not vm.delete(vm_name):
        logger.error('Failed to delete vm: {}'.format(vm_name))
        return False

    consumer.ack(message)
    return True


@consumer.on('UpdateVMRequest')  # type: ignore
def update_VM_request(message, auto_ack=False):
    # type: (pulsar.Message, bool) -> bool

    logger.info('Received UpdateVMRequest')

    data = message.data
    vm_name = data.id
    new_memory = data.memory_mb

    if not vm.update(vm_name, new_memory):
        logger.error('Failed to update vm: {}'.format(vm_name))
        return False

    consumer.ack(message)
    return True


@consumer.on('UpdateVMPowerStateRequest')  # type: ignore
def update_VM_power_state_request(message, auto_ack=False):
    # type: (pulsar.Message, bool) -> bool

    logger.info('Received UpdateVMPowerStateRequest')

    data = message.data
    vm_name = data.id
    status = data.status

    if status == 0:
        if not vm.start(vm_name):
            logger.error('[Failed] start vm: {}'.format(vm_name))
            return False
        logger.info('[Success] start vm: {}'.format(vm_name))

    elif status == 1:
        # TODO: REBOOT
        pass

    elif status == 2:
        if not vm.stop(vm_name):
            logger.error('[Failed] stop vm: {}'.format(vm_name))
            return False
        logger.info('[Success] stop vm: {}'.format(vm_name))

    elif status == 3:
        if not vm.force_stop(vm_name):
            logger.error('[Failed] force_stop vm: {}'.format(vm_name))
            return False
        logger.info('[Success] force_stop vm: {}'.format(vm_name))

    consumer.ack(message)
    return True


if __name__ == '__main__':
    client.listen()

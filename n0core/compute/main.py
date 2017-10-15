from mqhandler import MQHandler

mqhandler = MQHandler('pulsar://localhost:6550',
                      'persistent://sample/standalone/volumer/114514',
                      subscription_name='compute')


@mqhandler.handle('CreateVMRequest')
def create_VM_handler(inner_msg, messenger):
    print('create vm')
    print(inner_msg)


@mqhandler.handle('DeleteVMRequest')
def delete_vm_handler(inner_msg, messenger):
    print('delete vm')
    print(inner_msg)


if __name__ == '__main__':
    mqhandler.listen()

from mqhandler import MQHandler


mqhandler = MQHandler('pulsar://localhost:6650',
                      'persistent://sample/standalone/volumer/114514',
                      subscription_name='volumer')


@mqhandler.handle('CreateVolumeRequest')
def create_volume_handler(inner_msg, messenger):
    print('createvolume')
    print(inner_msg)


@mqhandler.handle('DeleteVolumeRequest')
def delete_volume_handler(inner_msg, messenger):
    print('deletevolume')
    print(inner_msg)


if __name__ == '__main__':
    mqhandler.listen()

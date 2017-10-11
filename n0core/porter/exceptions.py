class ReceivedUnsupportedMessage(Exception):
    """Raise when received not supported message.

    You sent wrong message or forgot the message implementation.

    Args:
        message_type: Set message type defined on protobuf.
    """
    
    def __init__(self, message_type):
        # type: (str) -> None
        self.message_type = message_type

    def __str__(self):
        # type: () -> str
        return "on {}".format(self.message_type)

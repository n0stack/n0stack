from typing import Tuple  # NOQA

from n0core.model import Model  # NOQA


class Target(object):
    """Target is application service to apply resources with some framework like KVM and iproute2.

    A target manage only one type `*/*/*` of resource like `resource/network/flat`.
    Directory structure and class name is ruled by resource type.
    For example, `resource/network/flat` define `class Flat` which is placed on `n0core.resource.network.flat`.

    Do not kill resource when target is killed.
    """

    def __init__(self):
        # type: () -> None
        pass

    def apply(self, model):
        # type: (Model) -> Tuple[Model, bool, str]
        """Apply resource with some framework.

        Args:
            model: Model which you want to apply.

        Return:
            Tuple of processed model, is_succeeded and description.
        """
        pass

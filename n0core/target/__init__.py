from typing import Tuple  # NOQA

from n0core.model import Model  # NOQA


class Target(object):
    """Application service to apply resources with some framework like KVM and iproute2.
    """

    def __init__(self):
        # type: () -> None
        pass

    def apply(self, model):
        # type: (Model) -> Tuple[Model, bool, str]
        """Apply resource with some framework.

        Args:
            model: Model is Model which you want to apply.

        Return:
            Tuple of processed model, is_succeeded and description.
        """
        pass

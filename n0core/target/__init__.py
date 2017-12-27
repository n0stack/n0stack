from typing import Tuple  # NOQA

from n0core.model import Model  # NOQA


class Target(object):
    def __init__(self):
        # type: () -> None
        pass

    def apply(self, model):
        # type: (Model) -> Tuple[Model, bool, str]
        """
        Args:
            model: model is Model which you want to apply.

        Return:
            Tuple of processed model, succeeded and description.
        """
        pass

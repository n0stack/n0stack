import unittest

from n0test.ga import Generation


class PreviousGeneration(object):
    def __init__(self):
        self._chosen = []
        self._previous = []
        self._previous_score = 0
        self._candidate = []
        self._score = []


class TestGeneration(unittest.TestCase):
    def test___init__(self):
        prev = Generation(None, None, previous=PreviousGeneration())
        prev._candidate = [{"foo": "bar"}]
        prev._score = [-1]

        gen = Generation(None, None, previous=prev)
        self.assertEqual(len(gen._previous), 0)

    def test__random_string(self):
        self.assertEqual(len(Generation._random_string(10)), 10)

    def test__random_value(self):
        self.assertEqual(type(Generation._random_value(1)), int)
        self.assertEqual(type(Generation._random_value(1.)), float)
        self.assertEqual(type(Generation._random_value("hoge")), str)
        self.assertEqual(type(Generation._random_value([])), list)
        self.assertEqual(type(Generation._random_value({})), dict)

    def test_cross(self):
        gen = Generation(None, None, [
            {
                "foo": "bar",
                "hoge": "hoge",
            },
            {
                "foo": "baa",
                "hoge": "hage",
            },
        ], previous=PreviousGeneration())
        self.assertNotEqual(gen.cross(), gen.cross(), msg="This test fails probability, try a few times")

    def test_mutation(self):
        seed = [{"foo": "bar"}]
        gen = Generation(None, None, seed, previous=PreviousGeneration())
        self.assertEqual(gen.mutation(), seed)

import unittest

from n0test.ga import Generation


class TestGeneration(unittest.TestCase):
    def test___init__(self):
        prev = Generation(None, None, _scoring_seed=False)
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
        ], _scoring_seed=False)
        self.assertNotEqual(gen.cross(), gen.cross(), msg="This test fails probability, try a few times")

    def test_mutation(self):
        seed = [{"foo": "bar"}]
        gen = Generation(None, None, seed, _scoring_seed=False)
        self.assertNotEqual(gen.mutation(), seed)

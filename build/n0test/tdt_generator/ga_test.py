import unittest

from n0test.tdt_generator.ga import Generation


class TestGeneration(unittest.TestCase):
    def test___init__(self):
        prev = Generation(seed={"foo": "bar"} , _scoring_seed=False)
        prev._score = [-1]

        gen = Generation(previous_generation=prev)
        self.assertEqual(len(gen._previous), 0)

    def test__next_value(self):
        seed = [{
            "foo": "bar",
        }]
        gen = Generation(seed=seed, _scoring_seed=False)

        self.assertEqual(gen._next_value(1, 1), 2)
        self.assertEqual(gen._next_value(1, -1), 0)
        self.assertEqual(gen._next_value("aaa", 1), "aab")
        self.assertEqual(gen._next_value("aab", -1), "aaa")
        self.assertEqual(gen._next_value({"key": "aaa"}, 1), {"key": "aab"})
        self.assertEqual(gen._next_value({"key": "aab"}, -1), {"key": "aaa"})
        self.assertEqual(gen._next_value(["aaa"], 1), ["aaa", "aaa"])
        self.assertEqual(gen._next_value(["aaa"], 2), ["aab"])
        self.assertEqual(gen._next_value(["aaa", "aaa"], -1), ["aaa"])

    def test_cross(self):
        seed = {
            "foo": "bar",
            "hoge": "hoge",
        }
        gen = Generation(seed=seed, _scoring_seed=False)

        gen._persistence = [
            {
                "foo": "bar",
                "hoge": "hoge",
            },
            {
                "foo": "baa",
                "hoge": "hage",
            },
        ]
        self.assertNotEqual(gen.cross()[0], gen._persistence[0], msg="This test fails probability, try a few times")

    def test_mutation(self):
        seed = {
            "foo": "bar",
        }
        result = {
            "foo": "bas",
        }
        gen = Generation(seed=seed, _scoring_seed=False)
        self.assertEqual(gen.mutate(1), result)

    def test__str_to_dec(self):
        seed = [{
            "foo": "bar",
        }]
        gen = Generation(seed=seed, _scoring_seed=False)

        self.assertEqual(gen._str_to_dec(""), 0)
        self.assertEqual(gen._str_to_dec("0"), 1)
        self.assertEqual(gen._str_to_dec("00"), 101)
        self.assertEqual(gen._str_to_dec("\f"), 100)


    def test__dec_to_str(self):
        seed = [{
            "foo": "bar",
        }]
        gen = Generation(seed=seed, _scoring_seed=False)

        self.assertEqual(gen._dec_to_str(0), "")
        self.assertEqual(gen._dec_to_str(1), "0")
        self.assertEqual(gen._dec_to_str(101), "00")
        self.assertEqual(gen._dec_to_str(100), "\f")

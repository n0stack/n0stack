import string
import random
import numpy.random as nprd
import numpy as np


class Generation:
    str_char = string.printable
    window_size = 32

    def __init__(self, seed=None, previous=[], previous_generation=None, list_operate_distance=100, _scoring_seed=True):
        self._seed = seed
        self._persistence = []
        self._selected = []
        self._previous = []
        self._previous_score = 0
        self._result = []

        if previous_generation:
            self._seed = previous_generation._seed
            self._persistence = previous_generation._persistence
            self._previous = previous_generation._selected
            self._previous_score = previous_generation.max_score
            self._list_operate_distance = previous_generation._list_operate_distance

        elif seed:
            self._previous.append(seed)
            self._list_operate_distance = list_operate_distance

            if previous:
                self._persistence = previous
            else:
                self._persistence = [seed]

            if _scoring_seed:
                out, err = self.run_get_proc(self._persistence).communicate()
                self._previous_score = self.run_get_score(out.decode('utf-8'), err.decode('utf-8'), self._persistence)

        else:
            raise Exception("set seed or previous_generation")

        self._candidate = []

    @property
    def persistence(self):
        return self._persistence

    @property
    def len(self):
        return len(self._candidate)

    @property
    def max_score(self):
        return max(map(lambda x: x["score"], self._result))

    @property
    def average_distance(self):
        return np.average(list(map(lambda x: x["distance"], self._result)))

    def run_get_proc(self, request_list):
        """
        Arguments:
            request_list[list<dict<str, any>>]:

        Returns:
            process[subprocess.Popen]
        """
        raise NotImplementedError

    def run_get_score(self, out, err, chosen):
        """
        Arguments:
            out[str]:
            err[str]:
            chosen[map<str, any>]

        Returns:
            score[float or int]
        """
        raise NotImplementedError

    def cross(self):
        """
        一様乱数でふたつ選択したものを一様交叉法(p=1/2)で交叉する
        """
        c1 = self._previous[nprd.randint(0, len(self._previous))].copy()
        c2 = self._previous[nprd.randint(0, len(self._previous))].copy()

        for k in c1.keys():
            if nprd.rand() < 0.5:
                tmp = c1[k]
                c1[k] = c2[k]
                c2[k] = tmp

        self._candidate.append(c1)
        self._candidate.append(c2)

        return c1, c2

    def mutate(self, distance):
        """
        一様乱数で選択したケースを、distanceだけすすめる
        """
        if distance == 0:
            raise Exception("Do not set distance as 0")

        mutating = self._previous[nprd.randint(0, len(self._previous))].copy()

        mutating = self._next_value(mutating, distance)
        self._candidate.append(mutating)

        return mutating

    def select(self, select_num=256):
        filtered = filter(lambda x: x["score"] != -1, self._result)
        selected = sorted(filtered, key=lambda x: x["distance"], reverse=True)[:select_num]
        self._selected = list(map(lambda x: x["case"], selected)) + self._persistence

        return self._selected

    def persistent(self):
        if self.max_score <= self._previous_score:
            return None

        for i, r in enumerate(self._result):
            if r["score"] == self.max_score:
                self._persistence.append(self._candidate[i])
                return self._candidate

    def run(self):
        self._result = []
        for window in [self._candidate[i:i+self.window_size] for i in range(0, len(self._candidate), self.window_size)]:
            procs = []
            for c in window:
                request_list = [c]
                request_list.extend(self._persistence)

                procs.append((c, self.run_get_proc(request_list)))

            for req, p in procs:
                out, err = p.communicate()
                score = self.run_get_score(out.decode('utf-8'), err.decode('utf-8'), req)
                self._result.append({
                    "distance": self._get_distance(req, self._seed),
                    "case": req,
                    "score": score
                })

    def __next_string(self, previous, distance):
        dec = self._str_to_dec(previous) + distance
        return self._dec_to_str(dec)

    def _next_value(self, previous, distance):
        """
        Supported types:
            - int -- return previous + random value (-1 or 1)
            - str -- return random string: length=exponential random value ([0, int32_max] avg=15)
            - list -- return previous value
                - append
                - delete
                - modify
            - dict -- return previous value
        """
        if type(previous) == int:
            return previous + distance

        if type(previous) == str:
            ret = previous
            return self.__next_string(ret, distance)

        if type(previous) == list:
            ret = previous.copy()

            if len(previous) == 0:
                return previous

            if 0 < distance and distance <= self._list_operate_distance:
                """
                append
                """
                if type(previous[-1]) == dict:
                    ret.append(previous[-1].copy())
                else:
                    ret.append(previous[-1])

            elif 0 > distance and distance >= -1 * self._list_operate_distance:
                """
                delete
                """
                ret.pop()

            else:
                distance -= self._list_operate_distance
                target_index = nprd.randint(0, len(previous))
                ret[target_index] = self._next_value(ret[target_index], distance)

            return ret

        if type(previous) == dict:
            if len(previous) == 0:
                return previous

            target_key = list(previous.keys())[nprd.randint(0, len(previous))]
            ret = previous.copy()

            ret[target_key] = self._next_value(ret[target_key], distance)
            return ret

        return None

    def _str_to_dec(self, s):
        dec = 0

        for c in s:
            dec *= len(self.str_char)
            dec += self.str_char.index(c) + 1

        return dec

    def _dec_to_str(self, d):
        s = []

        if d <= 0:
            return ""

        while True:
            s.insert(0, self.str_char[(d - 1) % len(self.str_char)])
            d = (d - 1) // len(self.str_char)

            if d == 0:
                break

        return "".join(s)

    def _get_distance(self, a, b):
        if type(a) == int:
            return abs(a-b)

        if type(a) == str:
            return abs(self._str_to_dec(a) - self._str_to_dec(b))

        if type(a) == list:
            len_distance = abs(len(a) - len(b))

            value_distance = 0
            for i in range(min(len(a), len(b))):
                value_distance += self._get_distance(a[i], b[i])

            return len_distance + value_distance

        if type(a) == dict:
            distance = 0
            for k in a.keys():
                distance += self._get_distance(a[k], b[k])

            return distance

        return 0

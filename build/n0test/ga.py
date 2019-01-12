import string
import random
import numpy.random as nprd


class Generation:
    def __init__(self, runner, callback, chosen=[], previous=None, _scoring_seed=True):
        self._runner = runner
        self._callback = callback

        self._chosen = []
        self._previous = []
        self._previous_score = 0

        if chosen:
            self._chosen = chosen

        if previous:
            self._chosen.extend(previous._chosen)
            self._previous = [v for i, v in enumerate(previous._candidate) if previous._score[i] != -1]
            self._previous_score = previous.max_score
        elif _scoring_seed:
            out, err = self._runner(self._chosen).communicate()
            self._previous_score = self._callback(out.decode('utf-8'), err.decode('utf-8'), self._chosen)

        self._candidate = []
        self._score = []

    def cross(self):
        """
        一様乱数でふたつ選択したものを一様交叉法(p=1/2)で交叉する
        """

        candidate = []
        candidate.extend(self._chosen)
        candidate.extend(self._previous)

        c1 = candidate[int(nprd.rand() * len(candidate))].copy()
        c2 = candidate[int(nprd.rand() * len(candidate))].copy()

        for k in c1.keys():
            if nprd.rand() < 0.5:
                tmp = c1[k]
                c1[k] = c2[k]
                c2[k] = tmp

        self._candidate.append(c1)
        self._candidate.append(c2)
        return c1, c2

    @staticmethod
    def _random_string(length):
        return ''.join([random.choice(string.printable) for i in range(length)])

    @staticmethod
    def _random_value(previous):
        """
        安全側にたおして、試行回数でカバーする

        Supported types:
            - int -- return previous + random value (-1 or 1)
            - float -- return previous + exponential random value (avg=0.1)
            - str -- return random string: length=exponential random value ([0, int32_max] avg=15)
            - list -- return previous value
                - append
                - delete
                - modify
            - dict -- return previous value
        """

        if type(previous) == int:
            return previous + (1 if nprd.rand() < 0.5 else -1)

        if type(previous) == float:
            return previous + float(nprd.exponential(0.1))

        if type(previous) == str:
            return Generation._random_string(int(nprd.exponential(1.) * 15))

        if type(previous) == list:
            t = nprd.rand()
            ret = previous.copy()

            if len(previous) == 0:
                return previous

            if t < 0.1:
                """
                delete
                """
                ret.pop()

            elif t < 0.5:
                """
                append
                """
                if type(previous[-1]) == dict:
                    ret.append(previous[-1].copy())
                else:
                    ret.append(previous[-1])

            else:
                """
                modify
                """
                target_index = nprd.randint(0, len(previous))
                ret[target_index] = Generation._random_value(ret[target_index])

            return ret

        if type(previous) == dict:
            if len(previous) == 0:
                return previous

            target_key = list(previous.keys())[nprd.randint(0, len(previous))]
            ret = previous.copy()

            ret[target_key] = Generation._random_value(ret[target_key])
            return ret
            # for k in previous.keys():
            #     previous[k] = Generation._random_value(previous[k])
            # return previous

        return None

    def mutation(self):
        """
        一様乱数で選択したものから、一様乱数で1つのパラメータを選択し、乱数に突然変異する
        """

        candidate = []
        candidate.extend(self._chosen)
        candidate.extend(self._previous)

        mutating = candidate[int(nprd.rand() * len(candidate))].copy()
        keys = list(mutating.keys())
        mutating_key = keys[int(nprd.rand() * len(keys))]

        mutating[mutating_key] = self._random_value(mutating[mutating_key])

        self._candidate.append(mutating)
        return mutating

    def selection(self):
        result = []

        if self.max_score <= self._previous_score:
            return result
        
        for i in range(self.len):
            if self._score[i] == self.max_score:
                result.append(self._candidate[i])
                break

        return result

    def run(self):
        window_size = 32
        for window in [self._candidate[i:i+window_size] for i in range(0, len(self._candidate), window_size)]:
            procs = []
            for c in window:
                request_list = [c]
                request_list.extend(self._chosen)

                procs.append((c, self._runner(request_list)))

            for req, p in procs:
                out, err = p.communicate()

                self._score.append(self._callback(out.decode('utf-8'), err.decode('utf-8'), req))

    @property
    def len(self):
        return len(self._candidate)

    @property
    def max_score(self):
        return max(self._score)

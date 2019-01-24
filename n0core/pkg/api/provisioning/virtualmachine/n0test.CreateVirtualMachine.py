import json
import subprocess
import numpy.random as nprd

from n0test.tdt_generator import Generation


ENV_KEY = "N0TEST_JSON_CreateVirtualMachine_REQUESTS"
# JSON_FILEPATH = "n0test.CreateVirtualMachine.json"
JSON_FILEPATH = "/tmp/n0test.selected.json"
JSON_SEED_FILEPATH = "n0test.CreateVirtualMachine.seed.json"
GENERATION_CASES = 512


class CreateVirtualMachineGeneration(Generation):
    def run_get_proc(self, request_list):
        p = subprocess.Popen([
            "virtualmachine.test",
            "-test.coverprofile=/dev/stderr"
        ], env={
            ENV_KEY: json.dumps(request_list),
        }, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

        return p

    def run_get_score(self, out, err, chosen):
        if "PASS" in out:
            lines = err.split('\n')
            coverage = 0

            for l in lines:
                if len(l) == 0:
                    continue

                if l[-1] == '1':
                    coverage += 1

            return coverage

        elif "N0TEST_OMIT" in out:
            return -1

        print("out={}\nerr={}\nrequest={}".format(out.decode('utf-8'), err.decode('utf-8'), json.dumps(request, indent=2)))

        # TODO: crash report, 保存したい
        return -1

    def select(self, select_num=256):
        filtered = filter(lambda x: x["score"] != -1, self._result)
        # selected = sorted(filtered, key=lambda x: x["distance"], reverse=True)[:select_num]
        # self._selected = list(map(lambda x: x["case"], selected)) + self._persistence
        self._selected = list(map(lambda x: x["case"], filtered)) + self._persistence
        # self._selected = list(map(lambda x: x["case"], selected))

        return self._selected


if __name__ == "__main__":
    # with open(JSON_FILEPATH) as f:
    #     prev = json.load(f)
    prev = None

    with open(JSON_SEED_FILEPATH) as f:
        seed = json.load(f)

    if prev:
        gen = CreateVirtualMachineGeneration(seed=seed, list_operate_distance=3, previous=prev)
    else:
        gen = CreateVirtualMachineGeneration(seed=seed, list_operate_distance=3)

    for i in range(500):
        while gen.len < GENERATION_CASES:
            if nprd.rand() < 0.75:
                gen.mutate(nprd.randint(1, 10) * (1 if nprd.rand() < 0.5 else -1))
                gen.mutate(nprd.randint(1, 10) * (1 if nprd.rand() < 0.5 else -1))
            else:
                gen.cross()

        gen.run()
        gen.persistent()
        gen.select(int(GENERATION_CASES * 0.75))

        print(json.dumps({
            "score(coverage)": gen.max_score,
            "persistenced": len(gen.persistence),
            "average_distance": gen.average_distance,
        }))

        with open(JSON_FILEPATH, mode='w') as f:
            json.dump(gen.persistence, f, indent=2)

        gen = CreateVirtualMachineGeneration(previous_generation=gen)
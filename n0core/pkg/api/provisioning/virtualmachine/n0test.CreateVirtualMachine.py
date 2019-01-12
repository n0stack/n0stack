import json
import subprocess
import numpy.random as nprd

from n0test.ga import Generation


ENV_KEY = "N0TEST_JSON_CreateVirtualMachine_REQUESTS"
JSON_FILEPATH = "n0test.CreateVirtualMachine.json"


def runner_CreateVirtualMachine(request_list):
    p = subprocess.Popen([
        "virtualmachine.test",
        "-test.coverprofile=/dev/stderr"
    ], env={
        ENV_KEY: json.dumps(request_list),
    }, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

    return p


def callback_CreateVirtualMachine(out, err, request):
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


if __name__ == "__main__":
    with open(JSON_FILEPATH) as f:
        seed = json.load(f)

    gen = Generation(runner_CreateVirtualMachine, callback_CreateVirtualMachine, chosen=seed)
    for i in range(100):
        while gen.len < 256:
            if nprd.rand() < 0.75:
                gen.mutation()
                gen.mutation()
            else:
                gen.cross()

        gen.run()
        print(json.dumps({
            "score": gen.max_score,
            "selected": len(gen._chosen),
        }))

        gen = Generation(runner_CreateVirtualMachine, callback_CreateVirtualMachine, chosen=gen.selection(), previous=gen)

    with open(JSON_FILEPATH, mode='w') as f:
        json.dump(gen._chosen, f, indent=2)

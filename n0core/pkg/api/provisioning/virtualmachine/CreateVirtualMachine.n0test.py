import json
import subprocess
import numpy.random as nprd

from n0test.ga import Generation


def runner_CreateVirtualMachine(request_list):
    p = subprocess.Popen([
        "virtualmachine.test",
        "-test.coverprofile=/dev/stderr"
    ], env={
        "N0TEST_JSON_CreateVirtualMachine_REQUESTS": json.dumps(request_list),
    }, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

    return p


def callback_CreateVirtualMachine(out, err, request):
    # print("out={}, err={}".format(out, err))
    # print("out={}, err={}".format(out, {}))

    if "PASS" in out:
        lines = err.split('\n')
        coverage = 0

        for l in lines:
            if len(l) == 0:
                continue

            if l[-1] == '1':
                coverage += 1

        # coverage = float(out.decode('utf-8').split("%")[0].split(" ")[1])
        # print("coverage: {}".format(coverage))
        return coverage

    elif "N0TEST_OMIT" in out:
        return -1

    # print(json.dumps(request_list))
    print("out={}\nerr={}\nrequest={}".format(out.decode('utf-8'), err.decode('utf-8'), json.dumps(request, indent=2)))

    # TODO: crash report, 保存したい
    return -1


if __name__ == "__main__":
    with open("CreateVirtualMachine.n0test.json") as f:
        seed = json.load(f)

    gen = Generation(runner_CreateVirtualMachine, callback_CreateVirtualMachine, chosen=seed)
    while True:
        while gen.len < 256:
            if nprd.rand() < 0.5:
                gen.mutation()
                gen.mutation()
            else:
                gen.cross()

        gen.run()
        print(json.dumps({
            "score": gen.max_score,
            "selected": len(gen._chosen),
        }))

        with open("/tmp/result", mode='w') as f:
            # json.dump(gen._chosen + gen._candidate, f)
            json.dump(gen._chosen, f)

        gen = Generation(runner_CreateVirtualMachine, callback_CreateVirtualMachine, chosen=gen.selection(), previous=gen)

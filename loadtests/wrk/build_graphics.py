import os

import matplotlib.pyplot as plt
import numpy as np

path = os.path.dirname(os.path.abspath(__file__))
output_dir = os.path.join(path, 'output')
users_dir = os.path.join(output_dir, 'users')

testing_module = 'search-users'
module_dir = os.path.join(users_dir, testing_module)


# testing_set = 'no_index_OR'
# testing_set = 'with_index_OR'

# testing_set = 'no_index_UNION'
testing_set = 'repl_idx_UNION'


files = os.listdir(os.path.join(module_dir, testing_set))

data = []
for fname in files:
    params = fname.split("-")
    if len(params) == 4 and params[1] == "d60" and fname.endswith('R1'):
        bar = {"conn": int([p[1:] for p in params if p.startswith('c')][0])}
        with open(os.path.join(os.path.join(module_dir, testing_set), fname), 'r') as f:
            for line in f.readlines():
                line = line.split()
                if line[0] == 'Latency':
                    bar['lat'] = line[3]
                    bar['lat'] = float(bar['lat'][:-2]) if bar['lat'].endswith('ms') or bar['lat'].endswith('us') \
                                                    else float(bar['lat'][:-1])*1000
                if line[0] == 'Requests/sec:': bar['req'] = float(line[1])
                if line[1] == 'requests': bar['total_reqs'] = int(line[0])
                if line[0] == 'Non-2xx': bar['err'] = int(line[-1])

        data.append(bar)

data = sorted(data, key=lambda d: d['conn'])

conns = [str(d['conn']) for d in data]
latency = [d['lat'] for d in data]
reqs = [d['req'] for d in data]
errs = [int(d.get('err')) if d.get('err') else 0 for d in data]
total_reqs = [int(d['total_reqs']) for d in data]
success_reqs = [t-e for t, e in zip(total_reqs, errs)]

plt.style.use('seaborn')
plt.figure(figsize=(12, 7))

plt.subplot(131)
plt.gca().set_title('Max latency (ms)')
plt.xlabel('connections')
plt.bar(conns, latency)

plt.subplot(132)
plt.gca().set_title('Req/sec')
plt.xlabel('connections')
plt.bar(conns, reqs)

plt.subplot(133)
plt.gca().set_title('Total requests')
plt.xlabel('connections')

p1 = plt.bar(conns, success_reqs, color='green')
p2 = plt.bar(conns, errs, hatch='//')


plt.legend((p1[0], p2[0]), ('Success requests', 'Non-2xx or 3xx responses'))

plt.suptitle(f'Latency & Requsts p/s [{testing_set}]')
plt.show()
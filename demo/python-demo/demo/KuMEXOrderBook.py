import platform
import os
import time
import sys

curPath = os.path.abspath(os.path.dirname(__file__))
rootPath = os.path.split(curPath)[0]
sys.path.append(rootPath)
from Level3.rpc import KuMEXRpc
from tools.config import rpc_config

cmd = ''
system_os = platform.system()

if system_os == 'Windows':
    cmd = 'cls'
elif system_os == 'Darwin' or system_os == 'Linux':
    cmd = 'clear'
else:
    raise Exception('unsupported system')

while True:
    kumex_rpc = KuMEXRpc(rpc_config['host'], rpc_config['port'], rpc_config['token'])
    data = kumex_rpc.get_ticker(100)
    asks = data['asks']
    bids = data['bids']
    price_list = [{}, {}]
    for ask in asks:
        if ask[1] not in price_list[0].keys():
            price_list[0].update({ask[1]: int(ask[2])})
        else:
            price_list[0].update({ask[1]: int(price_list[0][ask[1]]) + int(ask[2])})

        if len(price_list[0]) >= 13:
            price_list[0].pop(ask[1])
            break

    for bid in bids:
        if bid[1] not in price_list[1].keys():
            price_list[1].update({bid[1]: int(bid[2])})
        else:
            price_list[1].update({bid[1]: int(price_list[1][bid[1]]) + int(bid[2])})
        if len(price_list[1]) >= 13:
            price_list[1].pop(bid[1])
            break
    d1 = sorted(price_list[0].items(), key=lambda d: d[0], reverse=True)
    d2 = sorted(price_list[1].items(), key=lambda d: d[0], reverse=True)

    os.system(cmd)
    for d in d1:
        print("{} => {}".format(d[0], d[1]))

    print("---Spread---")
    for d in d2:
        print("{} => {}".format(d[0], d[1]))

    time.sleep(0.5)

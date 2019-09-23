import sys
import os
import json
import time
import uuid

curPath = os.path.abspath(os.path.dirname(__file__))
rootPath = os.path.split(curPath)[0]
sys.path.append(rootPath)

from tools.redis_client import RedisClient
from KumexApi.rest_api import KuMEXApi
from Level3.rpc import KuMEXRpc
from tools.config import redis_config, rpc_config, account_config


def deal_order(data, rdb, kumex_api):
    if data:
        data = json.loads(data)
        if data['type'] == 'match' and data['side'] == 'sell':
            rdb.client.set('matchOrder', json.dumps({'size': data['matchSize'], 'price': data['price']}))
            clientId = ''.join([each for each in str(uuid.uuid1()).split('-')])
            kumex_api.sandOrder('sell', '1', data['matchSize'], float(data['price']) + 1, clientOid=clientId)

    # TODO if you order is done go to hedger this order


if __name__ == '__main__':
    # channel definition by oneself and  Level3 Rpc add_event_client_id be the same
    kumex_rpc = KuMEXRpc(rpc_config['host'], rpc_config['port'], rpc_config['token'])
    kumex_api = KuMEXApi(account_config['key'], account_config['secret'],
                         account_config['passphrase'], False)
    channel = redis_config['clientId_channel']
    redis_client = RedisClient(redis_config['host'], redis_config['port'])
    rdb = redis_client.client
    ps = rdb.pubsub()
    ps.psubscribe([channel])

    for item in ps.listen():
        print(item)
        if item['type'] == 'pmessage':
            deal_order(item['data'], rdb, kumex_api)

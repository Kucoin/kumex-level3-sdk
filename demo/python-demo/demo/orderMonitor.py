import sys
import os
import json
import time
import uuid

curPath = os.path.abspath(os.path.dirname(__file__))
rootPath = os.path.split(curPath)[0]
sys.path.append(rootPath)

from KumexApi.rest_api import k_api
from Level3.rpc import krpc
from tools.redis_client import rdb
from tools.config import redis_config


def deal_order(data):
    if data:
        data = json.loads(data)

        if data['type'] == 'match':
            # judge side  buy or sell  you sandorder sell or buy
            rdb.set('matchOrder', json.dumps({'size': data['matchSize'], 'price': data['price']}))
            clientId = ''.join([each for each in str(uuid.uuid1()).split('-')])

            k_api.sandOrder('sell', '1', data['matchSize'], float(data['price']) + 1, clientOid=clientId)

    # TODO if you order is done go to hedger this order


if __name__ == '__main__':
    # channel definition by oneself and  Level3 Rpc add_event_client_id be the same


    channel = redis_config['clientId_channel']
    ps = rdb.pubsub()
    ps.psubscribe([channel])

    for item in ps.listen():
        print(item)
        if item['type'] == 'pmessage':
            deal_order(item['data'])

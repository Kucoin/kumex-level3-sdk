import json
import time
import uuid
import sys
import os
import random

curPath = os.path.abspath(os.path.dirname(__file__))
rootPath = os.path.split(curPath)[0]
sys.path.append(rootPath)
from KumexApi.rest_api import KuMEXApi
from Level3.rpc import KuMEXRpc
from tools.redis_client import RedisClient
from tools.config import redis_config, rpc_config, account_config


def main(kumex_rpc, kumex_api, price_dict):
    data = kumex_rpc.get_ticker(1)
    orderinfo = rdb.client.get('matchOrder')
    print(orderinfo)
    if orderinfo:
        orderinfo = json.loads(orderinfo)

    if data and data.get('asks') and data.get('bids'):
        ask1_price, bid1_price = int(data['asks'][0][1]), int(data['bids'][0][1])

        if not price_dict:
            sand_order(kumex_api, ask1_price, bid1_price, price_dict)
        else:
            if bid1_price != int(price_dict['price']):
                kumex_api.cancelOrder(price_dict['orderId'])
                price_dict.update({})
                time.sleep(0.1)
                sand_order(kumex_api, ask1_price, bid1_price, price_dict)
        if orderinfo and orderinfo.get('price') and orderinfo.get('price', 0) == price_dict['price'] and orderinfo.get(
                'size'):
            if orderinfo.get('size') == price_dict['size']:
                price_dict.update({})

            else:
                kumex_api.cancelOrder(price_dict[0])
                price_dict.update({})

                time.sleep(0.1)
                sand_order(kumex_api, ask1_price, bid1_price, price_dict)


def sand_order(kumex_api, ask1_price, bid1_price, price_dict):
    clientId = ''.join([each for each in str(uuid.uuid1()).split('-')])
    # price is int
    if ask1_price - bid1_price > 1:
        price = bid1_price + 1
    else:
        price = bid1_price
    kumex_rpc.add_event_client_id([clientId])
    size = random.randint(1, 5)
    orderId = kumex_api.sandOrder('buy', '1', size, price, clientOid=clientId)

    price_dict.update({'orderId': orderId, 'size': int(size), 'price': int(price)})


if __name__ == '__main__':
    rdb = RedisClient(redis_config['host'], redis_config['port'])
    kumex_rpc = KuMEXRpc(rpc_config['host'], rpc_config['port'], rpc_config['token'])
    kumex_api = KuMEXApi(account_config['key'], account_config['secret'],
                         account_config['passphrase'], False)
    # data is big use redis
    price_dict = {}
    while True:
        main(kumex_rpc, kumex_api, price_dict)
        time.sleep(0.2)

import json
import time
import uuid
import sys
import os
import random

curPath = os.path.abspath(os.path.dirname(__file__))
rootPath = os.path.split(curPath)[0]
sys.path.append(rootPath)
from KumexApi.rest_api import k_api
from Level3.rpc import krpc
from tools.redis_client import rdb
from tools.config import redis_config


def main():
    data = dict()
    try:
        data = krpc.get_ticker(1)
    except Exception as e:
        # cancel order
        pass
    orderinfo = rdb.get('matchOrder')
    if orderinfo:
        orderinfo = json.loads(orderinfo)

    if data and data.get('asks') and data.get('bids'):
        ask1_price, bid1_price = int(data['asks'][0][1]), int(data['bids'][0][1])

        if not price_dict:
            sand_order(ask1_price, bid1_price)
        else:
            if bid1_price != int(price_dict['price']):
                k_api.cancelOrder(price_dict['orderId'])
                price_dict.update({})
                time.sleep(0.1)
                sand_order(ask1_price, bid1_price,)
        if orderinfo and orderinfo.get('price') and orderinfo.get('price', 0) == price_dict['price'] and orderinfo.get(
                'size'):
            if orderinfo.get('size') == price_dict['size']:
                price_dict.update({})

            else:
                k_api.cancelOrder(price_dict[0])
                price_dict.update({})

                time.sleep(0.1)
                sand_order(ask1_price, bid1_price)
    else:
        # cancel order
        pass


def sand_order(ask1_price, bid1_price):
    clientId = ''.join([each for each in str(uuid.uuid1()).split('-')])
    # price is int
    if ask1_price - bid1_price > 1:
        price = bid1_price + 1
    else:
        price = bid1_price
    krpc.add_event_client_id([clientId], redis_config['clientId_channel'])
    size = random.randint(1, 5)
    orderId = k_api.sandOrder('buy', '1', size, price, clientOid=clientId)

    price_dict.update({'orderId': orderId, 'size': int(size), 'price': int(price)})


if __name__ == '__main__':

    # data is big use redis
    price_dict = {}
    while True:
        main()
        time.sleep(0.2)

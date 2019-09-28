import json
import requests
import hmac
import hashlib
import base64
import time
import uuid
import os
import sys

curPath = os.path.abspath(os.path.dirname(__file__))
rootPath = os.path.split(curPath)[0]
sys.path.append(rootPath)
from urllib.parse import urljoin
from tools.config import account_config



class KumexApi(object):

    def __init__(self, is_test=False):
        if is_test:
            self.url = 'https://sandbox-api.kumex.com'
        else:
            self.url = 'https://api.kumex.com'
        self.apiKey = account_config['key']
        self.secret = account_config['secret']
        self.Passphrase = account_config['passphrase']

    def requstdata(self, method, uri, params=None):
        url = urljoin(self.url, uri)
        if params:
            data_json = json.dumps(params)
        else:
            data_json = ''
        now_time = int(time.time()) * 1000
        str_to_sign = str(now_time) + method + uri + data_json
        sign = base64.b64encode(
            hmac.new(self.secret.encode('utf-8'), str_to_sign.encode('utf-8'), hashlib.sha256).digest())
        headers = {
            "KC-API-SIGN": sign,
            "KC-API-TIMESTAMP": str(now_time),
            "KC-API-KEY": self.apiKey,
            "KC-API-PASSPHRASE": self.Passphrase,
            "Content-Type": "application/json"
        }

        response_data = requests.request(method, url, headers=headers, data=data_json)
        if response_data and response_data.status_code == 200:
            data = response_data.json()
            if data and data.get('code'):
                if data.get('code') == '200000':
                    return data['data']
                else:
                    print('request err-{}'.format(json.dumps(data)))

    def getAccount(self):
        uri = '/api/v1/account-overview'
        method = 'GET'
        data = self.requstdata(method, uri)
        print(data)

    def getWsToker(self, is_public=True):
        if is_public:
            uri = '/api/v1/bullet-public'
        else:
            uri = '/api/v1/bullet-private'
        method = 'POST'
        data = self.requstdata(method, uri)
        print(data)

    def sandOrder(self, side, leverage, size, price, postOnly=True, symbol='XBTUSDM', clientOid=''):
        uri = '/api/v1/orders'
        method = 'POST'
        params = {
            'size': size,
            'side': side,
            'leverage': leverage,
            'symbol': symbol,
            'price': price,
            'postOnly': postOnly,
            'type': 'limit'
        }
        if clientOid:
            params['clientOid'] = clientOid
        orderId = self.requstdata(method, uri, params)
        return orderId

    def cancelOrder(self, orderId):
        uri = '/api/v1/orders/' + str(orderId)
        method = 'DELETE'
        self.requstdata(method, uri)

    def getTicker(self, symbol='XBTUSDM'):
        uri = '/api/v1/ticker' + '?symbol=' + symbol
        method = 'GET'
        ticker = self.requstdata(method, uri)
        # print(ticker)
        return ticker

    def getOrderBook(self, depth=1, symbol='XBTUSDM'):
        uri = '/api/v1/level3/snapshot' + '?symbol=' + symbol
        method = 'GET'
        orderBook = self.requstdata(method, uri)
        data = {}
        if orderBook:
            bids = orderBook['bids'][-depth]
            asks = orderBook['asks'][depth - 1]
            data['bids'] = [bids[2], bids[3]]
            data['asks'] = [asks[2], asks[3]]
        return data

k_api = KumexApi()


if __name__ == '__main__':
    pass

import json
import socket
import os
import sys

curPath = os.path.abspath(os.path.dirname(__file__))
rootPath = os.path.split(curPath)[0]
sys.path.append(rootPath)
from tools.config import rpc_config
# from tools.tools import singleton
from tools.tools import singleton

@singleton
class KuMEXRpc(object):
    def __init__(self, host, prot, token):
        self.prot = prot
        self.host = host
        self.token = token
        self.conn = None

    def get_connect(self):
        if not self.conn:
            self.conn = socket.create_connection((self.host, self.prot))

    def read_line(self):
        # return self.conn.makefile().readline()

        ret = b''
        while True:
            c = self.conn.recv(1)
            if c == b'\n' or c == b'':
                break
            else:
                ret += c

        return ret.decode("utf-8")

    def execute(self, data):
        data['id'] = 0
        msg = json.dumps(data)
        self.get_connect()
        self.conn.sendall(msg.encode())
        resp = self.read_line()
        resp = json.loads(resp)
        if resp.get('id') != 0:
            raise Exception("expected id=%s, received id=%s: %s"
                            % (0, resp.get('id'), resp.get('error')))
        if resp.get('error') is not None:
            raise Exception(resp.get('error'))
        data = json.loads(resp.get('result'))
        if data['code'] != '0':
            raise Exception("rpc get ticker fail: %s" % data['error'])
        return data

    def close(self):
        self.conn.close()

    def call(self, method, **kwargs):
        params = {
            'token': self.token,
        }
        # print(kwargs)
        if kwargs:
            params.update(kwargs)

        data = {
            'method': "Server." + method,
            'params': [params],
        }
        return self.execute(data)

    def get_ticker(self, num):
        data = self.call("GetPartOrderBook", number=num)
        ticker = json.loads(data['data'])
        if ticker['sequence'] == 0:
            raise Exception("rpc get ticker fail: sequence Is Null")
        return ticker

    def add_event_client_id(self, data, channel):
        datas = {}
        for i in data:
            datas[i] = [channel]
        return self.call("AddEventClientOidsToChannels", data=datas)

krpc = KuMEXRpc(rpc_config['host'], rpc_config['port'], rpc_config['token'])


if __name__ == '__main__':
    pass

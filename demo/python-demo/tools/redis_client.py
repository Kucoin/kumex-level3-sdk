import time
import aioredis
import redis
from tools.tools import singleton


@singleton
class RedisClient(object):
    def __init__(self, host, port, password=None, db=None):
        self.host = host
        self.port = port
        self.password = password
        self.db = db
        self.redis_client = None

    def init(self):
        pool = redis.ConnectionPool(host=self.host, port=self.port, password=self.password, db=self.db,
                                    socket_connect_timeout=15, socket_timeout=15)
        self.redis_client = redis.Redis(connection_pool=pool)

    @property
    def client(self):
        if not self.redis_client:
            for _ in range(2):
                self.init()
                time.sleep(2)
        if not self.redis_client:
            raise Exception('RedisClient is not')
        return self.redis_client


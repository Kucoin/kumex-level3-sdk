import time

import redis
import time
import redis
import os
import sys

curPath = os.path.abspath(os.path.dirname(__file__))
rootPath = os.path.split(curPath)[0]
sys.path.append(rootPath)
from tools.config import redis_config

rdb = redis.StrictRedis(host=redis_config['host'], port=redis_config['port'], db=None, password=None)

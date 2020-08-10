#!/usr/bin/env python3

import subprocess
import time
import json
import socket
import redis
from behave import *


register_type(int=int)
register_type(str=lambda x: x if x != "N/A" else "")
register_type(bool=lambda x: True if x == "true" else False)


config = {
    "prefix": "output/go-http",
    "service": {
        "port": 21680,
    },
    "redis": {
        "host": "127.0.0.1",
        "port": 16379,
        "password": '123456'
    }
}


def wait_for_port(port, host="localhost", timeout=5.0):
    start_time = time.perf_counter()
    while True:
        try:
            with socket.create_connection((host, port), timeout=timeout):
                break
        except OSError as ex:
            time.sleep(0.01)
            if time.perf_counter() - start_time >= timeout:
                raise TimeoutError("Waited too long for the port {} on host {} to start accepting connections.".format(
                    port, host
                )) from ex


def deploy():
    file = "{}/configs/server.json".format(config["prefix"])
    fp = open(file)
    cf = json.loads(fp.read())
    fp.close()
    cf["service"]["port"] = ':{}'.format(config["service"]["port"])
    cf["redis"]["addr"] = '{}:{}'.format(config["redis"]["host"], config["redis"]["port"])
    cf["redis"]["password"] = config["redis"]["password"]
    print(cf)
    fp = open(file, "w")
    fp.write(json.dumps(cf, indent=4))
    fp.close()


def start():
    subprocess.Popen(
        "cd {} && nohup bin/server &".format(config["prefix"]),  shell=True
    )
    wait_for_port(config["service"]["port"], timeout=5)


def stop():
    subprocess.getstatusoutput(
        "ps aux | grep bin/server | grep -v grep | awk '{print $2}' | xargs kill"
    )


# 在所有Scenario执行之前执行, 环境准备
def before_all(context):
    config["url"] = "http://127.0.0.1:{}".format(config["service"]["port"])
    deploy()
    start()
    context.config = config
    context.redis_client = redis.Redis(
        config["redis"]["host"], port=config["redis"]["port"], password=config["redis"]["password"]
    )


# 所有Scenario执行完成后执行
def after_all(context):
    # 清理工作，不影响下次测试
    context.redis_client.delete("go-http:17744581949")
    context.redis_client.close()
    stop()

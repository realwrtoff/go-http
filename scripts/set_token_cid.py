#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import getopt
from urllib import parse
import requests


def main():
    opts, args = getopt.getopt(sys.argv, "h", ["help"])
    if 2 > len(args):
        action = 'xx'
    else:
        action = args[1]
    #url = 'http://127.0.0.1:7060/{0}'.format(action)
    url = 'http://52.130.80.56:21668/{0}'.format(action)
    if action == 'channel':
        advertiser_url = 'http://advertise.equblock.com/click'
        params = {
            'pub': 'didazhuan',
            'cid': 'ddz_xxx',
            'status': 1,
            'advertiser_addr': advertiser_url,
            'advertiser_cid': 'guahao',
            'my_name': 'haina',
            'billing_type': 'active',
        }
    elif action == 'click':
        callback_url = 'http://channel.equblock.com/install?name=token&type=callback'
        params = {
            'pub': 'didazhuan',
            'cid': 'ddz_xxx',
            'ip': '192.168.1.111',
            'devicetype': 'iphone',
            'os': 'ios',
            'osversion': 'IOS-1.13',
            'idfa': 'fengmin-de-ipnone',
            'callback': callback_url
        }
    elif action in ['install', 'active']:
        params = {
            'click_id': args[2]
        }
    else:
        print('action {0} not support'.format(action))
        sys.exit(-1)

    print('request {0}'.format(url))
    print('params {0}'.format(params))
    r = requests.get(url, params=params)
    print('response status{0}, content [{1}]'.format(r.status_code, r.text))


if __name__ == '__main__':
    main()

Feature: geo 测试

  Scenario: geo distance
    When http 请求 GET /geo/distance
            """
            {
                "params": {
                    "latitude_1": 23.1378010917,
                    "longitude_1": 113.4022203113,
                    "latitude_2": 22.1191433172,
                    "longitude_2": 113.5826193044
                }
            }
            """
    Then http 检查 200
            """
            {
                "json": {
                    "data": 114.77281032762338
                }
            }
            """

  Scenario: geo address
    When http 请求 GET /geo/address
            """
            {
                "params": {
                    "address": "深圳湾科技生态园",
                    "city": "深圳"
                }
            }
            """
    Then http 检查 200
            """
            {
                "json": {
                    "code": 0,
                    "message": "query ok",
                    "data": {
                        "title":"深圳湾科技生态园",
                        "location": {
                            "lat":22.53091,
                            "lng":113.952385
                        },
                        "ad_info": {
                            "adcode": "440305"
                        },
                        "address_components": {
                            "province": "广东省",
                            "city": "深圳市",
                            "district": "南山区"
                        }
                    }
                }
            }
            """

  Scenario: geo adcode
    When http 请求 GET /geo/adcode
            """
            {
                "params": {
                    "latitude": 22.53091,
                    "longitude": 113.952385
                }
            }
            """
    Then http 检查 200
            """
            {
                "json": {
                    "code": 0,
                    "message": "query ok",
                    "data": {
                        "address":"广东省深圳市南山区高新南环路",
                        "location": {
                            "lat":22.53091,
                            "lng":113.952385
                        },
                        "ad_info": {
                            "nation_code":"156",
                            "adcode": "440305",
                            "city_code":"156440300",
                            "nation":"中国",
                            "province": "广东省",
                            "city": "深圳市",
                            "district": "南山区"
                        }
                    }
                }
            }
            """
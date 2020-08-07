Feature: echo 测试

    Scenario: echo
        When http 请求 GET /echo
            """
            {
                "params": {
                    "message": "hello world"
                }
            }
            """
        Then http 检查 200
            """
            {
                "json": {
                    "message": "hello world"
                }
            }
            """

    Scenario: location
        When http 请求 GET /location
            """
            {
                "params": {
                    "mobile": "1774458"
                }
            }
            """
        Then http 检查 200
            """
            {
                "json": {
                    "code": 200,
                    "message": "",
                    "data": {
                        "province": "北京",
                        "city": "北京",
                        "operator": "中国电信",
                        "post_code": "100000",
                        "area_code": "010"
                    }
                }
            }
            """
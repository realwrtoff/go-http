Feature: distinct 测试

  Scenario: distinct success
    When http 请求 GET /distinct
            """
            {
                "params": {
                    "business": "go-http",
                    "key": "17744581949"
                }
            }
            """
    Then http 检查 200
            """
            {
                "json": {
                    "code":200,
                    "message":"setnx go-http:17744581949 success"
                }
            }
            """
    Then redis exist "go-http:17744581949"


  Scenario: distinct fail
    When http 请求 GET /distinct
            """
            {
                "params": {
                    "business": "go-http",
                    "key": "17744581949"
                }
            }
            """
    Then http 检查 200
            """
            {
                "json": {
                    "code":200,
                    "message":"setnx go-http:17744581949 failed"
                }
            }
            """
    Then redis exist "go-http:17744581949"


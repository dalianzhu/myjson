package myjson

type testLongJsonStruct struct {
	Code int `json:"code"`
	Data struct {
		ChartInfo []struct {
			Cond struct {
				AppEnname      string   `json:"app_enname"`
				AppID          string   `json:"app_id"`
				AppMark        string   `json:"app_mark"`
				AppName        string   `json:"app_name"`
				BeginTime      string   `json:"begin_time"`
				CollectChartID int      `json:"collect_chart_id"`
				ConvergeTime   int      `json:"converge_time"`
				CustomTitle    string   `json:"custom_title"`
				EndTime        string   `json:"end_time"`
				Metric         []string `json:"metric"`
				MetricIDMap    struct {
					FrameworkRecvCounts int `json:"framework_recv_counts"`
				} `json:"metric_id_map"`
				MetricNameMap struct {
					FrameworkRecvCounts string `json:"framework_recv_counts"`
				} `json:"metric_name_map"`
				Order  int `json:"order"`
				TagSet struct {
					InstanceMark struct {
						IsNotEqual bool     `json:"is_not_equal"`
						Val        []string `json:"val"`
					} `json:"_instance_mark"`
					DataSourceType struct {
						IsNotEqual bool     `json:"is_not_equal"`
						Val        []string `json:"val"`
					} `json:"dataSourceType"`
				} `json:"tag_set"`
				Title string `json:"title"`
				UUID  string `json:"uuid"`
			} `json:"cond"`
			DetailDataList []struct {
				Current float64 `json:"current"`
				Time    string  `json:"time"`
			} `json:"detail_data_list"`
			Factor      float64 `json:"factor"`
			KeyDataList []struct {
				Current      float64 `json:"current"`
				CurrentTime  string  `json:"current_time,omitempty"`
				Desc         string  `json:"desc"`
				CurrentTotal int     `json:"current_total,omitempty"`
			} `json:"key_data_list"`
			Title string `json:"title"`
			Unit  string `json:"unit"`
		} `json:"chart_info"`
		MetricFactorDict struct {
			FrameworkRecvCounts float64 `json:"framework_recv_counts"`
		} `json:"metric_factor_dict"`
		MetricUnitDict struct {
			FrameworkRecvCounts string `json:"framework_recv_counts"`
		} `json:"metric_unit_dict"`
		PageNum  int `json:"page_num"`
		TotalNum int `json:"total_num"`
	} `json:"data"`
	Msg string `json:"msg"`
}

var longJsonVal = `
{
    "code": 0,
    "data": {
        "chart_info": [
            {
                "cond": {
                    "app_enname": "helloworld",
                    "app_id": "35",
                    "app_mark": "20_4533_helloworld",
                    "app_name": "测试服务",
                    "begin_time": "2020-12-10 11:32:21",
                    "collect_chart_id": 0,
                    "converge_time": 1,
                    "custom_title": "测试服务_framework接收数量",
                    "end_time": "2020-12-10 11:37:21",
                    "metric": [
                        "framework_recv_counts"
                    ],
                    "metric_id_map": {
                        "framework_recv_counts": 316
                    },
                    "metric_name_map": {
                        "framework_recv_counts": "framework接收数量"
                    },
                    "order": 0,
                    "tag_set": {
                        "_instance_mark": {
                            "is_not_equal": false,
                            "val": [
                                "Production"
                            ]
                        },
                        "dataSourceType": {
                            "is_not_equal": false,
                            "val": [
                                "this_is_a_test"
                            ]
                        }
                    },
                    "title": "测试服务_this_is_a_test_framework接收数量",
                    "uuid": "6e80de99-3a99-11eb-bd1d-6c0b84aed2314cqadaf99-3a99-11eb-bd1d-6c0b84aedcc50"
                },
                "detail_data_list": [
                    {
                        "current": 12007831.0,
                        "time": "2020-12-10 11:32"
                    },
                    {
                        "current": 12099001.0,
                        "time": "2020-12-10 11:33"
                    },
                    {
                        "current": 12259362.0,
                        "time": "2020-12-10 11:34"
                    },
                    {
                        "current": 12520209.0,
                        "time": "2020-12-10 11:35"
                    },
                    {
                        "current": 11648017.0,
                        "time": "2020-12-10 11:36"
                    },
                    {
                        "current": 11861225.0,
                        "time": "2020-12-10 11:37"
                    }
                ],
                "factor": 1.0,
                "key_data_list": [
                    {
                        "current": 12520209.0,
                        "current_time": "2020-12-10 11:35",
                        "desc": "最大值"
                    },
                    {
                        "current": 11648017.0,
                        "current_time": "2020-12-10 11:36",
                        "desc": "最小值"
                    },
                    {
                        "current": 11861225.0,
                        "current_time": "2020-12-10 11:37",
                        "desc": "最新值"
                    },
                    {
                        "current": 12065940.833333,
                        "desc": "平均值"
                    },
                    {
                        "current": 72395645.0,
                        "current_total": 6,
                        "desc": "累计值"
                    }
                ],
                "title": "测试服务_this_is_a_test_framework接收数量",
                "unit": "次"
            },
            {
                "cond": {
                    "app_enname": "helloworld",
                    "app_id": "35",
                    "app_mark": "20_4533_helloworld",
                    "app_name": "测试服务",
                    "begin_time": "2020-12-10 11:32:21",
                    "collect_chart_id": 0,
                    "converge_time": 1,
                    "custom_title": "测试服务_framework接收数量",
                    "end_time": "2020-12-10 11:37:21",
                    "metric": [
                        "framework_recv_counts"
                    ],
                    "metric_id_map": {
                        "framework_recv_counts": 316
                    },
                    "metric_name_map": {
                        "framework_recv_counts": "framework接收数量"
                    },
                    "order": 0,
                    "tag_set": {
                        "_instance_mark": {
                            "is_not_equal": false,
                            "val": [
                                "Production"
                            ]
                        },
                        "dataSourceType": {
                            "is_not_equal": false,
                            "val": [
                                "this_is_a_test"
                            ]
                        }
                    },
                    "title": "测试服务_this_is_a_test_framework接收数量",
                    "uuid": "6e80de99-3a99-11eb-bd1d-6asdf231431212349-3a99-11eb-bd1d-6c0b84aedcc50"
                },
                "detail_data_list": [
                    {
                        "current": 178691.0,
                        "time": "2020-12-10 11:32"
                    },
                    {
                        "current": 189804.0,
                        "time": "2020-12-10 11:33"
                    },
                    {
                        "current": 169061.0,
                        "time": "2020-12-10 11:34"
                    },
                    {
                        "current": 188334.0,
                        "time": "2020-12-10 11:35"
                    },
                    {
                        "current": 180259.0,
                        "time": "2020-12-10 11:36"
                    },
                    {
                        "current": 187791.0,
                        "time": "2020-12-10 11:37"
                    }
                ],
                "factor": 1.0,
                "key_data_list": [
                    {
                        "current": 189804.0,
                        "current_time": "2020-12-10 11:33",
                        "desc": "最大值"
                    },
                    {
                        "current": 169061.0,
                        "current_time": "2020-12-10 11:34",
                        "desc": "最小值"
                    },
                    {
                        "current": 187791.0,
                        "current_time": "2020-12-10 11:37",
                        "desc": "最新值"
                    },
                    {
                        "current": 182323.333333,
                        "desc": "平均值"
                    },
                    {
                        "current": 1093940.0,
                        "current_total": 6,
                        "desc": "累计值"
                    }
                ],
                "title": "测试服务_this_is_a_test_framework接收数量",
                "unit": "次"
            }
        ],
        "metric_factor_dict": {
            "framework_recv_counts": 1.0
        },
        "metric_unit_dict": {
            "framework_recv_counts": "次"
        },
        "page_num": 1,
        "total_num": 2
    },
    "msg": ""
}
`

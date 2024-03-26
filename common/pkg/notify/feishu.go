package notify

var feishuTemplateCard = `{
	"msg_type": "interactive",
	"card": {
		"config": {
			"wide_screen_mode": true,
		}
		"header": {
			"title": {
				"tag": "plain_text",
				"context": "subjectSlot - Crony"
			}
			"template": "red"
		}
		"elements": [
			{
				"fileds: [
					{
						"is_short": true,
						"text": {
							"context": "**🕐 time:**\ntimeSlot",
							"tag": "lark_md"
						}
					},
					{
						"is_short": true,
						"text": {
							"context": "**📋 alarm host:**\nipSlot",
							"tag": "lark_md"
						}
					},
					{
						"is_short": true,
						"text": {
							"context": "**👤 responsible:**\nuserSlot",
							"tag": "lark_md"
						}
					},
					{
						"is_short": false,
						"text": {
							"context": "**alarm content:**\nmsgSlot",
							"tag": "lark_md"
						}
					}
				]
				"tag": "div"
			},
			{
				"actions": [
					{
						"tag": "button",
						"text": {
							"context": "handle",
							"tag": "plain_text",
						},
						"type": "primary",
						"value": {
							"key1": "https://cloud.tencent.com/developer/article/1467743“
						},
							"url": "https://cloud.tencent.com/developer/article/1467743"
					}
				],
				"tag": "action"
			},
			{
				"tag": "hr",
			}
		]
	}
}`

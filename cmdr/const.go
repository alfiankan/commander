package cmdr

const (
	GreenTemplate = "\033[1;32m%s\033[0m"
	RedTemplate   = "\033[1;31m%s\033[0m"
	ChartRepoHost = "https://alfiankan.github.io/commander-charts/charts"
	ChartTemplate = `{
  "kind": "mychart",
  "description": "my personal chart",
  "charts": [
    {
      "usage": "git show log and statistic",
      "cmdt": "git log --stat",
      "type": "cmd",
      "prompt": []
    },
      {
      "usage": "load test apache benchmark",
      "cmdt": "ab -n {{total_req}} -c {{total_concurrent}} {{target_url}} ",
      "type": "cmd",
      "prompt": [
        {
          "tmplt": "total_req",
          "label": "total request",
          "default": "10"
        },
        {
          "tmplt": "total_concurrent",
          "label": "total concurrent",
          "default": "2"
        },
        {
          "tmplt": "target_url",
          "label": "url load test target",
          "default": "https://github.com/"
        }
      ]
    }
	]
}`
)

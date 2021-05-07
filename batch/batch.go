package batch

import "github.com/robfig/cron/v3"

func Start() {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("*/1 * * * * *", run)

	if err != nil {
		panic(err)
	}

	c.Start()
}

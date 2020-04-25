As we all know **HR** can manage jobs.
# features:
- run cron jobs that registered with [cron express](https://en.wikipedia.org/wiki/Cron#CRON_expression) syntax.
- when stops, will wait for running jobs to finish
- customizable `onStart` and `onEnd` 
- will not propagate panic which can be captured in `onEnd` as an error

# eg.
```go
hr := JobManager{}
hr.Register("hello job", "* * * * * *", func() {
    // job logic goes here
    log.Println("say hello once.")
})
hr.Start()
defer hr.Quit()
```

#dependency
- [cronexpr](https://github.com/gorhill/cronexpr) for parsing cron expression.

# TODO 

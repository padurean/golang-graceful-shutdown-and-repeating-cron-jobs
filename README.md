# Graceful shutdown with repeating "cron" jobs (running at a regular interval) in Go

Illustrates how to implement the following in Go:

- [x] run functions ("jobs) at a specified interval
- [x] gracefully shutdown a running process (and wait for any ongoing "job" runs to finish)

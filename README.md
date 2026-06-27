# Go Health Checker

Small concurrent health checker written in Go. It performs HTTP GET requests against a list of URLs and reports status, response time, content type/length and TLS expiry when available.

Quick usage

- Run normally:

```bash
make run
```

- Run in dev mode (prints a DEV marker):

```bash
make run-dev
```

What it does

- Loads the target URLs from `constants/urls.go`.
- Spawns a goroutine per URL to perform the check concurrently.
- Collects results with an indexed channel so output preserves the original URL order.
- Reports HTTP status code, elapsed time (ms), `Content-Type`, `Content-Length`, and TLS certificate expiry (if HTTPS).


Notes

- For long-running or higher-scale checks, consider adding a worker pool, retries, and better error handling.
- Feel free to edit `constants/urls.go` to change the checked sites.

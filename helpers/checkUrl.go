package helpers

	import (
		"fmt"
		"io"
		"net/http"
		"time"
	)

	type Result struct {
		Index int
		Msg   string
	}

	// CheckURL performs an HTTP GET with retry logic (up to 3 attempts).
	// It reports a Result containing an indexed message that includes the
	// number of attempts when relevant.
	func CheckURL(url string, i int, c chan Result) {
		const maxAttempts = 3
		var lastErr error
		var resp *http.Response
		var attemptsUsed int

		for attempt := 1; attempt <= maxAttempts; attempt++ {
			attemptsUsed = attempt
			client := &http.Client{Timeout: 5 * time.Second}
			start := time.Now()
			r, err := client.Get(url)
			elapsed := time.Since(start)

			if err != nil {
				lastErr = err
				if attempt < maxAttempts {
					time.Sleep(time.Duration(attempt) * 300 * time.Millisecond)
					continue
				}
				// final failure after retries
				c <- Result{i, fmt.Sprintf("❌ %s is down after %d attempts: %v", url, attempt, err)}
				return
			}

			// got a response
			resp = r

			// treat server error (5xx) as retryable
			if resp.StatusCode >= 500 && attempt < maxAttempts {
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				time.Sleep(time.Duration(attempt) * 300 * time.Millisecond)
				continue
			}

			// success (or non-retryable status)
			// measure elapsed from this successful request
			_ = elapsed
			break
		}

		if resp == nil {
			// no response and already reported final error above, but keep safe fallback
			c <- Result{i, fmt.Sprintf("❌ %s is down after %d attempts: %v", url, maxAttempts, lastErr)}
			return
		}

		// ensure body is closed and drained to allow connection reuse
		defer resp.Body.Close()
		io.Copy(io.Discard, resp.Body)

		status := resp.StatusCode
		contentType := resp.Header.Get("Content-Type")
		contentLen := resp.ContentLength

		tlsExpiry := ""
		if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
			tlsExpiry = resp.TLS.PeerCertificates[0].NotAfter.Format(time.RFC3339)
		}

		// include attempts info when more than one attempt was needed
		attemptsInfo := ""
		if attemptsUsed > 1 {
			attemptsInfo = fmt.Sprintf(" attempts=%d", attemptsUsed)
		}

		msg := fmt.Sprintf("✅ %s is up! status=%d type=%s len=%d%s", url, status, contentType, contentLen, attemptsInfo)
		if tlsExpiry != "" {
			msg = msg + " tls_expiry=" + tlsExpiry
		}

		c <- Result{i, msg}
	}
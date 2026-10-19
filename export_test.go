package flexops

import "time"

// SetRetryBaseDelay shortens the retry backoff so retry tests run fast.
func SetRetryBaseDelay(c *Client, d time.Duration) { c.http.baseDelay = d }

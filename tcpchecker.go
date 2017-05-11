package gopatrol

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

// TCPChecker implements a Checker for TCP endpoints.
type TCPChecker struct {
	Slug string `json:"slug" valid:"required"`
	// Name is the name of the endpoint.
	Name string `json:"name,omitempty" valid:"required"`
	// URL is the URL of the endpoint.
	URL  string `json:"url,omitempty" valid:"required"`
	Type string `json:"type,omitempty" valid:"required"`
	// ThresholdRTT is the maximum round trip time to
	// allow for a healthy endpoint. If non-zero and a
	// request takes longer than ThresholdRTT, the
	// endpoint will be considered unhealthy. Note that
	// this duration includes any in-between network
	// latency.
	ThresholdRTT time.Duration `json:"threshold_rtt,omitempty"`
	// Attempts is how many requests the client will
	// make to the endpoint in a single check.
	Attempts int `json:"attempts,omitempty"`
	// TLSEnabled controls whether to enable TLS or not.
	// If set, TLS is enabled.
	TLSEnabled bool `json:"tls,omitempty"`
	// TLSSkipVerify controls whether to skip server TLS
	// certificat validation or not.
	TLSSkipVerify bool `json:"tls_skip_verify,omitempty"`
	// TLSCAFile is the Certificate Authority used
	// to validate the server TLS certificate.
	TLSCAFile string `json:"tls_ca_file,omitempty"`
	// Timeout is the maximum time to wait for a
	// TCP connection to be established.
	Timeout     time.Duration `json:"timeout,omitempty"`
	LastChecked time.Time     `json:"last_checked"`
	LastStatus  string        `json:"last_status"`
}

func (c TCPChecker) GetName() string {
	return c.Name
}

func (c TCPChecker) GetURL() string {
	return c.URL
}

func (c TCPChecker) GetSlug() string {
	return c.Slug
}

func (c TCPChecker) GetLastChecked() time.Time {
	return c.LastChecked
}

func (c TCPChecker) GetLastStatus() string {
	return c.LastStatus
}

// Check performs checks using c according to its configuration.
// An error is only returned if there is a configuration error.
func (c TCPChecker) Check() (Result, error) {
	if c.Attempts < 1 {
		c.Attempts = 1
	}

	result := Result{Name: c.Name, URL: c.URL, Timestamp: time.Now().UTC(), Slug: c.Slug}
	result.Times = c.doChecks()

	return c.conclude(result), nil
}

// doChecks executes and returns each attempt.
func (c TCPChecker) doChecks() Attempts {
	var err error
	var conn net.Conn

	timeout := c.Timeout
	if timeout == 0 {
		timeout = 1 * time.Second
	}

	checks := make(Attempts, c.Attempts)
	for i := 0; i < c.Attempts; i++ {
		start := time.Now()

		if c.TLSEnabled {
			// Dialer with timeout
			dialer := &net.Dialer{
				Timeout: timeout,
			}

			// TLS config based on configuration
			var tlsConfig tls.Config
			tlsConfig.InsecureSkipVerify = c.TLSSkipVerify
			if c.TLSCAFile != "" {
				rootPEM, err := ioutil.ReadFile(c.TLSCAFile)
				if err != nil || rootPEM == nil {
					checks[i].Error = "failed to read root certificate"
				}
				pool := x509.NewCertPool()
				ok := pool.AppendCertsFromPEM([]byte(rootPEM))
				if !ok {
					checks[i].Error = "failed to parse root certificate"
				}
				tlsConfig.RootCAs = pool
			}
			if conn, err = tls.DialWithDialer(dialer, "tcp", c.URL, &tlsConfig); err == nil {
				conn.Close()
			}
		} else {
			if conn, err = net.DialTimeout("tcp", c.URL, c.Timeout); err == nil {
				conn.Close()
			}
		}
		checks[i].RTT = time.Since(start)
		if err != nil {
			checks[i].Error = err.Error()
			continue
		}
	}
	return checks
}

// conclude takes the data in result from the attempts and
// computes remaining values needed to fill out the result.
// It detects degraded (high-latency) responses and makes
// the conclusion about the result's status.
func (c TCPChecker) conclude(result Result) Result {
	result.ThresholdRTT = c.ThresholdRTT

	// Check errors (down)
	for i := range result.Times {
		if result.Times[i].Error != "" {
			result.Down = true
			result.Message = result.URL + " is down"
			return result
		}
	}

	// Check round trip time (degraded)
	if c.ThresholdRTT > 0 {
		stats := result.ComputeStats()
		if stats.Median > c.ThresholdRTT {
			result.Notice = fmt.Sprintf("median round trip time exceeded threshold (%s)", c.ThresholdRTT)
			result.Degraded = true
			return result
		}
	}

	result.Healthy = true
	result.Message = result.URL + " is healthy"
	return result
}

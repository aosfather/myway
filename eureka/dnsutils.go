package eureka

import (
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/franela/goreq"
	"github.com/miekg/dns"
	"github.com/op/go-logging"
	"time"
)

var log = logging.MustGetLogger("eureka")
var metadataLog = logging.MustGetLogger("eureka.metadata")
var marshalLog = logging.MustGetLogger("eureka.marshal")

func init() {
	logging.SetLevel(logging.WARNING, "eureka.metadata")
	logging.SetLevel(logging.WARNING, "eureka.marshal")
}

type unsuccessfulHTTPResponse struct {
	statusCode    int
	messagePrefix string
}

func (u *unsuccessfulHTTPResponse) Error() string {
	if len(u.messagePrefix) > 0 {
		return fmt.Sprint(u.messagePrefix, ", rcode = ", u.statusCode)
	}
	return fmt.Sprint("rcode = ", u.statusCode)
}

// HTTPResponseStatusCode extracts the HTTP status code for the response from Eureka that motivated
// the supplied error, if any. If the returned present value is true, the returned code is an HTTP
// status code.
func HTTPResponseStatusCode(err error) (code int, present bool) {
	if u, ok := err.(*unsuccessfulHTTPResponse); ok {
		return u.statusCode, true
	}
	return 0, false
}

type AppNotFoundError struct {
	specific string
}

func (e AppNotFoundError) Error() string {
	return "Application not found for name=" + e.specific
}

const azURL = "http://169.254.169.254/latest/meta-data/placement/availability-zone"

var ErrNotInAWS = fmt.Errorf("Not in AWS")

func discoverDNS(domain string, port int, urlBase string) (servers []string, ttl time.Duration, err error) {
	r, _ := region()

	// all DNS queries must use the FQDN
	domain = "txt." + r + "." + dns.Fqdn(domain)
	if _, ok := dns.IsDomainName(domain); !ok {
		err = fmt.Errorf("invalid domain name: '%s' is not a domain name", domain)
		return
	}
	regionRecords, ttl, err := retryingFindTXT(domain)
	if err != nil {
		return
	}

	for _, az := range regionRecords {
		instances, _, er := retryingFindTXT("txt." + dns.Fqdn(az))
		if er != nil {
			continue
		}
		for _, instance := range instances {
			// format the service URL
			servers = append(servers, fmt.Sprintf("http://%s:%d/%s", instance, port, urlBase))
		}
	}
	return
}

// retryingFindTXT will, on any DNS failure, retry for up to 15 minutes before
// giving up and returning an empty []string of records
func retryingFindTXT(fqdn string) (records []string, ttl time.Duration, err error) {
	err = backoff.Retry(
		func() error {
			records, ttl, err = findTXT(fqdn)
			if err != nil {
				log.Errorf("Retrying DNS query. Query failed with: %s", err.Error())
			}
			return err
		}, backoff.NewExponentialBackOff())
	return
}

func findTXT(fqdn string) ([]string, time.Duration, error) {
	defaultTTL := 120 * time.Second
	query := new(dns.Msg)
	query.SetQuestion(fqdn, dns.TypeTXT)
	dnsServerAddr, err := findDnsServerAddr()
	if err != nil {
		log.Errorf("Failure finding DNS server, err=%s", err.Error())
		return nil, defaultTTL, err
	}

	response, err := dns.Exchange(query, dnsServerAddr)
	if err != nil {
		log.Errorf("Failure resolving name %s err=%s", fqdn, err.Error())
		return nil, defaultTTL, err
	}
	if len(response.Answer) < 1 {
		err := fmt.Errorf("no Eureka discovery TXT record returned for name=%s", fqdn)
		log.Errorf("no answer for name=%s err=%s", fqdn, err.Error())
		return nil, defaultTTL, err
	}
	if response.Answer[0].Header().Rrtype != dns.TypeTXT {
		err := fmt.Errorf("did not receive TXT record back from query specifying TXT record. This should never happen.")
		log.Errorf("Failure resolving name %s err=%s", fqdn, err.Error())
		return nil, defaultTTL, err
	}
	txt := response.Answer[0].(*dns.TXT)
	ttl := response.Answer[0].Header().Ttl
	if ttl < 60 {
		ttl = 60
	}

	return txt.Txt, time.Duration(ttl) * time.Second, nil
}

func findDnsServerAddr() (string, error) {
	// Find a DNS server using the OS resolv.conf
	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		log.Errorf("Failure finding DNS server address from /etc/resolv.conf, err = %s", err)
		return "", err
	} else {
		return config.Servers[0] + ":" + config.Port, nil
	}
}

func region() (string, error) {
	zone, err := availabilityZone()
	if err != nil {
		log.Errorf("Could not retrieve availability zone err=%s", err.Error())
		return "us-east-1", err
	}
	return zone[:len(zone)-1], nil
}

// defaults to us-east-1 if there's a problem
func availabilityZone() (string, error) {
	response, err := goreq.Request{Uri: azURL}.Do()
	if err != nil {
		return "", err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := response.Body.ToString()
		return "", fmt.Errorf("bad response code: code %d does not indicate successful request, body=%s",
			response.StatusCode,
			body,
		)
	}
	zone, err := response.Body.ToString()
	if err != nil {
		return "", err
	}
	return zone, nil
}

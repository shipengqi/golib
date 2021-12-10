package sys

import "os"

func fqdn() (string, error) {
	return os.Hostname()
}

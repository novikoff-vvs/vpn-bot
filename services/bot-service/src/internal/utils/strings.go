package utils

import (
	"fmt"
	"os"
)

func BuildVlessLink(uuid string) string {
	return fmt.Sprintf(os.Getenv("VPN_SERVICE_BASE_SUBSCRIPTION_URL")+"/%s", uuid)
}

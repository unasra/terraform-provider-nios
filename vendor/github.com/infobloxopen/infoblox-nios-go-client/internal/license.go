package internal

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/infobloxopen/universal-ddi-go-client/inframgmt"
	"github.com/infobloxopen/universal-ddi-go-client/option"
)

var (
	licenseUIDCache   = map[string]string{}
	licenseUIDCacheMu sync.RWMutex
)

// fetchLicenseUID calls the Universal DDI Detail Hosts API to retrieve the
// host/license_uid tag for the NIOS host identified by its virtualnode_ip.
// Results are cached so the API is called at most once per (apiKey, hostIP) pair.
func fetchLicenseUID(apiKey, hostIP string) (string, error) {
	cacheKey := apiKey + "|" + hostIP

	licenseUIDCacheMu.RLock()
	if uid, ok := licenseUIDCache[cacheKey]; ok {
		licenseUIDCacheMu.RUnlock()
		//log.Printf("[fetchLicenseUID] returning cached license_uid for host %s", hostIP)
		return uid, nil
	}
	licenseUIDCacheMu.RUnlock()

	log.Printf("[fetchLicenseUID] creating inframgmt client with CSP URL: %s", cspURL)

	client := inframgmt.NewAPIClient(
		option.WithCSPUrl(cspURL),
		option.WithAPIKey(apiKey),
		option.WithDebug(true),
	)

	tfilter := fmt.Sprintf("\"host/virtualnode_ip\"=='%s'", hostIP)
	log.Printf("[fetchLicenseUID] calling detail hosts API with tfilter: %s", tfilter)

	resp, httpResp, err := client.DetailAPI.HostsList(context.Background()).
		Tfilter(tfilter).
		Execute()
	if err != nil {
		if httpResp != nil {
			log.Printf("[fetchLicenseUID] HTTP %d response from %s", httpResp.StatusCode, httpResp.Request.URL)
		}
		return "", fmt.Errorf("detail hosts API call failed: %w", err)
	}

	log.Printf("[fetchLicenseUID] HTTP %d from %s", httpResp.StatusCode, httpResp.Request.URL)

	results := resp.GetResults()
	log.Printf("[fetchLicenseUID] got %d result(s)", len(results))

	if len(results) == 0 {
		return "", fmt.Errorf("no host found with virtualnode_ip=%q", hostIP)
	}

	tags := results[0].GetTags()
	log.Printf("[fetchLicenseUID] host tags: %v", tags)

	licenseUID, ok := tags["host/license_uid"]
	if !ok {
		return "", fmt.Errorf("tag host/license_uid not found on host (available tags: %v)", tags)
	}

	uid, ok := licenseUID.(string)
	if !ok {
		return "", fmt.Errorf("tag host/license_uid has unexpected type %T", licenseUID)
	}

	licenseUIDCacheMu.Lock()
	licenseUIDCache[cacheKey] = uid
	licenseUIDCacheMu.Unlock()

	log.Printf("[fetchLicenseUID] resolved and cached license_uid for host %s", hostIP)
	return uid, nil
}

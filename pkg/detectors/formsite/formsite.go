package formsite

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/detectorspb"
)

type Scanner struct{}

// Ensure the Scanner satisfies the interface at compile time
var _ detectors.Detector = (*Scanner)(nil)

var (
	client = common.SaneHttpClient()

	//Make sure that your group is surrounded in boundry characters such as below to reduce false positives
	keyPat    = regexp.MustCompile(detectors.PrefixRegex([]string{"formsite"}) + `\b([a-zA-Z0-9]{32})\b`)
	serverPat = regexp.MustCompile(detectors.PrefixRegex([]string{"formsite"}) + `\b(fs[0-9]{1,4})\b`)
	userPat   = regexp.MustCompile(detectors.PrefixRegex([]string{"formsite"}) + `\b([a-zA-Z0-9]{6})\b`)
)

// Keywords are used for efficiently pre-filtering chunks.
// Use identifiers in the secret preferably, or the provider name.
func (s Scanner) Keywords() []string {
	return []string{"formsite"}
}

// FromData will find and optionally verify Formsite secrets in a given set of bytes..
func (s Scanner) FromData(ctx context.Context, verify bool, data []byte) (results []detectors.Result, err error) {
	dataStr := string(data)

	matches := keyPat.FindAllStringSubmatch(dataStr, -1)
	serverMatches := serverPat.FindAllStringSubmatch(dataStr, -1)
	userMatches := userPat.FindAllStringSubmatch(dataStr, -1)

	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		resMatch := strings.TrimSpace(match[1])
		for _, serverMatch := range serverMatches {
			if len(serverMatch) != 2 {
				continue
			}
			resServerMatch := strings.TrimSpace(serverMatch[1])
			for _, userMatch := range userMatches {
				if len(userMatch) != 2 {
					continue
				}
				resUserMatch := strings.TrimSpace(userMatch[1])
				s1 := detectors.Result{
					DetectorType: detectorspb.DetectorType_Formsite,
					Raw:          []byte(resMatch),
				}

				if verify {
					req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://%s.formsite.com/api/v2/%s/forms", resServerMatch, resUserMatch), nil)
					if err != nil {
						continue
					}
					req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", resMatch))
					res, err := client.Do(req)
					if err == nil {
						defer res.Body.Close()
						if res.StatusCode >= 200 && res.StatusCode < 300 {
							s1.Verified = true
						} else {
							//This function will check false positives for common test words, but also it will make sure the key appears 'random' enough to be a real key
							if detectors.IsKnownFalsePositive(resMatch, detectors.DefaultFalsePositives, true) {
								continue
							}
						}
					}
				}

				results = append(results, s1)
			}
		}

	}

	return detectors.CleanResults(results), nil
}

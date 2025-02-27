package common

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const EmailPattern = `\b(?:[a-z0-9!#$%&'*+/=?^_\x60{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_\x60{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])\b`
const SubDomainPattern = `\b([A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?)\b`
const UUIDPattern = `\b([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})\b`
const UUIDPatternUpperCase = `\b([0-9A-Z]{8}-[0-9A-Z]{4}-[0-9A-Z]{4}-[0-9A-Z]{4}-[0-9A-Z]{12})\b`

const RegexPattern = "0-9a-z"
const AlphaNumPattern = "0-9a-zA-Z"
const HexPattern = "0-9a-f"

// Custom Regex functions
func BuildRegex(pattern string, specialChar string, length int) string {
	return fmt.Sprintf(`\b([%s%s]{%s})\b`, pattern, specialChar, strconv.Itoa(length))
}

func BuildRegexJWT(firstRange, secondRange, thirdRange string) string {
	if RangeValidation(firstRange) || RangeValidation(secondRange) || RangeValidation(thirdRange) {
		log.Error("Min value should not be greater than or equal to max")
	}
	return fmt.Sprintf(`\b(ey[%s]{%s}.ey[%s-\/_]{%s}.[%s-\/_]{%s})\b`, AlphaNumPattern, firstRange, AlphaNumPattern, secondRange, AlphaNumPattern, thirdRange)
}

func RangeValidation(rangeInput string) bool {
	range_split := strings.Split(rangeInput, ",")
	range_min, _ := strconv.ParseInt(strings.TrimSpace(range_split[0]), 10, 0)
	range_max, _ := strconv.ParseInt(strings.TrimSpace(range_split[1]), 10, 0)
	return range_min >= range_max
}

func ToUpperCase(input string) string {
	return strings.ToUpper(input)
}

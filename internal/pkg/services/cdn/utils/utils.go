package utils

import (
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

func ParseGeofencing(p *print.Printer, geofencingInput []string) *map[string][]string { //nolint:gocritic // convenient for setting the SDK payload
	geofencing := make(map[string][]string)
	for _, in := range geofencingInput {
		firstSpace := strings.IndexRune(in, ' ')
		if firstSpace == -1 {
			p.Debug(print.ErrorLevel, "invalid geofencing entry (no space found): %q", in)
			continue
		}
		urlPart := in[:firstSpace]
		countriesPart := in[firstSpace+1:]
		geofencing[urlPart] = nil
		countries := strings.Split(countriesPart, ",")
		for _, country := range countries {
			country = strings.TrimSpace(country)
			geofencing[urlPart] = append(geofencing[urlPart], country)
		}
	}
	return &geofencing
}

func ParseOriginRequestHeaders(p *print.Printer, originRequestHeadersInput []string) *map[string]string { //nolint:gocritic // convenient for setting the SDK payload
	originRequestHeaders := make(map[string]string)
	for _, in := range originRequestHeadersInput {
		parts := strings.Split(in, ":")
		if len(parts) != 2 {
			p.Debug(print.ErrorLevel, "invalid origin request header entry (no colon found): %q", in)
			continue
		}
		originRequestHeaders[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return &originRequestHeaders
}

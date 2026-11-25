package bucket

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

// URIPrefix is the de-facto standard object storage URI scheme, accepted in
// addition to STACKIT Object Storage URLs
const URIPrefix = "s3://"

type Reference struct {
	Bucket string
	Key    string
	// Endpoint is the data-plane base endpoint (scheme + host) extracted from a
	// STACKIT https bucket URL. Empty for s3:// references, where it is derived
	// from the region.
	Endpoint string
	// Region is extracted from a STACKIT https bucket URL. Empty for s3:// references.
	Region string
}

var (
	bucketNameRegex    = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)
	virtualHostRegex   = regexp.MustCompile(`^([a-z0-9][a-z0-9.-]{1,61}[a-z0-9])\.object\.storage\.([a-z0-9-]+)\.onstackit\.cloud$`)
	pathStyleHostRegex = regexp.MustCompile(`^object\.storage\.([a-z0-9-]+)\.onstackit\.cloud$`)
)

func IsBucketURI(s string) bool {
	return strings.HasPrefix(s, URIPrefix)
}

func IsObjectRef(s string) bool {
	return IsBucketURI(s) || strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func (r Reference) IsBucket() bool {
	return r.Key == ""
}

func (r Reference) String() string {
	base := URIPrefix + r.Bucket
	if r.Endpoint != "" {
		base = r.Endpoint + "/" + r.Bucket
	}
	if r.IsBucket() {
		return base
	}
	return base + "/" + r.Key
}

func ParseBucketURI(uri string) (Reference, error) {
	switch {
	case IsBucketURI(uri):
		return parseBucketSchemeURI(uri)
	case strings.HasPrefix(uri, "http://"), strings.HasPrefix(uri, "https://"):
		return parseBucketURL(uri)
	default:
		return Reference{}, fmt.Errorf("%q is not a bucket URI (expected the %q prefix) or a STACKIT Object Storage URL (https://...)", uri, URIPrefix)
	}
}

func parseBucketSchemeURI(uri string) (Reference, error) {
	rest := strings.TrimPrefix(uri, URIPrefix)
	if rest == "" {
		return Reference{}, fmt.Errorf("missing bucket in URI %q", uri)
	}
	parts := strings.SplitN(rest, "/", 2)
	bucket := parts[0]
	if err := ValidateBucketName(bucket); err != nil {
		return Reference{}, fmt.Errorf("invalid bucket in URI %q: %w", uri, err)
	}
	key := ""
	if len(parts) == 2 {
		key = strings.TrimLeft(parts[1], "/")
	}
	return Reference{Bucket: bucket, Key: key}, nil
}

func parseBucketURL(uri string) (Reference, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return Reference{}, fmt.Errorf("parse URL %q: %w", uri, err)
	}
	host := u.Hostname()
	if host == "" {
		return Reference{}, fmt.Errorf("missing host in URL %q", uri)
	}
	base := u.Scheme + "://" + host
	path := strings.TrimLeft(u.Path, "/")

	// Virtual hosted style: https://<bucket>.object.storage.<region>.onstackit.cloud[/<key>]
	// The client always uses path style, so the bucket label is stripped from the host.
	if m := virtualHostRegex.FindStringSubmatch(host); m != nil {
		bucket := m[1]
		if err := ValidateBucketName(bucket); err != nil {
			return Reference{}, fmt.Errorf("invalid bucket in URL %q: %w", uri, err)
		}
		return Reference{
			Bucket:   bucket,
			Key:      path,
			Endpoint: u.Scheme + "://" + host[len(bucket)+1:],
			Region:   m[2],
		}, nil
	}

	// Path style: https://object.storage.<region>.onstackit.cloud/<bucket>[/<key>]
	if m := pathStyleHostRegex.FindStringSubmatch(host); m != nil {
		parts := strings.SplitN(path, "/", 2)
		if parts[0] == "" {
			return Reference{}, fmt.Errorf("missing bucket in URL %q", uri)
		}
		if err := ValidateBucketName(parts[0]); err != nil {
			return Reference{}, fmt.Errorf("invalid bucket in URL %q: %w", uri, err)
		}
		key := ""
		if len(parts) == 2 {
			key = parts[1]
		}
		return Reference{
			Bucket:   parts[0],
			Key:      key,
			Endpoint: base,
			Region:   m[1],
		}, nil
	}

	return Reference{}, fmt.Errorf("%q is not a valid STACKIT Object Storage URL (expected https://BUCKET.object.storage.REGION.onstackit.cloud[/KEY] or https://object.storage.REGION.onstackit.cloud/BUCKET[/KEY])", uri)
}

func ValidateBucketName(name string) error {
	if len(name) < 3 || len(name) > 63 {
		return fmt.Errorf("bucket name must be between 3 and 63 characters (got %d)", len(name))
	}
	if !bucketNameRegex.MatchString(name) {
		return fmt.Errorf("bucket name %q is invalid (must contain only lowercase letters, digits, dots and hyphens, and start and end with a lowercase letter or digit)", name)
	}
	if strings.Contains(name, "..") {
		return fmt.Errorf("bucket name %q must not contain consecutive dots", name)
	}
	if strings.Contains(name, ".-") || strings.Contains(name, "-.") {
		return fmt.Errorf("bucket name %q must not contain dots and hyphens next to each other", name)
	}
	if ip := net.ParseIP(name); ip != nil && ip.To4() != nil {
		return fmt.Errorf("bucket name %q must not look like an IP address", name)
	}
	return nil
}

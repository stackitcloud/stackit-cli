package bucket

import (
	"encoding/xml"
	"time"
)

type listAllMyBucketsResult struct {
	XMLName xml.Name      `xml:"ListAllMyBucketsResult"`
	Buckets []bucketEntry `xml:"Buckets>Bucket"`
}

type bucketEntry struct {
	Name         string    `xml:"Name"`
	CreationDate time.Time `xml:"CreationDate"`
}

type listBucketResult struct {
	XMLName               xml.Name       `xml:"ListBucketResult"`
	IsTruncated           bool           `xml:"IsTruncated"`
	NextContinuationToken string         `xml:"NextContinuationToken"`
	Contents              []objectEntry  `xml:"Contents"`
	CommonPrefixes        []commonPrefix `xml:"CommonPrefixes"`
}

type objectEntry struct {
	Key          string    `xml:"Key"`
	Size         int64     `xml:"Size"`
	LastModified time.Time `xml:"LastModified"`
}

type commonPrefix struct {
	Prefix string `xml:"Prefix"`
}

type deleteRequest struct {
	XMLName xml.Name      `xml:"Delete"`
	Quiet   bool          `xml:"Quiet"`
	Objects []deleteEntry `xml:"Object"`
}

type deleteEntry struct {
	Key string `xml:"Key"`
}

type deleteResult struct {
	XMLName xml.Name         `xml:"DeleteResult"`
	Deleted []deletedEntry   `xml:"Deleted"`
	Errors  []deleteEntryErr `xml:"Error"`
}

type deletedEntry struct {
	Key string `xml:"Key"`
}

type deleteEntryErr struct {
	Key     string `xml:"Key"`
	Code    string `xml:"Code"`
	Message string `xml:"Message"`
}

type createBucketConfiguration struct {
	XMLName            xml.Name `xml:"CreateBucketConfiguration"`
	LocationConstraint string   `xml:"LocationConstraint"`
}

package delete

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

var projectIdFlag = globalflags.ProjectIdFlag
var regionFlag = globalflags.RegionFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &objectstorage.APIClient{}
var testProjectId = uuid.NewString()
var testRegion = "eu01"
var testBucketName = "my-bucket"

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testBucketName,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		regionFlag:    testRegion,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
			Region:    testRegion,
		},
		BucketName: testBucketName,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *objectstorage.ApiDeleteBucketRequest)) objectstorage.ApiDeleteBucketRequest {
	request := testClient.DeleteBucket(testCtx, testProjectId, testRegion, testBucketName)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argValues     []string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "bucket name invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(p)
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			err = cmd.ValidateArgs(tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating args: %v", err)
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd, tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing input: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest objectstorage.ApiDeleteBucketRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestS3API(t *testing.T) {
	ctx := context.Background()
	client := s3.New(s3.Options{
		AppID:                        "stackit",
		BaseEndpoint:                 utils.Ptr("https://object.storage.eu01.onstackit.cloud"),
		ClientLogMode:                5,
		ContinueHeaderThresholdBytes: 0,
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     "GTVECKOU1GSVR393LIN0",
				SecretAccessKey: "mAtH/sP7SWYXSexbzKr/CILpBkWOPypUHIddlDkr",
			}, nil
		}),
		Region: "eu01",
	})
	buckets, err := client.ListBuckets(ctx, nil)
	if err != nil {
		t.Fatalf("cannot list buckets: %v", err)
	}
	for _, bucket := range buckets.Buckets {
		log.Printf("%s", *bucket.Name)
		list, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:            bucket.Name,
			ContinuationToken: utils.Ptr("foobarbaz"),
			MaxKeys:           utils.Ptr[int32](50),
		})
		if err != nil {
			log.Fatalf("cannot list bucket: %v", err)
		}
		i := 0
		for {
			for _, obj := range list.Contents {
				var builder strings.Builder
				builder.WriteString(fmt.Sprintf("%03d: ", i))
				if val := obj.Key; val != nil {
					builder.WriteString(fmt.Sprintf("key=%s ", *val))
				}
				if val := obj.ETag; val != nil {
					builder.WriteString(fmt.Sprintf("etag=%s ", *val))
				}
				if val := obj.Size; val != nil {
					builder.WriteString(fmt.Sprintf("size=%d ", *val))
				}
				if val := obj.Owner; val != nil && val.DisplayName != nil {
					builder.WriteString(fmt.Sprintf("size=%d ", val.DisplayName))
				}
				if val := obj.LastModified; val != nil {
					builder.WriteString(fmt.Sprintf("last modified=%s ", val))
				}
				t.Log(builder.String())
				i++
			}

			if !*list.IsTruncated {
				break
			}
			list, err = client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
				Bucket:            bucket.Name,
				ContinuationToken: list.NextContinuationToken,
				MaxKeys:           utils.Ptr[int32](100),
			})
			if err != nil {
				log.Fatalf("cannot continue paging: %v", err)
			}
		}
	}
}

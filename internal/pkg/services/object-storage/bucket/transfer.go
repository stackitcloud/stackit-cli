package bucket

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const maxDeleteBatchSize = 1000

type BucketInfo struct {
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creation_date"`
}

type ObjectInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
}

func ListBuckets(ctx context.Context, client *Client) ([]BucketInfo, error) {
	return client.ListBuckets(ctx)
}

func ListObjects(ctx context.Context, client *Client, ref Reference) ([]ObjectInfo, error) {
	if ref.IsBucket() {
		return nil, fmt.Errorf("cannot list objects of a bare bucket reference %q without a key", ref.String())
	}
	objects := make([]ObjectInfo, 0)
	token := ""
	for {
		page, err := client.ListObjectsV2(ctx, ref.Bucket, ref.Key, "", token)
		if err != nil {
			return nil, fmt.Errorf("list objects in %q: %w", ref.String(), err)
		}
		objects = append(objects, page.Objects...)
		if !page.IsTruncated {
			break
		}
		token = page.NextContinuationToken
	}
	return objects, nil
}

func ListPrefixes(ctx context.Context, client *Client, ref Reference) (prefixes []string, objects []ObjectInfo, err error) {
	page, err := client.ListObjectsV2(ctx, ref.Bucket, ref.Key, "/", "")
	if err != nil {
		return nil, nil, fmt.Errorf("list objects in %q: %w", ref.String(), err)
	}
	prefixes = append(prefixes, page.CommonPrefixes...)
	objects = append(objects, page.Objects...)
	return prefixes, objects, nil
}

func HeadObjectExists(ctx context.Context, client *Client, ref Reference) (bool, error) {
	if ref.IsBucket() {
		return false, fmt.Errorf("cannot head a bare bucket reference %q without a key", ref.String())
	}
	err := client.HeadObject(ctx, ref.Bucket, ref.Key)
	if IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check object %q: %w", ref.String(), err)
	}
	return true, nil
}

func DownloadFile(ctx context.Context, client *Client, ref Reference, localPath string) error {
	exists, err := HeadObjectExists(ctx, client, ref)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("object %q does not exist: if it is a prefix, use --recursive", ref.String())
	}
	return downloadObject(ctx, client, ref, localPath)
}

func downloadObject(ctx context.Context, client *Client, ref Reference, localPath string) error {
	body, err := client.GetObject(ctx, ref.Bucket, ref.Key)
	if err != nil {
		return fmt.Errorf("download object %q: %w", ref.String(), err)
	}
	defer func() {
		if cerr := body.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "close object body: %v\n", cerr)
		}
	}()

	if err := os.MkdirAll(filepath.Dir(localPath), 0o750); err != nil {
		return fmt.Errorf("create local directory for %q: %w", localPath, err)
	}
	// Write to a temp file and rename, so a failed copy never leaves a
	// partial file at the destination
	file, err := os.CreateTemp(filepath.Dir(localPath), filepath.Base(localPath)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create local file %q: %w", localPath, err)
	}
	tmpPath := file.Name()
	defer func() {
		// No-op after a successful rename; removes the temp file on failure
		_ = os.Remove(tmpPath)
	}()
	if _, err := io.Copy(file, body); err != nil {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "close local file: %v\n", cerr)
		}
		return fmt.Errorf("write local file %q: %w", localPath, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close local file %q: %w", localPath, err)
	}
	if err := os.Rename(tmpPath, localPath); err != nil {
		return fmt.Errorf("write local file %q: %w", localPath, err)
	}
	return nil
}

func DownloadPrefix(ctx context.Context, client *Client, ref Reference, localDir string) (int, error) {
	objects, err := ListObjects(ctx, client, ref)
	if err != nil {
		return 0, err
	}
	for _, obj := range objects {
		relPath := strings.TrimPrefix(obj.Key, ref.Key)
		relPath = strings.TrimPrefix(relPath, "/")
		if relPath == "" {
			continue
		}
		localPath := filepath.Join(localDir, relPath)
		if err := downloadObject(ctx, client, Reference{Bucket: ref.Bucket, Key: obj.Key}, localPath); err != nil {
			return 0, err
		}
	}
	return len(objects), nil
}

func UploadFile(ctx context.Context, client *Client, localPath string, ref Reference) error {
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open local file %q: %w", localPath, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "close local file: %v\n", cerr)
		}
	}()
	if err := client.PutObject(ctx, ref.Bucket, ref.Key, file); err != nil {
		return fmt.Errorf("upload %q to %q: %w", localPath, ref.String(), err)
	}
	return nil
}

func UploadDirectory(ctx context.Context, client *Client, localDir string, ref Reference) (int, error) {
	count := 0
	prefix := strings.TrimSuffix(ref.Key, "/")
	err := filepath.WalkDir(localDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk %q: %w", path, walkErr)
		}
		if d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return fmt.Errorf("resolve relative path of %q: %w", path, err)
		}
		key := filepath.ToSlash(relPath)
		if prefix != "" {
			key = prefix + "/" + key
		}
		if err := UploadFile(ctx, client, path, Reference{Bucket: ref.Bucket, Key: key}); err != nil {
			return err
		}
		count++
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func DeleteObjects(ctx context.Context, client *Client, refs []Reference) (int, error) {
	if len(refs) == 0 {
		return 0, nil
	}
	deleted := 0
	for start := 0; start < len(refs); start += maxDeleteBatchSize {
		end := min(start+maxDeleteBatchSize, len(refs))
		batch := refs[start:end]
		keys := make([]string, 0, len(batch))
		for _, ref := range batch {
			keys = append(keys, ref.Key)
		}
		n, err := client.DeleteObjects(ctx, batch[0].Bucket, keys)
		if err != nil {
			return deleted, fmt.Errorf("delete %d objects from %q: %w", len(batch), batch[0].Bucket, err)
		}
		deleted += n
	}
	return deleted, nil
}

func ResolveObjectRefs(ctx context.Context, client *Client, ref Reference, recursive bool) ([]Reference, error) {
	if !recursive {
		exists, err := HeadObjectExists(ctx, client, ref)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("object %q does not exist: if it is a prefix, use --recursive", ref.String())
		}
		return []Reference{ref}, nil
	}
	objects, err := ListObjects(ctx, client, ref)
	if err != nil {
		return nil, err
	}
	if len(objects) == 0 {
		return nil, fmt.Errorf("no objects found under %q", ref.String())
	}
	refs := make([]Reference, 0, len(objects))
	for _, obj := range objects {
		refs = append(refs, Reference{Bucket: ref.Bucket, Key: obj.Key})
	}
	return refs, nil
}

func PresignGetObject(_ context.Context, client *Client, ref Reference, expires time.Duration) (string, error) {
	url, err := client.PresignGetObject(ref.Bucket, ref.Key, expires)
	if err != nil {
		return "", fmt.Errorf("presign URL for %q: %w", ref.String(), err)
	}
	return url, nil
}

func CopyObject(ctx context.Context, client *Client, src, dst Reference) error {
	if err := client.CopyObject(ctx, src, dst); err != nil {
		return fmt.Errorf("copy %q to %q: %w", src.String(), dst.String(), err)
	}
	return nil
}

func localDestinationForObject(dst, key string) string {
	base := filepath.Base(key)
	if dst == "" {
		return base
	}
	if strings.HasSuffix(dst, "/") || strings.HasSuffix(dst, string(filepath.Separator)) {
		return filepath.Join(dst, base)
	}
	if stat, err := os.Stat(dst); err == nil && stat.IsDir() {
		return filepath.Join(dst, base)
	}
	return dst
}

func Download(ctx context.Context, client *Client, ref Reference, localPath string, recursive bool) (int, error) {
	if ref.IsBucket() {
		return 0, fmt.Errorf("cannot download a whole bucket %q: specify an object or prefix", ref.String())
	}
	if recursive {
		// Normalize to a path prefix so that "dir" does not also match
		// "dir-other/..." and map to a bogus local path
		ref.Key = strings.TrimSuffix(ref.Key, "/") + "/"
		return DownloadPrefix(ctx, client, ref, localPath)
	}
	local := localDestinationForObject(localPath, ref.Key)
	err := DownloadFile(ctx, client, ref, local)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func Upload(ctx context.Context, client *Client, localPath string, ref Reference, recursive bool) (int, error) {
	stat, err := os.Stat(localPath)
	if err != nil {
		return 0, fmt.Errorf("stat local path %q: %w", localPath, err)
	}
	if stat.IsDir() {
		if !recursive {
			return 0, fmt.Errorf("local path %q is a directory: use --recursive", localPath)
		}
		return UploadDirectory(ctx, client, localPath, ref)
	}
	key := ref.Key
	if key == "" || strings.HasSuffix(key, "/") {
		key += filepath.Base(localPath)
	}
	err = UploadFile(ctx, client, localPath, Reference{Bucket: ref.Bucket, Key: key})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

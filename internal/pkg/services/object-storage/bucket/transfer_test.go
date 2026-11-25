package bucket

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type fakeStorage struct {
	mu          sync.Mutex
	objects     map[string][]byte
	pageSize    int
	listCalls   int
	deleteCalls [][]string
	copySources map[string]string
}

func newFakeStorage(pageSize int) *fakeStorage {
	return &fakeStorage{
		objects:     map[string][]byte{},
		pageSize:    pageSize,
		copySources: map[string]string{},
	}
}

func (f *fakeStorage) put(bucket, key string, body []byte) {
	f.objects[bucket+"/"+key] = body
}

func (f *fakeStorage) handler(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 2)
	bucket := parts[0]
	key := ""
	if len(parts) == 2 {
		key = parts[1]
	}

	switch {
	case r.Method == http.MethodGet && key == "" && r.URL.Query().Get("list-type") == "2":
		f.listCalls++
		f.writeList(w, r, bucket)
	case r.Method == http.MethodHead:
		if _, ok := f.objects[bucket+"/"+key]; ok {
			w.WriteHeader(http.StatusOK)
		} else {
			writeNoSuchKey(w)
		}
	case r.Method == http.MethodGet:
		if body, ok := f.objects[bucket+"/"+key]; ok {
			_, _ = w.Write(body)
		} else {
			writeNoSuchKey(w)
		}
	case r.Method == http.MethodPut:
		if src := r.Header.Get("X-Amz-Copy-Source"); src != "" {
			f.copySources[bucket+"/"+key] = src
			w.WriteHeader(http.StatusOK)
		} else {
			body, _ := io.ReadAll(r.Body)
			f.objects[bucket+"/"+key] = body
			w.WriteHeader(http.StatusOK)
		}
	case r.Method == http.MethodPost && r.URL.Query().Has("delete"):
		body, _ := io.ReadAll(r.Body)
		var req deleteRequest
		if err := xml.Unmarshal(body, &req); err != nil {
			http.Error(w, "bad delete request: "+err.Error(), http.StatusBadRequest)
			return
		}
		keys := make([]string, 0, len(req.Objects))
		for _, obj := range req.Objects {
			keys = append(keys, obj.Key)
			delete(f.objects, bucket+"/"+obj.Key)
		}
		f.deleteCalls = append(f.deleteCalls, keys)
		_, _ = w.Write([]byte(`<DeleteResult></DeleteResult>`))
	default:
		http.Error(w, "unexpected request: "+r.Method+" "+r.URL.Path, http.StatusNotImplemented)
	}
}

func (f *fakeStorage) writeList(w http.ResponseWriter, r *http.Request, bucket string) {
	query := r.URL.Query()
	prefix := query.Get("prefix")
	delimiter := query.Get("delimiter")

	keys := make([]string, 0)
	for key := range f.objects {
		if !strings.HasPrefix(key, bucket+"/") {
			continue
		}
		rel := strings.TrimPrefix(key, bucket+"/")
		if strings.HasPrefix(rel, prefix) {
			keys = append(keys, rel)
		}
	}
	sort.Strings(keys)

	type entry struct {
		full     string
		isPrefix bool
	}
	var entries []entry
	seen := map[string]bool{}
	for _, key := range keys {
		rel := strings.TrimPrefix(key, prefix)
		if delimiter != "" {
			if i := strings.Index(rel, delimiter); i >= 0 {
				common := prefix + rel[:i+len(delimiter)]
				if !seen[common] {
					seen[common] = true
					entries = append(entries, entry{full: common, isPrefix: true})
				}
				continue
			}
		}
		entries = append(entries, entry{full: key})
	}

	offset, _ := strconv.Atoi(query.Get("continuation-token"))
	end := len(entries)
	if f.pageSize > 0 && offset+f.pageSize < len(entries) {
		end = offset + f.pageSize
	}
	page := entries[offset:end]

	var b strings.Builder
	b.WriteString(`<ListBucketResult xmlns="http://bucket.amazonaws.com/doc/2006-03-01/">`)
	if end < len(entries) {
		b.WriteString(`<IsTruncated>true</IsTruncated><NextContinuationToken>`)
		b.WriteString(strconv.Itoa(end))
		b.WriteString(`</NextContinuationToken>`)
	} else {
		b.WriteString(`<IsTruncated>false</IsTruncated>`)
	}
	for _, e := range page {
		if e.isPrefix {
			b.WriteString(`<CommonPrefixes><Prefix>`)
			b.WriteString(xmlEscape(e.full))
			b.WriteString(`</Prefix></CommonPrefixes>`)
			continue
		}
		b.WriteString(`<Contents><Key>`)
		b.WriteString(xmlEscape(e.full))
		b.WriteString(`</Key><Size>`)
		b.WriteString(strconv.Itoa(len(f.objects[bucket+"/"+e.full])))
		b.WriteString(`</Size><LastModified>2024-05-01T12:00:00Z</LastModified></Contents>`)
	}
	b.WriteString(`</ListBucketResult>`)
	_, _ = w.Write([]byte(b.String()))
}

func writeNoSuchKey(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(`<Error><Code>NoSuchKey</Code><Message>The specified key does not exist.</Message></Error>`))
}

func xmlEscape(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

func newFakeClient(t *testing.T, fake *fakeStorage) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(fake.handler))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{
		Region:          "eu01",
		Endpoint:        server.URL,
		AccessKeyId:     "AKIDEXAMPLE",
		SecretAccessKey: "secret",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

func TestListObjectsPaginates(t *testing.T) {
	fake := newFakeStorage(2)
	for i := 0; i < 5; i++ {
		fake.put("my-bucket", "p/file-"+strconv.Itoa(i)+".txt", []byte("x"))
	}
	client := newFakeClient(t, fake)

	objects, err := ListObjects(context.Background(), client, Reference{Bucket: "my-bucket", Key: "p/"})
	if err != nil {
		t.Fatalf("ListObjects: %v", err)
	}
	if len(objects) != 5 {
		t.Fatalf("ListObjects returned %d objects, want 5", len(objects))
	}
	if objects[0].Key != "p/file-0.txt" || objects[4].Key != "p/file-4.txt" {
		t.Fatalf("unexpected object order: %+v", objects)
	}
	if fake.listCalls != 3 {
		t.Fatalf("list calls = %d, want 3", fake.listCalls)
	}
}

func TestListObjectsBareBucket(t *testing.T) {
	fake := newFakeStorage(0)
	client := newFakeClient(t, fake)

	_, err := ListObjects(context.Background(), client, Reference{Bucket: "my-bucket"})
	if err == nil {
		t.Fatal("expected an error for a bare bucket reference")
	}
}

func TestListPrefixes(t *testing.T) {
	fake := newFakeStorage(0)
	fake.put("my-bucket", "a/1.txt", []byte("x"))
	fake.put("my-bucket", "a/2.txt", []byte("x"))
	fake.put("my-bucket", "b/3.txt", []byte("x"))
	fake.put("my-bucket", "top.txt", []byte("x"))
	client := newFakeClient(t, fake)

	prefixes, objects, err := ListPrefixes(context.Background(), client, Reference{Bucket: "my-bucket"})
	if err != nil {
		t.Fatalf("ListPrefixes: %v", err)
	}
	wantPrefixes := []string{"a/", "b/"}
	if len(prefixes) != len(wantPrefixes) || prefixes[0] != wantPrefixes[0] || prefixes[1] != wantPrefixes[1] {
		t.Fatalf("prefixes = %v, want %v", prefixes, wantPrefixes)
	}
	if len(objects) != 1 || objects[0].Key != "top.txt" {
		t.Fatalf("objects = %+v, want [top.txt]", objects)
	}
}

func TestHeadObjectExists(t *testing.T) {
	fake := newFakeStorage(0)
	fake.put("my-bucket", "there.txt", []byte("x"))
	client := newFakeClient(t, fake)

	tests := []struct {
		description string
		ref         Reference
		expected    bool
		isValid     bool
	}{
		{description: "exists", ref: Reference{Bucket: "my-bucket", Key: "there.txt"}, expected: true, isValid: true},
		{description: "missing", ref: Reference{Bucket: "my-bucket", Key: "nowhere.txt"}, expected: false, isValid: true},
		{description: "bare bucket", ref: Reference{Bucket: "my-bucket"}, isValid: false},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			exists, err := HeadObjectExists(context.Background(), client, tt.ref)
			if (err != nil) != !tt.isValid {
				t.Fatalf("HeadObjectExists error = %v, wantErr %v", err, !tt.isValid)
			}
			if tt.isValid && exists != tt.expected {
				t.Fatalf("HeadObjectExists = %v, want %v", exists, tt.expected)
			}
		})
	}
}

func TestDownloadFile(t *testing.T) {
	fake := newFakeStorage(0)
	fake.put("my-bucket", "dir/file.txt", []byte("hello world"))
	client := newFakeClient(t, fake)
	localPath := filepath.Join(t.TempDir(), "out", "sub", "file.txt")

	err := DownloadFile(context.Background(), client, Reference{Bucket: "my-bucket", Key: "dir/file.txt"}, localPath)
	if err != nil {
		t.Fatalf("DownloadFile: %v", err)
	}
	content, err := os.ReadFile(localPath)
	if err != nil {
		t.Fatalf("read downloaded file: %v", err)
	}
	if string(content) != "hello world" {
		t.Fatalf("downloaded content = %q, want %q", content, "hello world")
	}
}

func TestDownloadFileMissingObject(t *testing.T) {
	fake := newFakeStorage(0)
	client := newFakeClient(t, fake)

	err := DownloadFile(context.Background(), client, Reference{Bucket: "my-bucket", Key: "nope.txt"}, filepath.Join(t.TempDir(), "nope.txt"))
	if err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("error = %v, want it to mention the missing object", err)
	}
}

func TestDownloadPrefix(t *testing.T) {
	fake := newFakeStorage(0)
	fake.put("my-bucket", "p/a.txt", []byte("a"))
	fake.put("my-bucket", "p/nested/b.txt", []byte("b"))
	fake.put("my-bucket", "p/c.txt", []byte("c"))
	fake.put("my-bucket", "other/d.txt", []byte("d"))
	client := newFakeClient(t, fake)
	localDir := t.TempDir()

	count, err := DownloadPrefix(context.Background(), client, Reference{Bucket: "my-bucket", Key: "p/"}, localDir)
	if err != nil {
		t.Fatalf("DownloadPrefix: %v", err)
	}
	if count != 3 {
		t.Fatalf("count = %d, want 3", count)
	}
	for _, rel := range []string{"a.txt", "nested/b.txt", "c.txt"} {
		if _, err := os.Stat(filepath.Join(localDir, rel)); err != nil {
			t.Fatalf("expected downloaded file %q: %v", rel, err)
		}
	}
}

func TestDownloadPrefixExcludesSiblingPrefixes(t *testing.T) {
	fake := newFakeStorage(0)
	fake.put("my-bucket", "dir/a.txt", []byte("a"))
	fake.put("my-bucket", "dir-other/x.txt", []byte("x"))
	client := newFakeClient(t, fake)
	localDir := t.TempDir()

	// A prefix without a trailing slash must behave like a path prefix
	count, err := Download(context.Background(), client, Reference{Bucket: "my-bucket", Key: "dir"}, localDir, true)
	if err != nil {
		t.Fatalf("Download: %v", err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}
	if _, err := os.Stat(filepath.Join(localDir, "a.txt")); err != nil {
		t.Fatalf("expected downloaded file a.txt: %v", err)
	}
	if _, err := os.Stat(filepath.Join(localDir, "-other")); !os.IsNotExist(err) {
		t.Fatalf("did not expect a local path for the sibling prefix dir-other")
	}
}

func TestUploadFile(t *testing.T) {
	fake := newFakeStorage(0)
	client := newFakeClient(t, fake)
	localPath := filepath.Join(t.TempDir(), "source.txt")
	if err := os.WriteFile(localPath, []byte("upload me"), 0o600); err != nil {
		t.Fatalf("write local file: %v", err)
	}

	err := UploadFile(context.Background(), client, localPath, Reference{Bucket: "my-bucket", Key: "dest/file.txt"})
	if err != nil {
		t.Fatalf("UploadFile: %v", err)
	}
	if string(fake.objects["my-bucket/dest/file.txt"]) != "upload me" {
		t.Fatalf("server received %q, want %q", fake.objects["my-bucket/dest/file.txt"], "upload me")
	}
}

func TestUploadDirectory(t *testing.T) {
	fake := newFakeStorage(0)
	client := newFakeClient(t, fake)
	localDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(localDir, "nested"), 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(localDir, "a.txt"), []byte("a"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(localDir, "nested", "b.txt"), []byte("b"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	count, err := UploadDirectory(context.Background(), client, localDir, Reference{Bucket: "my-bucket", Key: "pre/"})
	if err != nil {
		t.Fatalf("UploadDirectory: %v", err)
	}
	if count != 2 {
		t.Fatalf("count = %d, want 2", count)
	}
	if _, ok := fake.objects["my-bucket/pre/a.txt"]; !ok {
		t.Fatalf("missing uploaded key pre/a.txt: %v", fake.objects)
	}
	if _, ok := fake.objects["my-bucket/pre/nested/b.txt"]; !ok {
		t.Fatalf("missing uploaded key pre/nested/b.txt: %v", fake.objects)
	}
}

func TestUploadDirectoryToBareBucket(t *testing.T) {
	fake := newFakeStorage(0)
	client := newFakeClient(t, fake)
	localDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(localDir, "nested"), 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(localDir, "a.txt"), []byte("a"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(localDir, "nested", "b.txt"), []byte("b"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	count, err := UploadDirectory(context.Background(), client, localDir, Reference{Bucket: "my-bucket"})
	if err != nil {
		t.Fatalf("UploadDirectory: %v", err)
	}
	if count != 2 {
		t.Fatalf("count = %d, want 2", count)
	}
	if _, ok := fake.objects["my-bucket/a.txt"]; !ok {
		t.Fatalf("missing uploaded key a.txt: %v", fake.objects)
	}
	if _, ok := fake.objects["my-bucket/nested/b.txt"]; !ok {
		t.Fatalf("missing uploaded key nested/b.txt: %v", fake.objects)
	}
	for key := range fake.objects {
		if strings.HasPrefix(key, "my-bucket//") {
			t.Fatalf("unexpected leading-slash key %q", key)
		}
	}
}

func TestDeleteObjectsBatchesInGroupsOf1000(t *testing.T) {
	fake := newFakeStorage(0)
	client := newFakeClient(t, fake)
	refs := make([]Reference, 0, 2500)
	for i := range 2500 {
		refs = append(refs, Reference{Bucket: "my-bucket", Key: "k" + strconv.Itoa(i)})
	}

	deleted, err := DeleteObjects(context.Background(), client, refs)
	if err != nil {
		t.Fatalf("DeleteObjects: %v", err)
	}
	if deleted != 2500 {
		t.Fatalf("deleted = %d, want 2500", deleted)
	}
	if len(fake.deleteCalls) != 3 {
		t.Fatalf("delete calls = %d, want 3", len(fake.deleteCalls))
	}
	if len(fake.deleteCalls[0]) != 1000 || len(fake.deleteCalls[1]) != 1000 || len(fake.deleteCalls[2]) != 500 {
		t.Fatalf("batch sizes = %d/%d/%d, want 1000/1000/500",
			len(fake.deleteCalls[0]), len(fake.deleteCalls[1]), len(fake.deleteCalls[2]))
	}
}

func TestDeleteObjectsEmpty(t *testing.T) {
	fake := newFakeStorage(0)
	client := newFakeClient(t, fake)

	deleted, err := DeleteObjects(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("DeleteObjects: %v", err)
	}
	if deleted != 0 {
		t.Fatalf("deleted = %d, want 0", deleted)
	}
	if len(fake.deleteCalls) != 0 {
		t.Fatalf("delete calls = %d, want 0", len(fake.deleteCalls))
	}
}

func TestCopyObject(t *testing.T) {
	fake := newFakeStorage(0)
	fake.put("my-bucket", "src.txt", []byte("data"))
	client := newFakeClient(t, fake)

	err := CopyObject(context.Background(), client,
		Reference{Bucket: "my-bucket", Key: "src.txt"},
		Reference{Bucket: "my-bucket", Key: "dst.txt"})
	if err != nil {
		t.Fatalf("CopyObject: %v", err)
	}
	if got := fake.copySources["my-bucket/dst.txt"]; got != "my-bucket/src.txt" {
		t.Fatalf("copy source = %q, want %q", got, "my-bucket/src.txt")
	}
}

func TestPresignGetObject(t *testing.T) {
	fake := newFakeStorage(0)
	client := newFakeClient(t, fake)

	url, err := PresignGetObject(context.Background(), client, Reference{Bucket: "my-bucket", Key: "dir/file.txt"}, 15*time.Minute)
	if err != nil {
		t.Fatalf("PresignGetObject: %v", err)
	}
	if !strings.Contains(url, "my-bucket/dir/file.txt") || !strings.Contains(url, "X-Amz-Signature=") {
		t.Fatalf("presigned URL %q is missing the key or the signature", url)
	}
}

func TestLocalDestinationForObject(t *testing.T) {
	existingDir := t.TempDir()
	tests := []struct {
		description string
		dst         string
		key         string
		expected    string
	}{
		{description: "empty dst uses the base name", dst: "", key: "dir/file.txt", expected: "file.txt"},
		{description: "dst with trailing slash appends the base name", dst: "/tmp/x/", key: "dir/file.txt", expected: "/tmp/x/file.txt"},
		{description: "existing directory appends the base name", dst: existingDir, key: "dir/file.txt", expected: filepath.Join(existingDir, "file.txt")},
		{description: "file path is used as-is", dst: "/tmp/out.txt", key: "dir/file.txt", expected: "/tmp/out.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if got := localDestinationForObject(tt.dst, tt.key); got != tt.expected {
				t.Fatalf("localDestinationForObject = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDownload(t *testing.T) {
	fake := newFakeStorage(0)
	fake.put("my-bucket", "file.txt", []byte("single"))
	fake.put("my-bucket", "p/a.txt", []byte("a"))
	fake.put("my-bucket", "p/b.txt", []byte("b"))
	client := newFakeClient(t, fake)

	t.Run("bare bucket is rejected", func(t *testing.T) {
		_, err := Download(context.Background(), client, Reference{Bucket: "my-bucket"}, t.TempDir(), false)
		if err == nil {
			t.Fatal("expected an error for a bare bucket reference")
		}
	})

	t.Run("single object", func(t *testing.T) {
		localPath := filepath.Join(t.TempDir(), "out.txt")
		count, err := Download(context.Background(), client, Reference{Bucket: "my-bucket", Key: "file.txt"}, localPath, false)
		if err != nil {
			t.Fatalf("Download: %v", err)
		}
		if count != 1 {
			t.Fatalf("count = %d, want 1", count)
		}
		content, err := os.ReadFile(localPath)
		if err != nil {
			t.Fatalf("read downloaded file: %v", err)
		}
		if string(content) != "single" {
			t.Fatalf("content = %q, want %q", content, "single")
		}
	})

	t.Run("recursive prefix", func(t *testing.T) {
		localDir := t.TempDir()
		count, err := Download(context.Background(), client, Reference{Bucket: "my-bucket", Key: "p/"}, localDir, true)
		if err != nil {
			t.Fatalf("Download: %v", err)
		}
		if count != 2 {
			t.Fatalf("count = %d, want 2", count)
		}
	})
}

func TestUpload(t *testing.T) {
	fake := newFakeStorage(0)
	client := newFakeClient(t, fake)
	localFile := filepath.Join(t.TempDir(), "source.txt")
	if err := os.WriteFile(localFile, []byte("data"), 0o600); err != nil {
		t.Fatalf("write local file: %v", err)
	}

	t.Run("missing local path", func(t *testing.T) {
		_, err := Upload(context.Background(), client, filepath.Join(t.TempDir(), "missing.txt"), Reference{Bucket: "my-bucket"}, false)
		if err == nil {
			t.Fatal("expected an error for a missing local path")
		}
	})

	t.Run("directory without recursive flag", func(t *testing.T) {
		_, err := Upload(context.Background(), client, t.TempDir(), Reference{Bucket: "my-bucket"}, false)
		if err == nil || !strings.Contains(err.Error(), "--recursive") {
			t.Fatalf("error = %v, want it to mention --recursive", err)
		}
	})

	t.Run("trailing slash key appends the base name", func(t *testing.T) {
		_, err := Upload(context.Background(), client, localFile, Reference{Bucket: "my-bucket", Key: "pre/"}, false)
		if err != nil {
			t.Fatalf("Upload: %v", err)
		}
		if _, ok := fake.objects["my-bucket/pre/source.txt"]; !ok {
			t.Fatalf("missing uploaded key pre/source.txt: %v", fake.objects)
		}
	})

	t.Run("empty key appends the base name", func(t *testing.T) {
		_, err := Upload(context.Background(), client, localFile, Reference{Bucket: "my-bucket"}, false)
		if err != nil {
			t.Fatalf("Upload: %v", err)
		}
		if _, ok := fake.objects["my-bucket/source.txt"]; !ok {
			t.Fatalf("missing uploaded key source.txt: %v", fake.objects)
		}
	})
}

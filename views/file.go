package views

import (
	"errors"
	"fmt"
	"github.com/hujh/gmvc"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FileView struct {
	root string
}

func NewFileView(root string) *FileView {
	return &FileView{
		root: path.Clean(root),
	}
}

func (v *FileView) Render(c *gmvc.Context, name string, data interface{}) error {
	w := c.ResponseWriter
	r := c.Request

	root, err := filepath.Abs(v.root)
	if err != nil {
		return err
	}

	p := path.Join(root, name)

	if !strings.HasPrefix(p, root) {
		c.Status(http.StatusBadRequest)
		return nil
	}

	f, err := os.Open(p)

	if err != nil {
		c.Status(http.StatusNotFound)
		return nil
	}

	defer f.Close()

	d, err1 := f.Stat()
	if err1 != nil {
		c.Status(http.StatusNotFound)
		return nil
	}

	if d.IsDir() {
		c.Status(http.StatusNotFound)
		return nil
	}

	if v.checkLastModified(w, r, d.ModTime()) {
		return nil
	}

	code := http.StatusOK

	ctype := w.Header().Get("Content-Type")
	if ctype == "" {
		ctype = mime.TypeByExtension(filepath.Ext(name))
		if ctype != "" {
			w.Header().Set("Content-Type", ctype)
		}
	}

	size := d.Size()
	sendSize := size

	var sendContent io.Reader = f
	if size >= 0 {
		ranges, err := v.parseRange(r.Header.Get("Range"), size)
		if err != nil {
			c.ErrorStatus(err, http.StatusRequestedRangeNotSatisfiable)
			return nil
		}

		if v.sumRangesSize(ranges) >= size {
			ranges = nil
		}

		switch {
		case len(ranges) == 1:
			ra := ranges[0]
			if _, err := f.Seek(ra.start, os.SEEK_SET); err != nil {
				c.ErrorStatus(err, http.StatusRequestedRangeNotSatisfiable)
				return nil
			}
			sendSize = ra.length
			code = http.StatusPartialContent
			w.Header().Set("Content-Range", ra.contentRange(size))
		case len(ranges) > 1:
			for _, ra := range ranges {
				if ra.start > size {
					c.ErrorStatus(err, http.StatusRequestedRangeNotSatisfiable)
					return nil
				}
			}
			sendSize = v.rangesMIMESize(ranges, ctype, size)
			code = http.StatusPartialContent

			pr, pw := io.Pipe()
			mw := multipart.NewWriter(pw)
			w.Header().Set("Content-Type", "multipart/byteranges; boundary="+mw.Boundary())
			sendContent = pr
			defer pr.Close()
			go func() {
				for _, ra := range ranges {
					part, err := mw.CreatePart(ra.mimeHeader(ctype, size))
					if err != nil {
						pw.CloseWithError(err)
						return
					}
					if _, err := f.Seek(ra.start, os.SEEK_SET); err != nil {
						pw.CloseWithError(err)
						return
					}
					if _, err := io.CopyN(part, f, ra.length); err != nil {
						pw.CloseWithError(err)
						return
					}
				}
				mw.Close()
				pw.Close()
			}()
		}

		w.Header().Set("Accept-Ranges", "bytes")
		if w.Header().Get("Content-Encoding") == "" {
			w.Header().Set("Content-Length", strconv.FormatInt(sendSize, 10))
		}
	}

	if !c.WroteHeader() {
		w.WriteHeader(code)
	}

	if r.Method != "HEAD" {
		io.CopyN(w, sendContent, sendSize)
	}

	return nil
}

func (v *FileView) checkLastModified(w http.ResponseWriter, r *http.Request, mtime time.Time) bool {
	if mtime.IsZero() {
		return false
	}

	t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since"))

	if err == nil && mtime.Before(t.Add(1*time.Second)) {
		w.WriteHeader(http.StatusNotModified)
		return true
	}

	w.Header().Set("Last-Modified", mtime.UTC().Format(http.TimeFormat))
	return false
}

func (v *FileView) parseRange(s string, size int64) ([]httpRange, error) {
	if s == "" {
		return nil, nil
	}
	const b = "bytes="
	if !strings.HasPrefix(s, b) {
		return nil, errors.New("invalid range")
	}
	var ranges []httpRange
	for _, ra := range strings.Split(s[len(b):], ",") {
		ra = strings.TrimSpace(ra)
		if ra == "" {
			continue
		}
		i := strings.Index(ra, "-")
		if i < 0 {
			return nil, errors.New("invalid range")
		}
		start, end := strings.TrimSpace(ra[:i]), strings.TrimSpace(ra[i+1:])
		var r httpRange
		if start == "" {
			// If no start is specified, end specifies the
			// range start relative to the end of the file.
			i, err := strconv.ParseInt(end, 10, 64)
			if err != nil {
				return nil, errors.New("invalid range")
			}
			if i > size {
				i = size
			}
			r.start = size - i
			r.length = size - r.start
		} else {
			i, err := strconv.ParseInt(start, 10, 64)
			if err != nil || i > size || i < 0 {
				return nil, errors.New("invalid range")
			}
			r.start = i
			if end == "" {
				// If no end is specified, range extends to end of the file.
				r.length = size - r.start
			} else {
				i, err := strconv.ParseInt(end, 10, 64)
				if err != nil || r.start > i {
					return nil, errors.New("invalid range")
				}
				if i >= size {
					i = size - 1
				}
				r.length = i - r.start + 1
			}
		}
		ranges = append(ranges, r)
	}
	return ranges, nil
}

func (v *FileView) sumRangesSize(ranges []httpRange) (size int64) {
	for _, ra := range ranges {
		size += ra.length
	}
	return
}

func (v *FileView) rangesMIMESize(ranges []httpRange, contentType string, contentSize int64) (encSize int64) {
	var w countingWriter
	mw := multipart.NewWriter(&w)
	for _, ra := range ranges {
		mw.CreatePart(ra.mimeHeader(contentType, contentSize))
		encSize += ra.length
	}
	mw.Close()
	encSize += int64(w)
	return
}

type httpRange struct {
	start, length int64
}

func (r httpRange) contentRange(size int64) string {
	return fmt.Sprintf("bytes %d-%d/%d", r.start, r.start+r.length-1, size)
}

func (r httpRange) mimeHeader(contentType string, size int64) textproto.MIMEHeader {
	return textproto.MIMEHeader{
		"Content-Range": {r.contentRange(size)},
		"Content-Type":  {contentType},
	}
}

type countingWriter int64

func (w *countingWriter) Write(p []byte) (n int, err error) {
	*w += countingWriter(len(p))
	return len(p), nil
}

package protocol

import (
	"net/http"
	"strconv"

	tusdS3Handler "github.com/tus/tusd/v2/pkg/handler"
)

type draftVersion string

// These are the different interoperability versions defines in the different
// versions of the resumable uploads draft from the HTTP working group.
// See https://datatracker.ietf.org/doc/draft-ietf-httpbis-resumable-upload/
const (
	interopVersion3 draftVersion = "3" // From draft version -01
	interopVersion4 draftVersion = "4" // From draft version -02
	interopVersion5 draftVersion = "5" // From draft version -03
	interopVersion6 draftVersion = "6" // From draft version -04 and -05
)

func getIETFDraftUploadLength(r *http.Request) (length int64, lengthIsDeferred bool, err error) {
	var lengthFromUploadLength int64
	hasLengthFromUploadLength := false
	var lengthFromContentLength int64
	hasLengthFromContentLength := false

	willCompleteUpload := isIETFDraftUploadComplete(r)
	if getIETFDraftInteropVersion(r) == "" && r.Method == http.MethodPost {
		willCompleteUpload = true
	}
	if willCompleteUpload && r.ContentLength != -1 {
		lengthFromContentLength = r.ContentLength
		hasLengthFromContentLength = true
	}

	uploadLengthStr := r.Header.Get("Upload-Length")
	if uploadLengthStr != "" {
		var err error
		lengthFromUploadLength, err = strconv.ParseInt(uploadLengthStr, 10, 64)
		if err != nil {
			return 0, false, tusdS3Handler.ErrInvalidUploadLength
		}

		hasLengthFromUploadLength = true
	}

	// If both lengths are set, they must match
	if hasLengthFromContentLength && hasLengthFromUploadLength && lengthFromUploadLength != lengthFromContentLength {
		return 0, false, tusdS3Handler.ErrInvalidUploadLength
	}

	// Return whichever length is set
	if hasLengthFromUploadLength {
		return lengthFromUploadLength, false, nil
	}
	if hasLengthFromContentLength {
		return lengthFromContentLength, false, nil
	}

	// No length set, so it's deferred
	return 0, true, nil
}

func isIETFDraftUploadComplete(r *http.Request) bool {
	currentUploadDraftInteropVersion := getIETFDraftInteropVersion(r)
	switch currentUploadDraftInteropVersion {
	case interopVersion4, interopVersion5, interopVersion6:
		return r.Header.Get("Upload-Complete") == "?1"
	case interopVersion3:
		return r.Header.Get("Upload-Incomplete") == "?0"
	default:
		return false
	}
}

func getIETFDraftInteropVersion(r *http.Request) draftVersion {
	version := draftVersion(r.Header.Get("Upload-Draft-Interop-Version"))
	switch version {
	case interopVersion3, interopVersion4, interopVersion5, interopVersion6:
		return version
	default:
		return ""
	}
}

// getIETFDraftUploadLimits returns the Upload-Limit header for a given upload
// according to the set resumable upload draft version from IETF.
func (p Protocol) getIETFDraftUploadLimits(info FileInfo) string {
	limits := "min-size=0"
	if p.cfg.Protocol.MaxSize > 0 {
		limits += ",max-size=" + strconv.FormatInt(p.cfg.Protocol.MaxSize, 10)
	} else if !info.SizeIsDeferred {
		limits += ",max-size=" + strconv.FormatInt(info.Size, 10)
	}

	return limits
}

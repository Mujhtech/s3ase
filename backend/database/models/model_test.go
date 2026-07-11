package models

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadataJSONAndDatabaseRoundTrip(t *testing.T) {
	original := Metadata{"events": []interface{}{"upload.completed"}}
	encoded, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded Metadata
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Equal(t, original, decoded)

	value, err := original.Value()
	require.NoError(t, err)
	var scanned Metadata
	require.NoError(t, scanned.Scan(value))
	require.Equal(t, original, scanned)
}

func TestMetadataScansLegacyBase64JSON(t *testing.T) {
	legacy := []byte(`"` + base64.StdEncoding.EncodeToString([]byte(`{"enabled":true}`)) + `"`)
	var metadata Metadata
	require.NoError(t, metadata.Scan(legacy))
	require.Equal(t, true, metadata["enabled"])
}

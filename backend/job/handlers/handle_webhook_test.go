package handlers

import (
	"testing"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/stretchr/testify/require"
)

func TestSubscribesTo(t *testing.T) {
	metadata := models.Metadata{"events": []interface{}{"upload.completed", "upload.deleted"}}
	require.True(t, subscribesTo(metadata, "upload.completed"))
	require.False(t, subscribesTo(metadata, "upload.failed"))
}

package models

import (
	"testing"

	"snippetbox.jonnevuorela.com/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name   string
		userId int
		want   bool
	}{
		{
			name:   "Valid Id",
			userId: 1,
			want:   true,
		},
		{
			name:   "Zero Id",
			userId: 0,
			want:   false,
		},
		{
			name:   "Non-existent Id",
			userId: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := UserModel{db}

			exists, err := m.Exists(tt.userId)

			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}

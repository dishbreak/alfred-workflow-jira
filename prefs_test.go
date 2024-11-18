package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPrefs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Prefs
		wantErr bool
	}{
		{
			name: "valid prefs",
			input: `
<plist version="1.0">
<dict>
	<key>jira_board_id</key>
	<string>123</string>
	<key>jira_email</key>
	<string>test@example.com</string>
	<key>jira_url</key>
	<string>https://example.com</string>
</dict>
</plist>`,
			want: &Prefs{
				MyBoardId: "123",
				JiraEmail: "test@example.com",
				JiraUrl:   "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid prefs",
			input: `
<plist version="1.0">
<dict>
	<key>jira_board_id</key>
	<string>invalid</string>
</dict>
</plist>`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader([]byte(tt.input))
			got, err := loadPrefs(r)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

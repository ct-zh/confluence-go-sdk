package template

import (
	"strings"
	"testing"
)

func TestBuiltinLoader(t *testing.T) {
	loader := NewBuiltinLoader()

	tests := []struct {
		name    string
		tplName string
		wantErr bool
	}{
		{
			name:    "Load existing swagger template",
			tplName: SwaggerDoc,
			wantErr: false,
		},
		{
			name:    "Load non-existent template",
			tplName: "non_existent.html",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loader.Load(tt.tplName)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuiltinLoader.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("BuiltinLoader.Load() returned empty content")
			}
			if !tt.wantErr && !strings.Contains(got, "{{.Title}}") {
				t.Errorf("BuiltinLoader.Load() content seems invalid")
			}
		})
	}
}

package schema

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	schemaPath         = "test_data/schema"
	templateConfigPath = "test_data/templates/template_config.yaml"
	templatesPath      = "test_data/templates"

	allPath = "test_output/all.yml"

	hogeGoPath    = "test_output/hoge.go"
	hogeProtoPath = "test_output/hoge.proto"
	hogeSQLPath   = "test_output/hoge.sql"
)

func TestApplyTemplatesAddsGenerationPrefix(t *testing.T) {
	tests := []struct {
		name           string
		filePath       string
		expectedPrefix string
	}{
		{
			name:           "given YAML file",
			filePath:       allPath,
			expectedPrefix: "^# Code generated by contrailschema tool .* DO NOT EDIT.",
		},
		{
			name:           "given Go file",
			filePath:       hogeGoPath,
			expectedPrefix: "^// Code generated by contrailschema tool .* DO NOT EDIT.",
		},
		{
			name:           "given Proto file",
			filePath:       hogeProtoPath,
			expectedPrefix: "^// Code generated by contrailschema tool .* DO NOT EDIT.",
		},
		{
			name:           "given SQL file",
			filePath:       hogeSQLPath,
			expectedPrefix: "^-- Code generated by contrailschema tool .* DO NOT EDIT.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ApplyTemplates(makeAPI(t), filepath.Dir(templatesPath), loadTemplates(t), &TemplateOption{})

			assert.Nil(t, err)
			assert.Regexp(t, tt.expectedPrefix, loadString(t, tt.filePath))
		})
	}
}

func makeAPI(t *testing.T) *API {
	api, err := MakeAPI([]string{schemaPath}, "")
	assert.Nil(t, err)

	return api
}

func loadTemplates(t *testing.T) []*TemplateConfig {
	c, err := LoadTemplates(templateConfigPath)
	assert.Nil(t, err)

	return c
}

func loadString(t *testing.T, path string) string {
	data, err := ioutil.ReadFile(path)
	assert.Nil(t, err)

	return string(data)
}

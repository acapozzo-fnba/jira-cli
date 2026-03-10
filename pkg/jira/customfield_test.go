package jira

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomFieldRichTextSerialization(t *testing.T) {
	rt := customFieldTypeRichText("h2. Section Title\n* Item one\n* Item two\n")
	b, err := json.Marshal(rt)
	assert.NoError(t, err)
	assert.JSONEq(t, `"h2. Section Title\n* Item one\n* Item two\n"`, string(b))
}

func TestCustomFieldRichTextSetSerialization(t *testing.T) {
	rts := customFieldTypeRichTextSet{
		Set: customFieldTypeRichText("*Bold text*"),
	}
	b, err := json.Marshal(rts)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"set":"*Bold text*"}`, string(b))
}

func TestCustomFieldRichTextSetArraySerialization(t *testing.T) {
	// Edit operations wrap in an array of set operations
	arr := []customFieldTypeRichTextSet{
		{Set: customFieldTypeRichText("h2. Title\n")},
	}
	b, err := json.Marshal(arr)
	assert.NoError(t, err)
	assert.JSONEq(t, `[{"set":"h2. Title\n"}]`, string(b))
}

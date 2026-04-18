package folks

import (
	"encoding/json"
	"ethno/internal/models"
	"fmt"
	"unicode/utf8"
)

func ParseFolkContent(raw json.RawMessage) (*models.FolkContent, error) {
    if raw == nil {
        return nil, nil
    }

    var content models.FolkContent
    if err := json.Unmarshal(raw, &content); err == nil && content.Text != "" {
        return &content, nil
    }

    var text string
    if err := json.Unmarshal(raw, &text); err == nil {
        return &models.FolkContent{Text: text}, nil
    }

    if utf8.Valid(raw) {
        return &models.FolkContent{Text: string(raw)}, nil
    }

    return nil, fmt.Errorf("failed to parse folk summary: invalid encoding")
}
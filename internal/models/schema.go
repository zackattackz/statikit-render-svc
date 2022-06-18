package models

import "html/template"

type Schema struct {
	Data    map[string]any           // Variable names->raw data to be substituted, comes directly from schema
	FileSub map[string]template.HTML // Variable names->html to be substituted, comes from a file
}

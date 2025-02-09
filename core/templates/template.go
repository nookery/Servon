package templates

import "embed"

//go:embed systemd_service.tmpl
var TemplateFS embed.FS

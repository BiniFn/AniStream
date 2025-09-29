package template

import _ "embed"

//go:embed email/forgot-password.html
var ForgetPasswordEmailTemplate string

//go:embed admin/index.html
var AdminPanelTemplate string

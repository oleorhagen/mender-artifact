// Copyright 2020 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

var ttyWidth = 80

func init() {
	// Get the width of the terminal
}

// SubcommandHelpTemplate is the text template for the subcommand help topic.
// cli.go uses text/template to render templates.
// This template is used for sub-commands with one or more required flags present.
// TODO - 80 is the hard-coded line-length. Make this configureable.
const CustomSubcommandHelpTemplate = `.TH FORMATWITHGROFF
.ll 80
.SH NAME:
{{.HelpName}} - {{if .Description}}{{.Description}}{{else}}{{.Usage}}{{end}}
.SH USAGE:
 {{if .UsageText}}
{{.UsageText}}
{{else}}
{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
{{end}}
{{if .VisibleFlags}}{{$reqs:=false}}{{range .VisibleFlags}}{{if .IsRequired}}{{$reqs = true}}{{end}}{{end}}
{{- if $reqs}}
.SH REQUIRED ARGUMENTS:
{{range .VisibleFlags}}{{if .IsRequired}}
.TP
\fB{{.Name}}\fP
{{.Usage}}{{end}}{{end}}{{end}}
.SH OPTIONS:
{{range .VisibleFlags}}{{if not .IsRequired}}
.TP
.sp -1
\fB--{{.Name}}\FP
{{.Usage}}{{end}}{{end}}{{end}}

foobar
`


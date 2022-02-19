//go:build wasm
// +build wasm

package app

import (
	"encoding/xml"
	"html/template"
	"strings"
	"syscall/js"
)

func WASMRender(data string) string {
	protocol := Protocol{}
	err := xml.Unmarshal([]byte(data), &protocol)
	if err != nil {
		return (err.Error())
	}
	tmpl, err := template.New("").Funcs(funcMap).Parse(`
	{{ define "description" }}
	{{ if . }}
		<p>
			{{ . }}
		</p>
	{{ end }}
	{{ end }}
	
	{{ define "()" }}
	{{ if . }}
		({{ . }})
	{{ end }}
	{{ end }}
	
	{{ define "()-light" }}
	{{ if . }}
		<small class="text-muted">({{ . }})</small>
	{{ end }}
	{{ end }}
	
	{{ define "arguments-table" }}
	{{ with . }}
	<h5 class="font-bold">Arguments</h5>
	<table>
		<thead>
			<tr>
				<th>
					Name
				</th>
				<th>
					Type
				</th>
				<th>
					Description
				</th>
			</tr>
		</thead>
		<tbody>
			{{ range $arg := . }}
			<tr>
				<td>{{ $arg.Name }}</td>
				<td>{{ $arg.Type }}{{- with $arg.Interface -}}[{{ . }}]{{ end }}</td>
				<td>{{ $arg.Summary }}</td>
			</tr>
			{{ end }}
		</tbody>
	</table>
	{{ end }}
	{{ end }}
	
	{{ define "enum-table" }}
	{{ with . }}
	<h5 class="font-bold">Entries</h5>
	<table>
		<thead>
			<tr>
				<th>
					Name
				</th>
				<th>
					Value
				</th>
				<th>
					Description
				</th>
			</tr>
		</thead>
		<tbody>
			{{ range $enum := . }}
			<tr>
				<td>{{ $enum.Name }}</td>
				<td>{{ $enum.Value }}</td>
				<td>{{ $enum.Summary }}</td>
			</tr>
			{{ end }}
		</tbody>
	</table>
	{{ end }}
	{{ end }}
	
	{{ define "arguments" }}
	{{- if . -}}
		{{- $alen := len . -}}
		{{- range $idx, $arg := . -}}
			{{- $arg.Name }} <span class="symbol-structure">{{$arg.Type}}</span>
			{{- with $arg.Interface -}}
				[<span class="symbol-class">{{ . }}</span>]
			{{- end -}}
			{{- if not (isLast $idx $alen) }}, {{ end -}}
		{{ end -}}
	{{- end -}}
	{{ end }}
	
	{{ define "name" -}}
	{{- $req := . -}}
	{{- if hasField $req "Type" -}}
	{{- if eq $req.Type "destructor" -}}
		~
	{{- end -}}
	{{- end -}}
	<span class="symbol-method">{{- $req.Name -}}</span>
	{{- end }}
	
	{{ define "message-name" }}
	{{ $msg := . }}
	{{ $msg.Name }}
	{{ if hasField $msg "Type" }}
	{{ template "()-light" $msg.Type }}
	{{ end }}
	{{ end }}
	
	{{ define "message" }}
	{{ $msg := . }}
	<h4 id="{{ $msg.Anchor }}">{{ template "message-name" $msg }} <small class="text-muted">since version {{ $msg.Since }}</small></h4>
	<pre><code>{{ template "name" $msg }}({{ template "arguments" $msg.Arguments }})</code></pre>
	<div>
		<div>
			{{ template "description" $msg.Description.Body }}
			{{ template "arguments-table" $msg.Arguments }}
		</div>
	</div>
	<br>
	{{ end }}
	
	{{ define "requests" }}
	{{ if . }}
		<h3>Requests</h3>
		<div>
		{{ range $req := . }}
		{{ template "message" $req }}
		{{ end }}
		</div>
	{{ end }}
	{{ end }}
	
	{{ define "events" }}
	{{ if . }}
		<h3>Events</h3>
		<div>
		{{ range $req := . }}
		{{ template "message" $req }}
		{{ end }}
		</div>
	{{ end }}
	{{ end }}
	
	{{ define "enum-entries" }}
	{{ $enum := . }}
	{{ with $enum.Entries }}
		<pre><code><span class="symbol-keyword">enum</span> <span class="symbol-enum">{{ $enum.Name }}</span> {<br/>{{ range $entry := . }}    <span class="symbol-enum-member">{{ $entry.Name }}</span> = {{ $entry.Value -}},<br>{{ end }}}</code></pre>
	{{ end }}
	{{ end }}
	
	{{ define "enums" }}
	{{ if . }}
		<h3>Enums</h3>
		<div>
		{{ range $enum := . }}
			<div id="{{ $enum.Anchor }}">
				<h4>{{ if $enum.Bitfield }}Flagset{{ end }} {{ $enum.Name }} <small class="text-muted">since version {{ $enum.Since }}</small></h4>
				{{ template "description" $enum.Description.Body }}
				{{ template "enum-entries" $enum }}
			</div>
			{{ template "enum-table" $enum.Entries }}
		{{ end }}
		</div>
	{{ end }}
	{{ end }}
	{{ define "" }}
	<!DOCTYPE html>
	<html prefix="og: http://ogp.me/ns#">
		<head>
			<meta charset="UTF-8">
			<title>ReadWay</title>
			<meta property="og:title" content="ReadWay" />
			<meta property="og:description" content="Why read XML when you can use a website with well-designed typography?" />
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<script src="wasm_exec.js"></script>
			<link rel="stylesheet" href="main.css"/>
			<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@vscode/codicons@0.0.28/dist/codicon.css" />
		</head>
		<body class="bg-stone-100 dark:bg-stone-900 text-black dark:text-white">
			<main>
				<header class="bg-stone-200 dark:bg-stone-800 p-2">
					<nav class="flex flex-row space-x-2">
						<a href="index.html">ReadWay Protocol Browser</a>
						<span> / </span>
						<a href="#">{{ .Name }}</a>
					</nav>
				</header>
				
				<div class="flex flex-row w-full">
					<aside class="prose dark:prose-invert prose-sm w-full max-w-xs p-4">
						{{ range $iface := .Interfaces }}
							<h4><span class="codicon codicon-symbol-class symbol-class"></span><a href="#{{ $iface.Anchor }}">{{ $iface.Name }}</a><br></h4>
							{{ if $iface.Requests }}
							Requests:
							<ul>
							{{ range $req := $iface.Requests }}
								<li><span class="codicon codicon-symbol-method symbol-method"></span><a href="#{{ $req.Anchor }}">{{ $req.Name }}</a></li>
							{{ end }}
							</ul>
							{{ end }}
							{{ if $iface.Events }}
							Events:
							<ul>
							{{ range $ev := $iface.Events }}
								<li><span class="codicon codicon-symbol-event symbol-event"></span><a href="#{{ $ev.Anchor }}">{{ $ev.Name }}</a></li>
							{{ end }}
							</ul>
							{{ end }}
							{{ if $iface.Enums }}
							Enums:
							<ul>
							{{ range $en := $iface.Enums }}
								<li><span class="codicon codicon-symbol-enum symbol-enum"></span><a href="#{{ $en.Anchor }}">{{ $en.Name }}</a></li>
							{{ end }}
							</ul>
							{{ end }}
						{{ end }}
					</aside>
					<article class="prose dark:prose-invert w-full pt-4">
						<h1>{{ .Name }} <small class="text-muted">protocol</small></h1>
				
						{{ template "description" .Description.Body }}
				
						{{ range $iface := .Interfaces }}
							<h2 id="{{ $iface.Anchor }}">{{ $iface.Name }} <small class="text-muted">interface version {{ $iface.Version }}</small></h2>
				
							{{ template "description" $iface.Description.Body }}
				
							{{ template "requests" $iface.Requests }}
							
							{{ template "events" $iface.Events }}
				
							{{ template "enums" $iface.Enums }}
				
							<hr>
						{{ end }}
				
						{{ template "description" .Copyright }}
					</article>
				</div>
			</main>
		</body>
	</html>
	{{ end }}`)
	if err != nil {
		return (err.Error())
	}
	sb := strings.Builder{}
	tmpl.Execute(&sb, protocol)
	return sb.String()
}

func init() {
	js.Global().Set("render", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return WASMRender(args[0].String())
	}))
}

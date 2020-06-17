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
	tmpl, err := template.New("").Parse(`
	{{ define "description" }}
	{{ if . }}
		<p class="lead">
			{{ . }}
		</p>
	{{ end }}
	{{ end }}
	
	{{ define "description-light" }}
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
	
	{{ define "arguments" }}
	{{ if . }}
		<div class="list-group list-group-flush">
			{{ range $arg := . }}
				<span class="list-group-item">
					<h6 class="mb-1">{{ $arg.Name }} {{ template "()-light" $arg.Type }} {{ template "()-light" $arg.Interface }} {{ template "()-light" $arg.Enum }}</h6>
					{{ template "description-light" $arg.Summary }}
					{{ template "description-light" $arg.Description.Body }}
				</span>
			{{ end }}
		</div>
	{{ end }}
	{{ end }}
	
	{{ define "requests" }}
	{{ if . }}
		<h3>Requests</h3>
		<div>
		{{ range $req := . }}
			<div class="card {{ if eq $req.Type "destructor" }}text-white bg-danger{{ end }}" id="{{ $req.Anchor }}">
				<div class="card-header">
					{{ $req.Name }} {{ template "()-light" $req.Type }} <small class="text-muted">since version {{ $req.Since }}</small>
				</div>
				<div class="card-body">
					{{ template "description-light" $req.Description.Body }}
					{{ if $req.Arguments }}<h5>Arguments</h5>{{ end }}
				</div>
				{{ template "arguments" $req.Arguments }}
			</div>
			<br>
		{{ end }}
		</div>
	{{ end }}
	{{ end }}
	
	{{ define "events" }}
	{{ if . }}
		<h3>Events</h3>
		<div>
		{{ range $ev := . }}
			<div class="card" id="{{ $ev.Anchor }}">
				<div class="card-header">
					{{ $ev.Name }} <small class="text-muted">since version {{ $ev.Since }}</small>
				</div>
				<div class="card-body">
					{{ template "description-light" $ev.Description.Body }}
				</div>
				{{ template "arguments" $ev.Arguments }}
			</div>
			<br>
		{{ end }}
		</div>
	{{ end }}
	{{ end }}
	
	{{ define "enum-entries" }}
	{{ if . }}
		{{ range $entry := . }}
		<ul class="list-group list-group-flush">
			<li class="list-group-item">
				{{ $entry.Name }} ({{ $entry.Value }}) <small class="text-muted">since version {{ $entry.Since }}</small> <br>
				{{ $entry.Summary }}
			</li>
		</ul>
		{{ end }}
	{{ end }}
	{{ end }}
	
	{{ define "enums" }}
	{{ if . }}
		<h3>Enums</h3>
		<div>
		{{ range $enum := . }}
			<div class="card" id="{{ $enum.Anchor }}">
				<div class="card-header">
					{{ if $enum.Bitfield }}Flagset{{ end }} {{ $enum.Name }} <small class="text-muted">since version {{ $enum.Since }}</small>
				</div>
				<div class="card-body">
					{{ template "description-light" $enum.Description.Body }}
					{{ if $enum.Entries }}<h5>Entries</h5>{{ end }}
				</div>
				{{ template "enum-entries" $enum.Entries }}
			</div>
			<br>
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
			<link rel="stylesheet" type="text/css" href="https://cdn.kde.org/aether-devel/bootstrap.css" />
			<link rel="stylesheet" type="text/css" href="https://cdn.kde.org/aether-devel/aether-sidebar.css" />
			<style>
				.list-group-item {
					background-color: var(--base);
				}
				.text-white {
					color: white !important;
				}
				.text-white .text-muted {
					color: rgba(255,255,255,0.5) !important;
				}
				.bg-danger {
					background-color: #dc3545 !important;
				}
				@media (prefers-color-scheme: dark) {
					.menu-title h2 {
						color: #eff0f1 !important;
					}
				}
			</style>
		</head>
		<body>
			<main class="body">				
				<header class="header">
					<nav class="navbar">
						<ol class="breadcrumb navbar-nav">
							<button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarsExampleDefault"
								aria-controls="navbarsExampleDefault" aria-expanded="false" aria-label="Toggle navigation">
								<span class="navbar-toggler-icon"></span>
							</button>
							<li class="breadcrumb-link">
								<a href="/">ReadWay Protocol Browser</a>
							</li>
							<li class="breadcrumb-link">
								<a href="#">{{ .Name }}</a>
							</li>
						</ol>
					</nav>
				</header>
				
				<aside class="sidebar">
					<div id="sidebar-header" class="menu-box">
						<div class="menu-title">
							<h2>ReadWay</h2>
						</div>
					</div>
					<div class="menu-box">
						<ul>
						{{ range $iface := .Interfaces }}
							<li>
								<h4><a href="#{{ $iface.Anchor }}">{{ $iface.Name }}</a><br></h4>
								{{ if $iface.Requests }}
								Requests:
								<ul>
								{{ range $req := $iface.Requests }}
									<li><a href="#{{ $req.Anchor }}">{{ $req.Name }}</a></li>
								{{ end }}
								</ul>
								{{ end }}
								{{ if $iface.Events }}
								Events:
								<ul>
								{{ range $ev := $iface.Events }}
									<li><a href="#{{ $ev.Anchor }}">{{ $ev.Name }}</a></li>
								{{ end }}
								</ul>
								{{ end }}
								{{ if $iface.Enums }}
								Enums:
								<ul>
								{{ range $en := $iface.Enums }}
									<li><a href="#{{ $en.Anchor }}">{{ $en.Name }}</a></li>
								{{ end }}
								</ul>
								{{ end }}
								<br>
							</li>
						{{ end }}
						</ul>
					</div>
				</aside>
				
				<article class="content">
					<h1>{{ .Name }} <small class="text-muted">protocol</small></h1>
				
					{{ template "description" .Description.Body }}
				
					{{ range $iface := .Interfaces }}
						<h2 id="{{ $iface.Anchor }}">{{ $iface.Name }} <small class="text-muted">interface version {{ $iface.Version }}</small></h2>
				
						{{ template "description" $iface.Description.Body }}
				
						{{ template "requests" $iface.Requests }}
						
						{{ template "events" $iface.Events }}
				
						{{ template "enums" $iface.Enums }}
				
						<br>
						<hr>
						<br>
					{{ end }}
				</article>
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

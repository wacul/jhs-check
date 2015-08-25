jhs-check
===

A simple checker for [PRMD](//github.com/interagent/prmd) JSON Hyper-Schema

Installation
---
```sh
go get github.com/kyoh86/jhs-check/...
```

Usage
---
<pre>jhs-check [-p|--pattern=<var>file-name-pattern</var>] <var>source-directory</var></pre>

Checking YAML files in "schemata":

```sh
jhs-check -p "\.yaml$" schemata
```

Watching sources, use a [watch](//www.npmjs.com/package/watch) package:

```sh
watch "jhs-check -p '\.yaml$' schemata" schemata
```

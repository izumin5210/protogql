package protoresolvergen

var templateResolvers = `
{{ reserveImport "context" }}
{{ reserveImport "fmt" }}

{{ range $resolver := .Resolvers }}
	func (r *{{ $resolver.Object | resolverImplementationName }}) {{ $resolver.GoFieldName }}{{ $resolver.ShortProtoResolverDeclaration }} {
		panic(fmt.Errorf("not implemented"))
	}
{{ end }}

{{ range $obj := .Objects -}}
	// {{$obj.Name}} returns {{ $obj.ResolverInterface | ref }} implementation.
	func (r *{{$.ResolverType}}) {{$obj.Name}}() {{ $obj.ResolverInterface | ref }} { return &{{ $obj | $.ResolverAdapterName }}{&{{ $obj | resolverImplementationName }}{r}} }

	type {{ $obj | resolverImplementationName }} struct { *{{ $.ResolverType }} }
{{ end }}
`

var templateResolverAdapters = `
{{ reserveImport "context" }}

{{ range $file := .Files -}}
	{{ range $resolver := $file.Resolvers }}
		func (a *{{ $resolver.Object | $file.ResolverAdapterName }}) {{ $resolver.GoFieldName }}{{ $resolver.ShortResolverDeclaration }} {
			resp, err := a.protoResolver.{{ $resolver.GoFieldName }}({{ $resolver.ArgList }})
			if err != nil {
				return nil, err
			}
			{{ if $resolver.TypeReference.IsScalar }}
				return resp, nil
			{{ else }}
				return {{ $resolver.ResolverModelFromProtoFunc }}(resp), nil
			{{ end }}
		}
	{{ end }}
{{ end }}

{{ range $file := .Files -}}
	{{ range $obj := $file.Objects -}}
		type {{ $obj | $file.ResolverAdapterName }} struct { protoResolver *{{ $obj | $file.ResolverImplementationName }} }
	{{ end -}}
{{ end }}
`

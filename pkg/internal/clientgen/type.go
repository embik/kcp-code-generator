package clientgen

import (
	"io"
	"strings"
	"text/template"

	"k8s.io/code-generator/cmd/client-gen/types"

	"github.com/kcp-dev/code-generator/pkg/parser"
	"github.com/kcp-dev/code-generator/pkg/util"
)

type TypedClient struct {
	// Group is the group in this client.
	Group types.GroupVersionInfo

	// Kind is the kinds in this file.
	Kind parser.Kind

	// APIPackagePath is the root directory under which API types exist.
	// e.g. "k8s.io/api"
	APIPackagePath string

	// SingleClusterClientPackagePath is the root directory under which single-cluster-aware clients exist.
	// e.g. "k8s.io/client-go/kubernetes"
	SingleClusterClientPackagePath string
}

func (c *TypedClient) WriteContent(w io.Writer) error {
	templ, err := template.New("typedClient").Funcs(template.FuncMap{
		"upperFirst": util.UpperFirst,
		"lowerFirst": util.LowerFirst,
		"toLower":    strings.ToLower,
	}).Parse(typedClient)
	if err != nil {
		return err
	}

	m := map[string]interface{}{
		"group":                          c.Group,
		"kind":                           &c.Kind,
		"apiPackagePath":                 c.APIPackagePath,
		"singleClusterClientPackagePath": c.SingleClusterClientPackagePath,
	}
	return templ.Execute(w, m)
}

var typedClient = `
//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by kcp code-generator. DO NOT EDIT.

package {{.group.Version.PackageName}}

import (
	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v2"

{{- if .kind.SupportsListWatch }}
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	{{.group.PackageAlias}} "{{.apiPackagePath}}/{{.group.Group.PackageName}}/{{.group.Version.PackageName}}"
{{end}}

	{{.group.PackageAlias}}client "{{.singleClusterClientPackagePath}}/typed/{{.group.Group.PackageName}}/{{.group.Version.PackageName}}"
)

// {{.kind.Plural}}ClusterGetter has a method to return a {{.kind.String}}ClusterInterface.
// A group's cluster client should implement this interface.
type {{.kind.Plural}}ClusterGetter interface {
	{{.kind.Plural}}() {{.kind.String}}ClusterInterface
}

{{ if .kind.SupportsListWatch -}}
// {{.kind.String}}ClusterInterface can operate on {{.kind.Plural}} across all clusters,
// or scope down to one cluster and return a {{if .kind.IsNamespaced}}{{.kind.Plural}}Namespacer{{else}}{{.group.PackageAlias}}client.{{.kind.String}}Interface{{end}}.
{{ else -}}
// {{.kind.String}}ClusterInterface can scope down to one cluster and return a {{if .kind.IsNamespaced}}{{.kind.Plural}}Namespacer{{else}}{{.group.PackageAlias}}client.{{.kind.String}}Interface{{end}}.
{{ end -}}
type {{.kind.String}}ClusterInterface interface {
	Cluster(logicalcluster.Name) {{if .kind.IsNamespaced}}{{.kind.Plural}}Namespacer{{else}}{{.group.PackageAlias}}client.{{.kind.String}}Interface{{end}} 
{{- if .kind.SupportsListWatch }}
	List(ctx context.Context, opts metav1.ListOptions) (*{{.group.PackageAlias}}.{{.kind.String}}List, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
{{ end -}}
}

type {{.kind.Plural | lowerFirst}}ClusterInterface struct {
	clientCache kcpclient.Cache[*{{.group.PackageAlias}}client.{{.group.GroupGoName}}{{.group.Version}}Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *{{.kind.Plural | lowerFirst}}ClusterInterface) Cluster(name logicalcluster.Name) {{if .kind.IsNamespaced}}{{.kind.Plural}}Namespacer{{else}}{{.group.PackageAlias}}client.{{.kind.String}}Interface{{end}} {
	if name == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}
{{ if .kind.IsNamespaced }}
	return &{{.kind.Plural | lowerFirst}}Namespacer{clientCache: c.clientCache, name: name}
{{ else }}
	return c.clientCache.ClusterOrDie(name).{{.kind.Plural}}()
{{ end -}}
}

{{ if .kind.SupportsListWatch }}
// List returns the entire collection of all {{.kind.Plural}} across all clusters. 
func (c *{{.kind.Plural | lowerFirst}}ClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*{{.group.PackageAlias}}.{{.kind.String}}List, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).{{.kind.Plural}}({{if .kind.IsNamespaced }}metav1.NamespaceAll{{end}}).List(ctx, opts)
}

// Watch begins to watch all {{.kind.Plural}} across all clusters.
func (c *{{.kind.Plural | lowerFirst}}ClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).{{.kind.Plural}}({{if .kind.IsNamespaced }}metav1.NamespaceAll{{end}}).Watch(ctx, opts)
}
{{ end -}}

{{ if .kind.IsNamespaced -}}
// {{.kind.Plural}}Namespacer can scope to objects within a namespace, returning a {{.group.PackageAlias}}client.{{.kind.String}}Interface.
type {{.kind.Plural}}Namespacer interface {
	Namespace(string) {{.group.PackageAlias}}client.{{.kind.String}}Interface
}

type {{.kind.Plural | lowerFirst}}Namespacer struct {
	clientCache kcpclient.Cache[*{{.group.PackageAlias}}client.{{.group.GroupGoName}}{{.group.Version}}Client]
	name logicalcluster.Name
}

func (n *{{.kind.Plural | lowerFirst}}Namespacer) Namespace(namespace string) {{.group.PackageAlias}}client.{{.kind.String}}Interface {
	return n.clientCache.ClusterOrDie(n.name).{{.kind.Plural}}(namespace)
}
{{ end -}}
`

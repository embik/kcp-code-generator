package clientgen

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/code-generator/cmd/client-gen/types"

	"github.com/kcp-dev/code-generator/v2/pkg/parser"
	"github.com/kcp-dev/code-generator/v2/pkg/util"
)

type FakeTypedClient struct {
	// Group is the group in this client.
	Group types.GroupVersionInfo

	// Kind is the kinds in this file.
	Kind parser.Kind

	// PackagePath is the package under which this client-set will be exposed.
	// TODO(skuznets) we should be able to figure this out from the output dir, ideally
	PackagePath string

	// APIPackagePath is the root directory under which API types exist.
	// e.g. "k8s.io/api"
	APIPackagePath string

	// SingleClusterClientPackagePath is the root directory under which single-cluster-aware clients exist.
	// e.g. "k8s.io/client-go/kubernetes"
	SingleClusterClientPackagePath string

	// SingleClusterApplyConfigurationsPackagePath is the root directory under which single-cluster-aware clients exist.
	// e.g. "k8s.io/client-go/applyconfigurations"
	SingleClusterApplyConfigurationsPackagePath string
}

func aliasFor(path string) ([]string, error) {
	var alias []string
	parts := strings.Split(path, "/")
	switch len(parts) {
	case 0:
		return nil, fmt.Errorf("invalid override path: %v", path)
	case 1:
		alias = []string{parts[0]}
	default:
		alias = parts[len(parts)-2:]
	}
	return alias, nil
}

func (c *FakeTypedClient) WriteContent(w io.Writer) error {
	templ, err := template.New("fakeTypedClient").Funcs(template.FuncMap{
		"upperFirst": util.UpperFirst,
		"lowerFirst": util.LowerFirst,
		"toLower":    strings.ToLower,
	}).Parse(fakeTypedClient)
	if err != nil {
		return err
	}

	importAliases := map[string]string{} // import path -> alias
	extraImports := map[string]string{   // import path -> alias
		"": c.Group.PackageAlias, // unset paths on input/output tags use the default api package
	}
	extensionVerbs := sets.New[string]()
	for i, extension := range c.Kind.Extensions {
		extensionVerbs.Insert(extension.Verb)
		for _, importPath := range []string{extension.InputPath, extension.ResultPath} {
			if importPath == "" {
				continue
			}
			alias, err := aliasFor(importPath)
			if err != nil {
				return err
			}
			extraImports[importPath] = strings.Join(alias, "")
			if importPath != filepath.Join(c.APIPackagePath, c.Group.Group.PackageName(), c.Group.Version.PackageName()) {
				importAliases[importPath] = strings.Join(alias, "")
			}
		}
		if extension.Verb == "apply" {
			var alias []string
			if extension.InputPath != "" {
				alias, err = aliasFor(extension.InputPath)
				if err != nil {
					return err
				}
			} else {
				alias = []string{c.Group.Group.PackageName(), c.Group.Version.PackageName()}
			}
			hack := filepath.Join(append([]string{c.SingleClusterApplyConfigurationsPackagePath}, alias...)...)
			extraImports[hack] = "applyconfigurations" + strings.Join(alias, "")
			if extension.InputPath != "" {
				importAliases[hack] = "applyconfigurations" + strings.Join(alias, "")
			}
			c.Kind.Extensions[i].InputPath = hack
			c.Kind.Extensions[i].InputType += "ApplyConfiguration"
		}
	}
	allVerbs := c.Kind.SupportedVerbs.Union(extensionVerbs)

	groupName := c.Group.Group.String()
	if groupName == "core" {
		groupName = ""
	}

	m := map[string]interface{}{
		"group":                          c.Group,
		"groupName":                      groupName,
		"kind":                           &c.Kind,
		"extraImports":                   extraImports,
		"importAliases":                  importAliases,
		"hasMethods":                     c.Kind.SupportedVerbs.Len() > 0 || len(c.Kind.Extensions) > 0,
		"needsApply":                     allVerbs.Has("apply"),
		"needsList":                      allVerbs.Has("list"),
		"needsPatch":                     allVerbs.Has("patch"),
		"apiPackagePath":                 c.APIPackagePath,
		"packagePath":                    c.PackagePath,
		"singleClusterClientPackagePath": c.SingleClusterClientPackagePath,
		"singleClusterApplyConfigurationsPackagePath": c.SingleClusterApplyConfigurationsPackagePath,
		"generateApplyVerbs":                          len(c.SingleClusterApplyConfigurationsPackagePath) > 0,
	}
	return templ.Execute(w, m)
}

var fakeTypedClient = `
//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by kcp code-generator. DO NOT EDIT.

package fake

import (
	"github.com/kcp-dev/logicalcluster/v3"
	kcptesting "github.com/kcp-dev/client-go/third_party/k8s.io/client-go/testing"

	"k8s.io/apimachinery/pkg/runtime/schema"

{{- range $path, $alias := .importAliases }}
	{{if $path}}{{$alias}} "{{$path}}"{{end}}
{{end -}}

{{- if .hasMethods }}
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	{{.group.PackageAlias}} "{{.apiPackagePath}}/{{.group.Group.PackageName}}/{{.group.Version.PackageName}}"
{{end}}

{{- if "watch" | .kind.SupportedVerbs.Has }}
	"k8s.io/apimachinery/pkg/watch"
{{end}}

{{- if .needsList }}
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/testing"
{{end}}

{{- if and .generateApplyVerbs .needsApply }}
	"fmt"
	"encoding/json"
{{end}}
{{- if or (and .generateApplyVerbs .needsApply) .needsPatch }}
	"k8s.io/apimachinery/pkg/types"
{{end}}
{{- if and .generateApplyVerbs ("apply" | .kind.SupportedVerbs.Has) }}
	applyconfigurations{{.group.PackageAlias}} "{{.singleClusterApplyConfigurationsPackagePath}}/{{.group.Group.PackageName}}/{{.group.Version.PackageName}}"
{{end}}

{{- if .kind.IsNamespaced}}
	kcp{{.group.PackageAlias}} "{{.packagePath}}/typed/{{.group.Group.PackageName}}/{{.group.Version.PackageName}}"
{{end}}
	{{.group.PackageAlias}}client "{{.singleClusterClientPackagePath}}/typed/{{.group.Group.PackageName}}/{{.group.Version.PackageName}}"
)

var {{.kind.Plural | lowerFirst}}Resource = schema.GroupVersionResource{Group: "{{.groupName}}", Version: "{{.group.Version.String | toLower}}", Resource: "{{.kind.Plural | toLower}}"}
var {{.kind.Plural | lowerFirst}}Kind = schema.GroupVersionKind{Group: "{{.groupName}}", Version: "{{.group.Version.String | toLower}}", Kind: "{{.kind.String}}"}

type {{.kind.Plural | lowerFirst}}ClusterClient struct {
	*kcptesting.Fake
}

// Cluster scopes the client down to a particular cluster.
func (c *{{.kind.Plural | lowerFirst}}ClusterClient) Cluster(clusterPath logicalcluster.Path) {{if .kind.IsNamespaced}}kcp{{.group.PackageAlias}}.{{.kind.Plural}}Namespacer{{else}}{{.group.PackageAlias}}client.{{.kind.String}}Interface{{end}} {
	if clusterPath == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}
{{ if .kind.IsNamespaced }}
	return &{{.kind.Plural | lowerFirst}}Namespacer{Fake: c.Fake, ClusterPath: clusterPath}
{{ else }}
	return &{{.kind.Plural | lowerFirst}}Client{Fake: c.Fake, ClusterPath: clusterPath}
{{ end -}}
}

{{ if .kind.SupportsListWatch }}
// List takes label and field selectors, and returns the list of {{.kind.Plural}} that match those selectors across all clusters.
func (c *{{.kind.Plural | lowerFirst}}ClusterClient) List(ctx context.Context, opts metav1.ListOptions) (*{{.group.PackageAlias}}.{{.kind.String}}List, error) {
	obj, err := c.Fake.Invokes(kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}ListAction({{.kind.Plural | lowerFirst}}Resource, {{.kind.Plural | lowerFirst}}Kind, logicalcluster.Wildcard, {{if .kind.IsNamespaced}}metav1.NamespaceAll, {{end}}opts), &{{.group.PackageAlias}}.{{.kind.String}}List{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &{{.group.PackageAlias}}.{{.kind.String}}List{ListMeta: obj.(*{{.group.PackageAlias}}.{{.kind.String}}List).ListMeta}
	for _, item := range obj.(*{{.group.PackageAlias}}.{{.kind.String}}List).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested {{.kind.Plural}} across all clusters.
func (c *{{.kind.Plural | lowerFirst}}ClusterClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}WatchAction({{.kind.Plural | lowerFirst}}Resource, logicalcluster.Wildcard, {{if .kind.IsNamespaced}}metav1.NamespaceAll, {{end}}opts))
}
{{ end -}}

{{ if .kind.IsNamespaced -}}
type {{.kind.Plural | lowerFirst}}Namespacer struct {
	*kcptesting.Fake
	ClusterPath logicalcluster.Path
}

func (n *{{.kind.Plural | lowerFirst}}Namespacer) Namespace(namespace string) {{.group.PackageAlias}}client.{{.kind.String}}Interface {
	return &{{.kind.Plural | lowerFirst}}Client{Fake: n.Fake, ClusterPath: n.ClusterPath, Namespace: namespace}
}
{{ end -}}

type {{.kind.Plural | lowerFirst}}Client struct {
	*kcptesting.Fake
	ClusterPath logicalcluster.Path
	{{if .kind.IsNamespaced}}Namespace string{{end}}
}

{{if "create" | .kind.SupportedVerbs.Has}}
func (c *{{.kind.Plural | lowerFirst}}Client) Create(ctx context.Context, {{.kind.String | lowerFirst}} *{{.group.PackageAlias}}.{{.kind.String}}, opts metav1.CreateOptions) (*{{.group.PackageAlias}}.{{.kind.String}}, error) {
	obj, err := c.Fake.Invokes(kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}CreateAction({{.kind.Plural | lowerFirst}}Resource, c.ClusterPath, {{if .kind.IsNamespaced}}c.Namespace, {{end}}{{.kind.String | lowerFirst}}), &{{.group.PackageAlias}}.{{.kind.String}}{})
	if obj == nil {
		return nil, err
	}
	return obj.(*{{.group.PackageAlias}}.{{.kind.String}}), err
}
{{end -}}

{{if "update" | .kind.SupportedVerbs.Has}}
func (c *{{.kind.Plural | lowerFirst}}Client) Update(ctx context.Context, {{.kind.String | lowerFirst}} *{{.group.PackageAlias}}.{{.kind.String}}, opts metav1.UpdateOptions) (*{{.group.PackageAlias}}.{{.kind.String}}, error) {
	obj, err := c.Fake.Invokes(kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}UpdateAction({{.kind.Plural | lowerFirst}}Resource, c.ClusterPath, {{if .kind.IsNamespaced}}c.Namespace, {{end}}{{.kind.String | lowerFirst}}), &{{.group.PackageAlias}}.{{.kind.String}}{})
	if obj == nil {
		return nil, err
	}
	return obj.(*{{.group.PackageAlias}}.{{.kind.String}}), err
}
{{end -}}

{{if "updateStatus" | .kind.SupportedVerbs.Has}}
func (c *{{.kind.Plural | lowerFirst}}Client) UpdateStatus(ctx context.Context, {{.kind.String | lowerFirst}} *{{.group.PackageAlias}}.{{.kind.String}}, opts metav1.UpdateOptions) (*{{.group.PackageAlias}}.{{.kind.String}}, error) {
	obj, err := c.Fake.Invokes(kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}UpdateSubresourceAction({{.kind.Plural | lowerFirst}}Resource, c.ClusterPath, "status", {{if .kind.IsNamespaced}}c.Namespace, {{end}}{{.kind.String | lowerFirst}}), &{{.group.PackageAlias}}.{{.kind.String}}{})
	if obj == nil {
		return nil, err
	}
	return obj.(*{{.group.PackageAlias}}.{{.kind.String}}), err
}
{{end -}}

{{if "delete" | .kind.SupportedVerbs.Has}}
func (c *{{.kind.Plural | lowerFirst}}Client) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.Invokes(kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}DeleteActionWithOptions({{.kind.Plural | lowerFirst}}Resource, c.ClusterPath, {{if .kind.IsNamespaced}}c.Namespace, {{end}}name, opts), &{{.group.PackageAlias}}.{{.kind.String}}{})
	return err
}
{{end -}}

{{if "deleteCollection" | .kind.SupportedVerbs.Has}}
func (c *{{.kind.Plural | lowerFirst}}Client) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}DeleteCollectionAction({{.kind.Plural | lowerFirst}}Resource, c.ClusterPath, {{if .kind.IsNamespaced}}c.Namespace, {{end}}listOpts)

	_, err := c.Fake.Invokes(action, &{{.group.PackageAlias}}.{{.kind.String}}List{})
	return err
}
{{end -}}

{{if "get" | .kind.SupportedVerbs.Has}}
func (c *{{.kind.Plural | lowerFirst}}Client) Get(ctx context.Context, name string, options metav1.GetOptions) (*{{.group.PackageAlias}}.{{.kind.String}}, error) {
	obj, err := c.Fake.Invokes(kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}GetAction({{.kind.Plural | lowerFirst}}Resource, c.ClusterPath, {{if .kind.IsNamespaced}}c.Namespace, {{end}}name), &{{.group.PackageAlias}}.{{.kind.String}}{})
	if obj == nil {
		return nil, err
	}
	return obj.(*{{.group.PackageAlias}}.{{.kind.String}}), err
}
{{end -}}

{{if "list" | .kind.SupportedVerbs.Has}}
// List takes label and field selectors, and returns the list of {{.kind.Plural}} that match those selectors.
func (c *{{.kind.Plural | lowerFirst}}Client) List(ctx context.Context, opts metav1.ListOptions) (*{{.group.PackageAlias}}.{{.kind.String}}List, error) {
	obj, err := c.Fake.Invokes(kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}ListAction({{.kind.Plural | lowerFirst}}Resource, {{.kind.Plural | lowerFirst}}Kind, c.ClusterPath, {{if .kind.IsNamespaced}}c.Namespace, {{end}}opts), &{{.group.PackageAlias}}.{{.kind.String}}List{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &{{.group.PackageAlias}}.{{.kind.String}}List{ListMeta: obj.(*{{.group.PackageAlias}}.{{.kind.String}}List).ListMeta}
	for _, item := range obj.(*{{.group.PackageAlias}}.{{.kind.String}}List).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}
{{end -}}

{{if "watch" | .kind.SupportedVerbs.Has}}
func (c *{{.kind.Plural | lowerFirst}}Client) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}WatchAction({{.kind.Plural | lowerFirst}}Resource, c.ClusterPath, {{if .kind.IsNamespaced}}c.Namespace, {{end}}opts))
}
{{end -}}

{{if "patch" | .kind.SupportedVerbs.Has}}
func (c *{{.kind.Plural | lowerFirst}}Client) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*{{.group.PackageAlias}}.{{.kind.String}}, error) {
	obj, err := c.Fake.Invokes(kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}PatchSubresourceAction({{.kind.Plural | lowerFirst}}Resource, c.ClusterPath, {{if .kind.IsNamespaced}}c.Namespace, {{end}}name, pt, data, subresources...), &{{.group.PackageAlias}}.{{.kind.String}}{})
	if obj == nil {
		return nil, err
	}
	return obj.(*{{.group.PackageAlias}}.{{.kind.String}}), err
}
{{end -}}

{{if and .generateApplyVerbs ("apply" | .kind.SupportedVerbs.Has) }}
func (c *{{.kind.Plural | lowerFirst}}Client) Apply(ctx context.Context, applyConfiguration *applyconfigurations{{.group.PackageAlias}}.{{.kind.String}}ApplyConfiguration, opts metav1.ApplyOptions) (*{{.group.PackageAlias}}.{{.kind.String}}, error) {
	if applyConfiguration == nil {
		return nil, fmt.Errorf("applyConfiguration provided to Apply must not be nil")
	}
	data, err := json.Marshal(applyConfiguration)
	if err != nil {
		return nil, err
	}
	name := applyConfiguration.Name
	if name == nil {
		return nil, fmt.Errorf("applyConfiguration.Name must be provided to Apply")
	}
	obj, err := c.Fake.Invokes(kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}PatchSubresourceAction({{.kind.Plural | lowerFirst}}Resource, c.ClusterPath, {{if .kind.IsNamespaced}}c.Namespace, {{end}}*name, types.ApplyPatchType, data), &{{.group.PackageAlias}}.{{.kind.String}}{})
	if obj == nil {
		return nil, err
	}
	return obj.(*{{.group.PackageAlias}}.{{.kind.String}}), err
}
{{end -}}

{{if and .generateApplyVerbs ("applyStatus" | .kind.SupportedVerbs.Has) }}
func (c *{{.kind.Plural | lowerFirst}}Client) ApplyStatus(ctx context.Context, applyConfiguration *applyconfigurations{{.group.PackageAlias}}.{{.kind.String}}ApplyConfiguration, opts metav1.ApplyOptions) (*{{.group.PackageAlias}}.{{.kind.String}}, error) {
	if applyConfiguration == nil {
		return nil, fmt.Errorf("applyConfiguration provided to Apply must not be nil")
	}
	data, err := json.Marshal(applyConfiguration)
	if err != nil {
		return nil, err
	}
	name := applyConfiguration.Name
	if name == nil {
		return nil, fmt.Errorf("applyConfiguration.Name must be provided to Apply")
	}
	obj, err := c.Fake.Invokes(kcptesting.New{{if not .kind.IsNamespaced}}Root{{end}}PatchSubresourceAction({{.kind.Plural | lowerFirst}}Resource, c.ClusterPath, {{if .kind.IsNamespaced}}c.Namespace, {{end}}*name, types.ApplyPatchType, data, "status"), &{{.group.PackageAlias}}.{{.kind.String}}{})
	if obj == nil {
		return nil, err
	}
	return obj.(*{{.group.PackageAlias}}.{{.kind.String}}), err
}
{{end -}}

{{range .kind.Extensions}}
{{if eq .Verb "create"}}
func (c *{{$.kind.Plural | lowerFirst}}Client) {{.Method}}(ctx context.Context, {{$.kind.String | lowerFirst}}Name string, {{if .InputType}}{{.InputType | lowerFirst}} *{{index $.extraImports .InputPath}}.{{.InputType}}{{else}}{{$.kind.String | lowerFirst}} *{{.group.PackageAlias}}.{{$.kind.String}}{{end}}, opts metav1.CreateOptions) (*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}, error) {
	obj, err := c.Fake.Invokes(kcptesting.New{{if not $.kind.IsNamespaced}}Root{{end}}CreateSubresourceAction({{$.kind.Plural | lowerFirst}}Resource, c.ClusterPath, {{$.kind.String | lowerFirst}}Name, "{{.Subresource}}", {{if $.kind.IsNamespaced}}c.Namespace, {{end}}{{if .InputType}}{{.InputType | lowerFirst}}{{else}}{{$.kind.String | lowerFirst}}{{end}}), &{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}{})
	if obj == nil {
		return nil, err
	}
	return obj.(*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}), err
}
{{end -}}

{{if eq .Verb "update"}}
func (c *{{$.kind.Plural | lowerFirst}}Client) {{.Method}}(ctx context.Context, {{$.kind.String | lowerFirst}}Name string, {{if .InputType}}{{.InputType | lowerFirst}} *{{index $.extraImports .InputPath}}.{{.InputType}}{{else}}{{$.kind.String | lowerFirst}} *{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}{{end}}, opts metav1.UpdateOptions) (*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}, error) {
	obj, err := c.Fake.Invokes(kcptesting.New{{if not $.kind.IsNamespaced}}Root{{end}}UpdateSubresourceAction({{$.kind.Plural | lowerFirst}}Resource, c.ClusterPath, "{{.Subresource}}", {{if $.kind.IsNamespaced}}c.Namespace, {{end}}{{if .InputType}}{{.InputType | lowerFirst}}{{else}}{{$.kind.String | lowerFirst}}{{end}}), &{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}{})
	if obj == nil {
		return nil, err
	}
	return obj.(*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}), err
}
{{end -}}

{{if eq .Verb "get"}}
func (c *{{$.kind.Plural | lowerFirst}}Client) {{.Method}}(ctx context.Context, {{$.kind.String | lowerFirst}}Name string, options metav1.GetOptions) (*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}, error) {
	obj, err := c.Fake.Invokes(kcptesting.New{{if not $.kind.IsNamespaced}}Root{{end}}GetSubresourceAction({{$.kind.Plural | lowerFirst}}Resource, c.ClusterPath, "{{.Subresource}}", {{if $.kind.IsNamespaced}}c.Namespace, {{end}}{{$.kind.String | lowerFirst}}Name), &{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}{})
	if obj == nil {
		return nil, err
	}
	return obj.(*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}), err
}
{{end -}}

{{if eq .Verb "list"}}
// List takes label and field selectors, and returns the list of {{$.kind.Plural}} that match those selectors.
func (c *{{$.kind.Plural | lowerFirst}}Client) {{.Method}}(ctx context.Context, {{$.kind.String | lowerFirst}}Name string, opts metav1.ListOptions) (*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}List, error) {
	obj, err := c.Fake.Invokes(kcptesting.New{{if not $.kind.IsNamespaced}}Root{{end}}ListAction({{$.kind.Plural | lowerFirst}}Resource, {{$.kind.Plural | lowerFirst}}Kind, c.ClusterPath, {{if $.kind.IsNamespaced}}c.Namespace, {{end}}opts), &{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}List{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}List{ListMeta: obj.(*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}List).ListMeta}
	for _, item := range obj.(*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}List).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}
{{end -}}

{{if eq .Verb "patch"}}
func (c *{{$.kind.Plural | lowerFirst}}Client) {{.Method}}(ctx context.Context, {{$.kind.String | lowerFirst}}Name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}, error) {
	obj, err := c.Fake.Invokes(kcptesting.New{{if not $.kind.IsNamespaced}}Root{{end}}PatchSubresourceAction({{$.kind.Plural | lowerFirst}}Resource, c.ClusterPath, {{if $.kind.IsNamespaced}}c.Namespace, {{end}}{{$.kind.String | lowerFirst}}Name, pt, data, subresources...), &{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}{})
	if obj == nil {
		return nil, err
	}
	return obj.(*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}), err
}
{{end -}}

{{if and $.generateApplyVerbs (eq .Verb "apply") }}
func (c *{{$.kind.Plural | lowerFirst}}Client) {{.Method}}(ctx context.Context, {{$.kind.String | lowerFirst}}Name string, applyConfiguration {{if .InputType}}*{{index $.extraImports .InputPath}}.{{.InputType}}{{else}}*applyconfigurations{{$.group.PackageAlias}}.{{$.kind.String}}ApplyConfiguration,{{end}}, opts metav1.ApplyOptions) (*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}, error) {
	if applyConfiguration == nil {
		return nil, fmt.Errorf("applyConfiguration provided to Apply must not be nil")
	}
	data, err := json.Marshal(applyConfiguration)
	if err != nil {
		return nil, err
	}
	name := applyConfiguration.Name
	if name == nil {
		return nil, fmt.Errorf("applyConfiguration.Name must be provided to Apply")
	}
	obj, err := c.Fake.Invokes(kcptesting.New{{if not $.kind.IsNamespaced}}Root{{end}}PatchSubresourceAction({{$.kind.Plural | lowerFirst}}Resource, c.ClusterPath, {{if $.kind.IsNamespaced}}c.Namespace, {{end}}*name, types.ApplyPatchType, data), &{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}{})
	if obj == nil {
		return nil, err
	}
	return obj.(*{{index $.extraImports .ResultPath}}.{{if .ResultType}}{{.ResultType}}{{else}}{{$.kind.String}}{{end}}), err
}
{{end -}}
{{end -}}
`

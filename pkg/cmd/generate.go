package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"milkyway/pkg/tosca"
	"milkyway/pkg/tosca/definitions"
	"net/url"
	"reflect"
)

type generateCmd struct {
	*baseBuilderCmd
}

const (
	ToscaString    = "string"
	ToscaInteger   = "integer"
	ToscaFloat     = "float"
	ToscaBool      = "boolean"
	ToscaTimestamp = "timestamp"
	ToscaNull      = "null"
)
const metaMain = `
dependencies: []

galaxy_info:
  role_name: nginx
  author: geerlingguy
  description: Nginx installation for Linux, FreeBSD and OpenBSD.
  company: "Midwestern Mac, LLC"
  license: "license (BSD, MIT)"
  min_ansible_version: 2.4
  platforms:
    - name: EL
      versions:
        - 7
        - 8
    - name: Debian
      versions:
        - all
    - name: Ubuntu
      versions:
        - trusty
        - xenial
        - focal
    - name: Archlinux
      versions:
        - all
    - name: FreeBSD
      versions:
        - 10.3
        - 10.2
        - 10.1
        - 10.0
        - 9.3
    - name: OpenBSD
      versions:
        - 5.9
        - 6.0
  galaxy_tags:
    - development
    - web
    - nginx
    - reverse
    - proxy
    - load
    - balancer
`
const defMain = `
---
# Used only for Debian/Ubuntu installation, as the -t option for apt.
nginx_default_release: ""

# Used only for Redhat installation, enables source Nginx repo.
nginx_yum_repo_enabled: true

# Use the official Nginx PPA for Ubuntu, and the version to use if so.
nginx_ppa_use: false
nginx_ppa_version: stable

# The name of the nginx package to install.
nginx_package_name: "nginx"

nginx_service_state: started
nginx_service_enabled: true

nginx_conf_template: "nginx.conf.j2"
nginx_vhost_template: "vhost.j2"

nginx_worker_processes: >-
  "{{ ansible_processor_vcpus | default(ansible_processor_count) }}"
nginx_worker_connections: "1024"
nginx_multi_accept: "off"

nginx_error_log: "/var/log/nginx/error.log warn"
nginx_access_log: "/var/log/nginx/access.log main buffer=16k flush=2m"

nginx_sendfile: "on"
nginx_tcp_nopush: "on"
nginx_tcp_nodelay: "on"

nginx_keepalive_timeout: "65"
nginx_keepalive_requests: "100"

nginx_server_tokens: "on"

nginx_client_max_body_size: "64m"

nginx_server_names_hash_bucket_size: "64"

nginx_proxy_cache_path: ""

nginx_extra_conf_options: ""
# Example extra main options, used within the main nginx's context:
#   nginx_extra_conf_options: |
#     env VARIABLE;
#     include /etc/nginx/main.d/*.conf;

nginx_extra_http_options: ""
# Example extra http options, printed inside the main server http config:
#    nginx_extra_http_options: |
#      proxy_buffering    off;
#      proxy_set_header   X-Real-IP $remote_addr;
#      proxy_set_header   X-Scheme $scheme;
#      proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
#      proxy_set_header   Host $http_host;

nginx_remove_default_vhost: false

# Listen on IPv6 (default: true)
nginx_listen_ipv6: true

nginx_vhosts: []
# Example vhost below, showing all available options:
# - listen: "80" # default: "80"
#   server_name: "example.com" # default: N/A
#   root: "/var/www/example.com" # default: N/A
#   index: "index.html index.htm" # default: "index.html index.htm"
#   filename: "example.com.conf" # Can be used to set the vhost filename.
#
#   # Properties that are only added if defined:
#   server_name_redirect: "www.example.com" # default: N/A
#   error_page: ""
#   access_log: ""
#   error_log: ""
#   extra_parameters: "" # Can be used to add extra config blocks (multiline).
#   state: "absent" # To remove the vhost configuration.

nginx_upstreams: []
# - name: myapp1
#   strategy: "ip_hash" # "least_conn", etc.
#   keepalive: 16 # optional
#   servers: {
#     "srv1.example.com",
#     "srv2.example.com weight=3",
#     "srv3.example.com"
#   }

nginx_log_format: |-
  '$remote_addr - $remote_user [$time_local] "$request" '
  '$status $body_bytes_sent "$http_referer" '
  '"$http_user_agent" "$http_x_forwarded_for"'
`

// AnsibleRole struct for Ansible role v2.9 and less
type AnsibleRole struct {
	TemplatesMain []byte
	TasksMain     []byte
	VarsMain      []byte
	DefaultsMain  []byte
	HandlersMain  []byte
	MetaMain      []byte
	FilesMain     []byte
	//LibraryMain   []byte
	Templates []byte
	Tasks     []byte
	Vars      []byte
	Defaults  []byte
	Handlers  []byte
	Meta      []byte
	Files     []byte
	//Library       []byte
	Version string
}
type ansibleRoleMeta struct {
	Meta         GalaxyMeta `yaml:"galaxy_info,omitempty"`
	Dependencies []string   `yaml:"dependencies,omitempty"`
}

type GalaxyMeta struct {
	RoleName          string     `yaml:"role_name,omitempty"`
	Author            string     `yaml:"author,omitempty"`
	Description       string     `yaml:"description,omitempty"`
	Platforms         []Platform `yaml:"platforms,omitempty"`
	MinAnsibleVersion string     `yaml:"min_ansible_version,omitempty"`
}
type Platform struct {
	Name     string   `yaml:"name,omitempty"`
	Versions []string `yaml:"versions,omitempty"`
}
type Data interface {
}
type Datayaml struct {
	Data []Data `yaml:",inline"`
}

func (b *cmdsBuilder) newGenerateCmd() *generateCmd {
	cc := &generateCmd{}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate tosca types",
		Run:   cc.generateType,
	}
	cc.baseBuilderCmd = b.newBuilderCmd(cmd)
	return cc
}
func GitHubConnect(path string) (connection GithubConnection) {
	connection = *NewConnectionBuilder(path)
	err := connection.getContents("", "")
	if err != nil {
		log.Fatal(err)
	} else if NilFields(connection.ansibleRole) {
		log.Fatal(&cmdError{
			s:         "Please, make sure that Ansible Role is correct",
			userError: true,
		})
	}
	return
}
func (c *generateCmd) generateType(cmd *cobra.Command, args []string) {
	u, err := url.Parse(c.roleURL)
	if err != nil || len(u.Path) < 3 {
		log.Fatal(err)
	}
	con := GitHubConnect(u.Path)
	generatedType := con.ansibleRole.parseRole()

	//log.Print(ansibleRoleMeta)
	//testrole := AnsibleRole{
	//	TemplatesMain: nil,
	//	TasksMain:     nil,
	//	VarsMain:      nil,
	//	DefaultsMain:  []byte(defMain),
	//	HandlersMain:  nil,
	//	MetaMain:      nil,
	//	FilesMain:     nil,
	//	Templates:     nil,
	//	Tasks:         nil,
	//	Vars:          nil,
	//	Defaults:      nil,
	//	Handlers:      nil,
	//	Meta:          []byte(metaMain),
	//	Files:         nil,
	//	Version:       "",
	//}
	//generatedType := testrole.parseRole()
	//generatedType.Properties["role_name"].Default.Value = testrole.Meta.RoleName
	//m["ansible.role."+testrole.Meta.RoleName] = generatedType
	//log.Printf("--- m:\n%v\n\n", string(connection.ansibleRole.MetaMain))
	b, _ := yaml.Marshal(generatedType)
	log.Println(string(b))

}

func (ar AnsibleRole) parseRole() map[string]tosca.NodeType {
	m := &ansibleRoleMeta{}
	errYaml := yaml.Unmarshal(ar.MetaMain, m)

	var someStruct map[string]interface{}
	errYaml2 := yaml.Unmarshal(ar.DefaultsMain, &someStruct)
	if errYaml != nil || errYaml2 != nil {
		log.Fatal(errYaml)
	}

	nt := tosca.NodeType{
		Type: tosca.Type{
			Base: tosca.Node, DerivedFrom: "tosca.nodes.Root", Version: "0.0", Description: m.Meta.Description,
			Metadata: map[string]string{"author": m.Meta.Author, "min_ansible_version": m.Meta.MinAnsibleVersion}},
		Properties:   make(map[string]definitions.PropertyDefinition),
		Attributes:   nil,
		Requirements: nil,
		Capabilities: nil,
		Interfaces:   nil,
		Artifacts:    nil,
	}
	nt.AddProperty("role_name",
		definitions.PropertyDefinition{
			Type: ToscaString, Description: "RoleName", Required: &[]bool{true}[0],
			Default: m.Meta.RoleName})
	for key, value := range someStruct {
		nt.AddProperty(key,
			definitions.PropertyDefinition{
				Type: ToscaString, Required: new(bool),
				Default: value})
	}
	return map[string]tosca.NodeType{"ansible.role." + m.Meta.RoleName: nt}
}
func NilFields(x AnsibleRole) bool {
	rv := reflect.ValueOf(&x).Elem()

	for i := 0; i < rv.NumField(); i++ {
		if !rv.Field(i).IsNil() {
			return false
		}
	}
	return true
}

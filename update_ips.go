//go:build tools
// +build tools

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"sort"
	"text/template"

	"golang.org/x/tools/imports"
)

var ipURL = flag.String("url", "https://api.github.com/meta", "URL to fetch ips from")
var localFile = flag.String("file", "", "Local file to read instead of fetching from URL")

func main() {
	flag.Parse()

	var meta GithubMeta

	if *localFile != "" {
		// Read from local file
		data, err := os.ReadFile(*localFile)
		if err != nil {
			log.Fatalf("read file err: %s", err)
		}
		err = json.Unmarshal(data, &meta)
		if err != nil {
			log.Fatalf("decode json err: %s", err)
		}
	} else {
		// Fetch from URL
		resp, err := http.Get(*ipURL)
		if err != nil {
			log.Fatalf("http fetch err: %s", err)
		}
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&meta)
		if err != nil {
			log.Fatalf("decode json err: %s", err)
		}
	}

	// Process IP ranges by service
	type IPRangeEntry struct {
		Prefix  string
		Service string
	}

	var ranges []IPRangeEntry
	services := []struct {
		name string
		ips  []string
	}{
		{"hooks", meta.Hooks},
		{"web", meta.Web},
		{"api", meta.Api},
		{"git", meta.Git},
		{"github_enterprise_importer", meta.GithubEnterpriseImporter},
		{"packages", meta.Packages},
		{"pages", meta.Pages},
		{"importer", meta.Importer},
		{"actions", meta.Actions},
		{"copilot", meta.Copilot},
		{"dependabot", meta.Dependabot},
		{"docker", meta.Docker},
	}

	// Collect all ranges with their services
	for _, svc := range services {
		for _, ip := range svc.ips {
			if ip != "" {
				ranges = append(ranges, IPRangeEntry{
					Prefix:  ip,
					Service: svc.name,
				})
			}
		}
	}

	// Sort by prefix for consistent output
	sort.Slice(ranges, func(i, j int) bool {
		if ranges[i].Prefix == ranges[j].Prefix {
			return ranges[i].Service < ranges[j].Service
		}
		return ranges[i].Prefix < ranges[j].Prefix
	})

	// Group by prefix to handle ranges that belong to multiple services
	groupedRanges := make(map[string][]string)
	for _, r := range ranges {
		groupedRanges[r.Prefix] = append(groupedRanges[r.Prefix], r.Service)
	}

	// Create sorted list of unique prefixes
	var uniquePrefixes []string
	for prefix := range groupedRanges {
		uniquePrefixes = append(uniquePrefixes, prefix)
	}
	sort.Strings(uniquePrefixes)

	// Prepare data for template
	type TemplateData struct {
		Ranges []struct {
			Prefix   string
			Services []string
		}
	}

	var templateData TemplateData
	for _, prefix := range uniquePrefixes {
		services := groupedRanges[prefix]
		// Remove duplicates and sort services
		serviceMap := make(map[string]bool)
		for _, s := range services {
			serviceMap[s] = true
		}
		var uniqueServices []string
		for s := range serviceMap {
			uniqueServices = append(uniqueServices, s)
		}
		sort.Strings(uniqueServices)

		templateData.Ranges = append(templateData.Ranges, struct {
			Prefix   string
			Services []string
		}{
			Prefix:   prefix,
			Services: uniqueServices,
		})
	}

	tmpl := template.Must(template.New("ips").Parse(ipTmpl))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, templateData)
	if err != nil {
		log.Fatalf("template err: %s", err)
	}

	fmted, err := imports.Process("ips.gen.go", buf.Bytes(), nil)
	if err != nil {
		log.Fatalf("goimports error: %s", err)
	}

	fout, err := os.Create("ips.gen.go")
	if err != nil {
		log.Fatalf("open ips.gen.go err: %s", err)
	}

	fout.Write(fmted)
	fout.Close()
}

var ipTmpl = `
// Code generated by update_ips.go. DO NOT EDIT.

package githubip

import (
	"net/netip"

	"github.com/gaissmai/cidrtree"
)

var cidrTbl = new(cidrtree.Table[IPRange])

func init() {
	var r IPRange
{{- range .Ranges}}
	r = IPRange{
		Prefix: netip.MustParsePrefix("{{.Prefix}}"),
		Services: []string{ {{range .Services}}"{{.}}", {{end}}},
	}
	cidrTbl.Insert(r.Prefix, r)
{{- end}}
}
`

type GithubMeta struct {
	VerifiablePasswordAuthentication bool              `json:"verifiable_password_authentication"`
	SSHKeyFingerprints               map[string]string `json:"ssh_key_fingerprints"`
	SSHKeys                          []string          `json:"ssh_keys"`
	Hooks                            []string          `json:"hooks"`
	Web                              []string          `json:"web"`
	Api                              []string          `json:"api"`
	Git                              []string          `json:"git"`
	GithubEnterpriseImporter         []string          `json:"github_enterprise_importer"`
	Packages                         []string          `json:"packages"`
	Pages                            []string          `json:"pages"`
	Importer                         []string          `json:"importer"`
	Actions                          []string          `json:"actions"`
	Copilot                          []string          `json:"copilot"`
	Dependabot                       []string          `json:"dependabot"`
	Docker                           []string          `json:"docker"`
}

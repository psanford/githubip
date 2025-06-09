# githubip: a Go package to check if an IP belongs to GitHub

githubip is a Go package that allows you to determine if an IP address belongs to GitHub.

A cli tool is also included in `cmd/githubip` for easily checking the status of an ip address.

## Example:

```
func ExampleIsGithubIP() {
	ips := []netip.Addr{
		netip.MustParseAddr("192.30.252.1"),
		netip.MustParseAddr("127.0.0.1"),
	}
	for _, ip := range ips {
		if IsGithubIP(ip) {
			fmt.Printf("%s is GitHub\n", ip)
		} else {
			fmt.Printf("%s is NOT GitHub\n", ip)
		}
	}
	// Output:
	// 192.30.252.1 is GitHub
	// 127.0.0.1 is NOT GitHub
}

func ExampleRange() {
	ip := netip.MustParseAddr("192.30.252.1")
	r := Range(ip)
	fmt.Println(r.Prefix)
	fmt.Println(r.Services)
	// Output:
	// 192.30.252.0/22
	// [api copilot git github_enterprise_importer hooks web]
}
```

CLI:
```
$ ./githubip 192.30.252.1
{
  "Prefix": "192.30.252.0/22",
  "Services": ["api", "copilot", "git", "github_enterprise_importer", "hooks", "web"]
}
```

## Updating the ip ranges

To update the ip ranges run: `go generate`. This will fetch from https://api.github.com/meta and regenerate the `ips.gen.go` file.

## License

MIT
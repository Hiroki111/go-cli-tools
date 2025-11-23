# pScan

A CLI tool that uses subcommands, similar to Git or Kubernetes. 
This tool executes a TCP port scan on a list of hosts similarly to the Nmap command. It allows you to add, list, and delete hosts from the list using the subcommand hosts.
It executes the scan on selected ports using the subcommand scan. Users can specify the ports using a command-line flag.
It also features command completion using the subcommand completion and manual page generation with the subcommand docs. 

To get started, install Cobra (if you haven't yet).

```bash
# Test if cobra is alreayd installed
cobra --help
```

Then, run the following to ensure that the correct directory is included in the $PATH so you can execute cobra directly.
```bash
export PATH=$(go env GOPATH)/bin:$PATH
```

## Example usage

```bash
go build

# Add localhost to the host list to scan to pScan.hosts
./pScan hosts add localhost
```
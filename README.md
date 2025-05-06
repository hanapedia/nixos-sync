# nixos-sync

A minimal CLI tool to deploy NixOS configurations to one or more remote hosts over SSH using `nixos-rebuild`.

## Features

- Reads host-to-SSH mappings from a TOML file
- Supports targeting one or more hosts
- Builds and switches configuration **on the remote host**
- Works across different architectures
- Uses `--use-remote-sudo` for privilege escalation on the target

## Requirements

- Go (for building the CLI)
- NixOS on target machines
- SSH access to each host
- `nixos-rebuild` and flakes enabled

## Example `hosts.toml`

```toml
[hosts]
host1 = "user@192.168.1.10"
host2 = "user@192.168.1.11"
```

## Usage
```bash
# Build the CLI
go build -o bin/nixos-sync

# Deploy to all hosts defined in hosts.toml
./bin/nixos-sync -t ./my/hosts.toml -f ../nixos-flake

# Deploy to one or more specific hosts defined in hosts.toml
./bin/nixos-sync -t ./my/hosts.toml -f ../nixos-flake -h host1
./bin/nixos-sync -t ./my/hosts.toml -f ../nixos-flake -h host1 -h host2
```

## Flags
| Flag       | Description                                      |
|------------|--------------------------------------------------|
| `-t`, `--toml` | Path to the hosts TOML file (default: `hosts.toml`) |
| `-f`, `--flake` | Path to the flake directory (default: `.`)         |
| `-h`, `--host`  | Hostname(s) to deploy (repeatable)                |

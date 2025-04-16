# WAGMI (We're All Gonna Make It)

A Go-based tool for automating BlockPI RPC setup and management.

## Features

- 🚀 Automated BlockPI account creation
- 📧 Temporary email handling for verification
- 🔄 Proxy management with automatic failover
- 🔑 RPC/WSS endpoint generation
- ⚡ Concurrent proxy testing for optimal performance

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/hauchu1196/wagmi.git
cd wagmi

# Install dependencies
make install

# Build the project
make build

# The binary will be available at bin/wagmi
```

### From Release

Download the latest release from the [releases page](https://github.com/hauchu1196/wagmi/releases) and extract it:

```bash
# Linux
tar -xzf wagmi-linux-amd64.tar.gz

# macOS
tar -xzf wagmi-darwin-amd64.tar.gz

# Windows
# Extract wagmi-windows-amd64.zip
```

## Commands

### Setup RPC

Create a new BlockPI account and generate RPC/WSS endpoints.

```bash
# Basic usage
./bin/wagmi setup-rpc --chain base

# Show help
./bin/wagmi setup-rpc --help
```

Options:
- `--chain, -c`: Chain name (required)

### Proxy

Test and manage proxies.

```bash
# List available proxies
./bin/wagmi proxy list

# Test proxy connection
./bin/wagmi proxy test
```

### Supported Chains

The `--chain` flag is required and must be one of the following supported chains:

| Chain Name      | Chain ID | Description           |
|----------------|----------|-----------------------|
| base           | 2030     | Base Mainnet         |
| base-sepolia   | 2041     | Base Sepolia Testnet |
| ethereum       | 1006     | Ethereum Mainnet     |
| ethereum-sepolia | 1011   | Ethereum Sepolia Testnet |

Example with validation:
```bash
# Valid
./bin/wagmi setup-rpc --chain base

# Invalid (will show error)
./bin/wagmi setup-rpc
./bin/wagmi setup-rpc --chain invalid
```

## Development

### Building

The project uses Make for building:

```bash
# Build for current platform
make build

# Build for all platforms (Linux, macOS, Windows)
make build-all

# Create release package
make release

# Clean build artifacts
make clean
```

### Version Information

The binary includes version information from git tags:
```bash
./bin/wagmi --version
```

## How It Works

1. **Account Creation**
   - Generates temporary email
   - Registers BlockPI account
   - Handles email verification

2. **Proxy Management**
   - Fetches and tests multiple proxies concurrently
   - Automatically selects working proxies
   - Implements failover for registration attempts

3. **RPC Generation**
   - Authenticates with BlockPI
   - Generates RPC/WSS endpoints
   - Returns connection details

## License

MIT 
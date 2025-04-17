# Wagmi

A CLI tool and API server for managing BlockPI RPC endpoints.

## Features

- Create BlockPI accounts
- Verify email addresses
- Get RPC/WSS endpoints
- Manage proxy lists
- REST API support

## Installation

### From Source

```bash
git clone https://github.com/hauchu1196/wagmi.git
cd wagmi
make build
```

### From Release

Download the appropriate binary for your platform from the [releases page](https://github.com/hauchu1196/wagmi/releases).

## Usage

### CLI

#### Setup RPC

```bash
# Setup RPC for Base chain
./wagmi setup-rpc --chain base

# Setup RPC for Base Sepolia chain
./wagmi setup-rpc --chain base-sepolia

# Setup RPC for Ethereum chain
./wagmi setup-rpc --chain ethereum

# Setup RPC for Ethereum Sepolia chain
./wagmi setup-rpc --chain ethereum-sepolia
```

#### API Server

```bash
# Start API server on default port (8080)
./wagmi api

# Start API server on custom port
./wagmi api --port 3000
```

### API

#### Setup RPC

```bash
curl -X POST http://localhost:8080/setup-rpc \
  -H "Content-Type: application/json" \
  -d '{"chain": "base"}'
```

Response:
```json
{
  "email": "generated-email",
  "password": "generated-password",
  "rpc": "rpc-endpoint",
  "wss": "wss-endpoint"
}
```

## Supported Chains

- Base (2030)
- Base Sepolia (2041)
- Ethereum (1006)
- Ethereum Sepolia (1011)

## Development

### Build

```bash
make build
```

### Build for All Platforms

```bash
make build-all
```

### Create Release Package

```bash
make release
```

### Clean Build Artifacts

```bash
make clean
```

## License

MIT 
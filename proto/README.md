# Flightpath Protocol Definitions

Protocol Buffer definitions for the Flightpath drone control platform.

## Overview

This folder contains the API contract definitions for Flightpath. The generated code is published as:
- **Go module**: `github.com/flightpath-dev/flightpath` - importable by Go projects
- **npm package**: `@flightpath/flightpath` - importable by TypeScript/JavaScript projects

## Using the Go Module

This repository publishes generated Go code as a Go module that can be imported by other Go projects.

### Installation

Add the module to your Go project:

```bash
go get github.com/flightpath-dev/flightpath@latest
```

Or specify a specific version:

```bash
go get github.com/flightpath-dev/flightpath@v1.0.0
```

### Importing

Import the generated code in your Go files:

```go
import (
    "github.com/flightpath-dev/flightpath/gen/go/flightpath"
    "github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)
```

### Example Usage

```go
package main

import (
    "context"
    "log"
    "net/http"
    "github.com/flightpath-dev/flightpath/gen/go/flightpath"
    "github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
    "connectrpc.com/connect"
)

func main() {
    // TODO: Add an example here
}
```

## Using the TypeScript/JavaScript Package

This repository publishes generated TypeScript/JavaScript code as an npm package that can be imported by other TypeScript/JavaScript projects.

### Installation

Install the package using npm:

```bash
npm install @flightpath/flightpath
```

Or using yarn:

```bash
yarn add @flightpath/flightpath
```

Or using pnpm:

```bash
pnpm add @flightpath/flightpath
```

### Importing

Import the generated code in your TypeScript/JavaScript files:

```typescript
import { ConnectionService } from '@flightpath/flightpath/gen/ts/flightpath/connection_connect.js';
import { ConnectRequest, ConnectResponse } from '@flightpath/flightpath/gen/ts/flightpath/connection_pb.js';
import { Position } from '@flightpath/flightpath/gen/ts/flightpath/types_pb.js';
```

### Example Usage

```typescript
import { createPromiseClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { ConnectionService } from '@flightpath/flightpath/gen/ts/flightpath/connection_connect.js';
import { ConnectRequest, ConnectResponse } from '@flightpath/flightpath/gen/ts/flightpath/connection_pb.js';
import { Position } from '@flightpath/flightpath/gen/ts/flightpath/types_pb.js';

// TODO: Add an example here
```

### Required Dependencies

Make sure to install the required peer dependencies:

```bash
npm install @bufbuild/protobuf @connectrpc/connect
```

## Development

### Prerequisites

- [Buf CLI](https://buf.build/docs/installation)
- Go 1.21 or later (for Go code generation)
- Node.js 20 or later (for TypeScript code generation and npm publishing)

### Linting
```bash
buf lint
```

### Code Generation for Go

```bash
buf generate --template buf.gen.go.yaml
```

Generates code to `gen/go/` with the correct import paths for the published module.

**Important:** After generating code, commit the generated files to the repository.

#### Go Module Dependencies

**Important:** The module path in `go.mod` and the `go_package` prefix in `buf.gen.go.yaml` must match the repository name. The import paths in the generated code must match the actual file locations in the repository.

**After first publish:**
- You can run `go mod tidy` to update dependencies and `go.sum`
- Commit the updated `go.sum` file.

### Code Generation for TypeScript

```bash
buf generate --template buf.gen.ts.yaml
```

Generates code to `gen/ts/` with the correct import paths for the published npm package.

**Important:** After generating code, commit the generated files to the repository.

### Breaking Change Detection
```bash
buf breaking --against '.git#branch=main'
```

## Releasing New Versions

This repository uses semantic versioning (v1.0.0, v1.1.0, etc.) and manual releases.

### Release Process

1. **Make your changes** to the proto files in the `flightpath/` directory

2. **Test locally**:
   ```bash
   buf lint
   buf breaking --against '.git#branch=main'
   ```

3. **Generate code**:
   ```bash
   buf generate --template buf.gen.go.yaml
   buf generate --template buf.gen.ts.yaml
   ```

4. **Update version** (if needed):
   - Update `package.json` version for npm releases
   - The Go module version is determined by the git tag

5. **Commit all changes**:
   ```bash
   git add protos/ gen/go/ gen/ts/ package.json
   git commit -m "feat: add new field to ConnectionRequest"
   ```

6. **If this is after the first publish, update go.sum**:
   ```bash
   go mod tidy
   git add go.sum
   git commit -m "chore: update go.sum"
   ```

7. **Create and push the tag**:
   ```bash
   git tag -a v1.1.0 -m "Release v1.1.0"
   git push origin main
   git push origin v1.1.0
   ```

8. **Publish to npm** (if needed):
   ```bash
   npm publish --access public
   ```

9. **Users can then update**:
   
   **Go:**
   ```bash
   go get github.com/flightpath-dev/flightpath@v1.1.0
   ```
   
   **TypeScript/JavaScript:**
   ```bash
   npm install @flightpath/flightpath@1.1.0
   ```

### Versioning Guidelines

- **MAJOR** (v2.0.0): Breaking changes to the API
- **MINOR** (v1.1.0): New features, backward compatible
- **PATCH** (v1.0.1): Bug fixes, backward compatible

### Publishing to npm

To publish the package to npm:

1. **Log in to npm** (if not already):
   ```bash
   npm login
   ```

2. **Publish the package**:
   ```bash
   npm publish --access public
   ```

   Make sure the version in `package.json` matches your git tag (without the `v` prefix).

## Publishing to Buf Schema Registry (Optional)
```bash
buf push
```

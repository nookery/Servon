# Servon - Server Management Tool

[![Coffic](https://img.shields.io/badge/Coffic-green)](https://coffic.cn)
[![Maintainer](https://img.shields.io/badge/Maintainer-blue)](https://github.com/nookery)
[![中文文档](https://img.shields.io/badge/中文文档-orange)](README-CN.md)
[![MAKEFILE](https://img.shields.io/badge/MAKEFILE-gray)](README-MAKEFILE.md)
[![Plugins](https://img.shields.io/badge/Plugins-red)](README-PLUGINS.md)
[![Update](https://img.shields.io/badge/Update-orange)](README-UPDATE.md)

Servon is a multi-functional server management tool that provides project deployment, software installation, and visual management panel features.

Currently in development stage, features may be unstable, please use with caution.

## Features

- Project Deployment (`servon deploy`)
  - Support for rapid deployment of multiple project types
  - Automatic runtime environment configuration
  - Deployment process visualization

- Software Management (`servon install`)
  - One-click installation of common server software (such as Caddy, Nginx, etc.)
  - Automatic configuration and optimization
  - Version management

- Visual Management Panel (`servon serve`)
  - System resource monitoring (CPU, memory, disk usage)
  - Website management (create, configure, deploy)
  - Docker container management
  - Intuitive Web operation interface

## Quick Installation

### Method 1: One-click Installation (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/nookery/servon/main/install.sh | bash
```

### Method 2: Manual Installation

Download the pre-compiled binary file suitable for your system from the [GitHub Releases](https://github.com/nookery/servon/releases) page:

```bash
# Download binary file (Linux amd64 example)
curl -LO https://github.com/nookery/servon/releases/latest/download/servon-linux-amd64
chmod +x servon-linux-amd64
sudo mv servon-linux-amd64 /usr/local/bin/servon
```

## Usage

### Command Line Interface

- Start management panel: `servon serve`
  - Options:
    - `-p, --port`: Specify port number (default: 8080)

- View system information: `servon info`

  - Options:
    - `-f, --format`: Output format (formatted|json|plain)

- Real-time monitoring: `servon monitor`

  - Options:
    - `-i, --interval`: Monitoring interval (seconds)

- View version: `servon version`

### Web Interface

After starting the service, visit `http://localhost:8080` to use the Web management interface.

## System Requirements

- Operating System: Linux, macOS
- Recommended Memory: >= 512MB
- Disk Space: >= 200MB

## Build Instructions

This project uses Makefile to manage the build process, see [Makefile](Makefile) for more details.

## Contributing

Welcome to submit Pull Requests or create Issues!

## License

MIT License
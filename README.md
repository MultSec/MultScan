<div align="center">
  <img width="125px" src="assets/MultScan.png" />
  <h1>MultScan</h1>
  <br/>

  <p><i>MultScan is a self-hosted, open-source, and easy-to-use malware scanner, created by <a href="https://infosec.exchange/@Pengrey">@Pengrey</a>.</i></p>
  <br />
  
</div>

>:warning: **This project is still in development, and is not ready for use.**

## Quick Start

MultScan requires the usage of Docker, although it is possible to run it without it. If you do not have Docker installed, you can install it [here](https://docs.docker.com/get-docker/).

### Installation

```bash
git clone <repo>
cd MultScan
docker build -t multscan .
```

### Features

- [ ] Dockerized
- [ ] Web UI
- [ ] File Upload
- [ ] File Hashing
- [ ] REST API
- Scanning
    - [ ] Modular Architecture
    - [ ] Proxmox API
    - Static Analysis Trigger
        - [ ] Binary Search
    - Dynamic Analysis Trigger
        - [ ] On Execution
        - [ ] On Finish

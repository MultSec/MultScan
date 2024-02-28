<div align="center">
  <img width="125px" src="assets/MultScan.png" />
  <h1>MultScan</h1>
  <br/>

  <p><i>MultScan is a self-hosted, open-source, and easy-to-use malware scanner, created by <a href="https://infosec.exchange/@Pengrey">@Pengrey</a>.</i></p>
  <br />
  
</div>

>:warning: **This project is still in development, and is not ready for use.**

## Demo

## Quick Start

MultScan requires the usage of Docker, although it is possible to run it without it. If you do not have Docker installed, you can install it [here](https://docs.docker.com/get-docker/).

### Installation

```bash
git clone <repo>
cd MultScan
docker build -t multscan .
```

### Configuration
The project requires a configuration file to be created in the webapp directory. The file should be named `config.yml` and should contain the following:

```yaml
# config.yml
config:
  connector:
    connector_type: "proxmox"
    connector_url: "https://proxmox.example.com:8006/api2/json"
    connector_user: "root@pam"
    connector_password: "password"
  machines:
    - machine_name: "machine1"
      machine_ip: "10.10.10.10"
      scan_timeout: 60
      scan_bytes: false
    - machine_name: "machine2"
      machine_ip: "10.10.10.11"
      scan_timeout: 60
      scan_bytes: true
```

#### Connector
The connector section is used to define the connector type and the credentials required to connect to the connector. As of now, the only supported connector is Proxmox but more can be easily added. The connector should contain the following fields:

- `connector_type`: The type of connector to use.
- `connector_url`: The URL of the connector.
- `connector_user`: The username to use to connect to the connector.
- `connector_password`: The password to use to connect to the connector.

#### Machines
The machines section is used to define the machines that MultScan will use for scanning. The machines should contain the following fields:

- `machine_name`: The name of the machine.
- `machine_ip`: The IP address of the machine.
- `scan_timeout`: The timeout for the scan in seconds.
- `scan_bytes`: Whether to scan the the bad bytes of the file with a scanner present on the machine.

### To-Do

- [ ] Dockerized
- [x] Web UI
- [x] File Upload
- File Info
  - [ ] Type
  - [ ] Hashing
  - [ ] Entropy
  - Public Presence
    - [ ] Check IBM X-Force
    - [ ] Check Virustotal
- [ ] REST API
- Scanning
    - [ ] Modular Architecture
    - [ ] Proxmox API
    - Static Analysis Trigger
        - [x] Binary Search
    - Dynamic Analysis Trigger
        - [ ] On Execution
        - [ ] On Finish
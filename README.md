<div align="center">
  <img width="125px" src="assets/MultScan.png" />
  <h1>MultScan</h1>
  <br/>

  <p><i>MultScan is a self-hosted, open-source, and easy-to-use malware scanner, created by <a href="https://infosec.exchange/@Pengrey">@Pengrey</a>.</i></p>
  <br />
  
</div>

## Demo

https://github.com/user-attachments/assets/543336a6-398d-4ca2-96f1-b60182e8ecb3

## Quick Start

MultScan requires the usage of Docker, although it is possible to run it without it. If you do not have Docker installed, you can install it [here](https://docs.docker.com/get-docker/).

### Installation

```bash
git clone <repo>
cd MultScan
docker build -t multscan .
docker run -p 80:8080 multscan
```

### Configuration
The project requires a configuration file to be created in the webapp directory. The file should be named `config.yml` and should contain the following:

```yaml
config:
  connector:
    type: "example"
    url: "https://example.example.com:8006/api2/json"
    user: "root@pam"
    password: "password"
  machines:
    - name: "machine1"
      ip: "10.10.10.10"
    - name: "machine2"
      ip: "10.10.10.11"
```

#### Connector
The connector section is used to define the connector type and the credentials required to connect to the connector. As of now, the only supported connector is Proxmox but more can be easily added. The connector should contain the following fields:

- `type`: The type of connector to use.
- `url`: The URL of the connector.
- `user`: The username to use to connect to the connector.
- `password`: The password to use to connect to the connector.

#### Machines
The machines section is used to define the machines that MultScan will use for scanning. The machines should contain the following fields:

- `name`: The name of the machine.
- `ip`: The IP address of the machine.

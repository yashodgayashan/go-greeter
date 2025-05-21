# Choreo Samples

## Sample Go Greeter:

#### Use the following config when creating this component in Choreo:

- Dockerfile: `go/greeter/Dockerfile`
- Docker context: `go/greeter/`

#### Running the sample locally

```shell
go run main.go
```

```mermaid
sequenceDiagram
 autonumber
 box Control Plane
 participant GraphQL
 participant CICD
 participant Secret Manager
 participant STS
 participant Mizzen
 end
 participant KV
 GraphQL->>+CICD: Trigger Build<br/>[M-IN][C-Low]
 CICD->>+Secret Manager: Fetch Secrets <br/>[M-NT][C-Low]
 Secret Manager->>+KV: Fetch Secrets <br/>[M-IN][C-Low]
 Note right of KV: This communication is<br/>done via mizzen -> mizzen agent<br>-> KV resolver -> KV
 KV-->>-Secret Manager: Secrets <br/>[M-IN][C-High]
 Secret Manager-->>-CICD: Secrets <br/>[M-NT][C-High]
 CICD->>+Mizzen: Apply Secret<br/>[M-NT][C-High]
 Mizzen-->>-CICD: Ack<br/>[M-IN][C-Low]
 Note right of CICD: Usual Build trigger process<br/>will continue
 CICD-->>-GraphQL: Acknowledge <br/>[M-IN][C-Low]
```

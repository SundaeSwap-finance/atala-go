## Atala PRISM API Wrapper

Interact directly with Atala PRISM Enterprise Agent using Golang.

See Official Atala PRISM site for more: [![Atala PRISM](https://atalaprism.io/images/atala-prism-logo-suite.svg)](https://atalaprism.io/)  


From the root directory of the repo ```atala-go```, you can run the following sample commands:
```
go run main.go getDIDs
go run main.go doCreateDID
go run main.go getDID someDid
go run main.go invitation
go run main.go issueCreds
go run main.go doCredSchema
go run main.go all
```

## Run Tests

```
cd atala
go test -v
```
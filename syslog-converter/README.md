Steps to Rebuild and Run:
Clean Build Cache:
go clean -cache -modcache -i -r

Build Filterlog Plugin:
go build -buildmode=plugin -o patterns/filterlog.so patterns/filterlog.go

Build Suricata Plugin:
go build -buildmode=plugin -o patterns/suricata.so patterns/suricata.go

Build Main Application:
go build -o syslog-converter main.go

Run the Main Application:
./syslog-converter

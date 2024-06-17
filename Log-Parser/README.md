Steps to Rebuild and Run:
Clean the Build:
$ go clean -modcache

Rebuild the Plugin:
$ go build -buildmode=plugin -o patterns/my_pattern.so my_pattern.go

Rebuild the Main Application:
$ go build -o syslog-converter main.go

Run the Main Application:
$ ./syslog-converter

Verification:
You can use the tail -f command to continuously monitor the log file:
tail -n 0 -f filterlog.log

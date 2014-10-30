#statsd-udp-splitter
A proxy server written in Go for broadcasting UDP-packets using the statsd protocol to a Graphite and Elasticsearch instance.

#BUILD
make deps
make

#TEST
make test

#CONFIG
Use ./config.json or provide another configuration file.

Configuration:
* elasticsearch instance
* graphite instance

#RUN
./gostats -p 1234

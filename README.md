#statsd-udp-splitter
A proxy server written in Go for broadcasting UDP-packets using the statsd protocol to a Graphite and Elasticsearch instance.

#BUILD IT
make deps
make

#TEST IT
make test

#CONFIG IT
Use ./config.json or provide another configuration file.

Configuration:
* elasticsearch instance
* graphite instance

#RUN IT
./gostats -p 1234

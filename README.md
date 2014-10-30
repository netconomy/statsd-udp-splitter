#UDP Stats Splitter
A proxy server written in Go for broadcasting statsd-UDP-packets to a Graphite and Elasticsearch instance.

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

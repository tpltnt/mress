mress
=====

An IRC bot to facilitate communication in groups. The bot automates certain tasks 
to keep up a good community. It is intended to help people to do good things (in 
groups) and reduce the (overall) suffering in the world. Therefore it is released 
under both the AGPLv3 (or later) and the HESSLA together to form a 
[human rights license](http://wiki.creativecommons.org/Human_rights_license) (and 
limiting freedom 0). You are not allowed to only pick the AGPL in order to execute 
or cover up human rights violations including spying on users.

license: AGPLv3 (or later) + [HESSLA](http://www.hacktivismo.com/about/hessla.php)
```
if you use this code
you and your children’s children
must make your source free
```
+ no use or modification of the software to violate human rights or spy on its user

commands
--------
* (direct message) "tell <nick>: message" - Leave a message for other offline users. It gets delivered as soon as the recipient joins the channel monitored by this mress instance.

get mress up and running
------------------------
* install a [Go toolchain](http://golang.org/doc/install)
* install go-ircevent: $go get github.com/thoj/go-ircevent
* install go-sqlite3: $go get github.com/mattn/go-sqlite3
* install goini: $go get github.com/jurka/goini
* run (optional) tests: $go test
* create executable: $go build mress.go
* check usage: $./mress --help
* run mress with flags of your choice

resources
---------
* [go-ircevent](https://github.com/thoj/go-ircevent): an event based IRC client library
* [RFC 2810: Internet Relay Chat: Architecture](https://tools.ietf.org/html/rfc2810)
* [RFC 2811: Internet Relay Chat: Channel Management](https://tools.ietf.org/html/rfc2811)
* [RFC 2812: Internet Relay Chat: Client Protocol](https://tools.ietf.org/html/rfc2812)
* [RFC 2813: Internet Relay Chat: Server Protocol](https://tools.ietf.org/html/rfc2813)
* [go-sqlite3](https://github.com/mattn/go-sqlite3): a Go-binding for [sqlite](https://sqlite.org/)
* [goini](https://github.com/jurka/goini): a [ini-style](https://en.wikipedia.org/wiki/INI_file) configfile parser ([documentation](http://godoc.org/github.com/jurka/goini)
* [M'Ress](https://en.wikipedia.org/wiki/M%27Ress)

pger -- PostgreSQL Database Driver for the Go Programming Language
==================================================================

pger is a Go language module providing rich interaction with PostgreSQL
databases. Its main features are:

- implements the sql.driver interface
- wraps the libpq driver
- uses asynchronous operations for concurrency
- can be extended with extra data types


Current status
--------------

Nothing of the above yet :) I'm learning go with this project. I'll let you
know when it'll be working.


Testing
-------

$ export GOPATH=`pwd`
$ PGER_TESTDSN="dbname=pger_test" go test pger

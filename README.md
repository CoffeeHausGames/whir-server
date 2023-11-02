# tower-of-babel_server

## Development

### Code Organization
To start developing a `Go` project the `go mod init` command needs to be run with a domain of choice. This project is using `github.com/lucas-kern/tower-of-babel_server` as the GitHub repository so has been initialized with that domain.

In order to import a package into another package for use it will need to be in the format of

`import github.com/lucas-kern/tower-of-babel_server/{PATH}/{TO}{FILE}`

With the {PATH}/{TO}{FILE} being the directory structure of the imported package.

To add other packages from GitHub that are not in this repository use `go get {PACKAGE-DOMAIN}`

If needed can run `go mod tidy` to update current packages and remove unused ones.

Read [`How to Write Go Code`](https://go.dev/doc/code) for a simple understanding of how to organize the code

### Starting the project

EXPORT ENV VARIABLES FIRST

To start the project run `go run app/main.go`

Stop the server by hitting `ctrl + c`

## Database

We will be using MongoDB community Server for development.

### Mac Instructions

#### Install

  1. Install [homebrew](https://brew.sh/) if you don't already have it.
  1. Then follow the [MongoDB homebrew install instructions](https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-os-x/)
  1. Install the Community version 6.0

#### Running MongoDB
  - EXPORT ENV VARIABLES FIRST
  - using a terminal `mongod --config /usr/local/etc/mongod.conf`
    - NOTE: Will use the terminal window to print logs
  - Run as a service `brew services start mongodb-community@6.0`

#### Connecting to Mongo instance
  
  - Run `mongosh` in a different terminal to connect to the running instance 
    - Opens up the Mongo shell

#### Disconenct from Mongo instance
  - Type .exit, exit, or exit().
  - Type quit or quit().
  - Press Ctrl + D.
  - Press Ctrl + C twice.

#### To Stop MongoDB Service
  - hit `ctrl + c` on the service and the shell if ran in the shell
  - run `brew services stop mongodb-community@6.0` to stop the backgroun service

### Windows Instructions

#### Install on Windows

  1. Donwload the MSI file for the community server version 6.0 from [MongoDB](https://www.mongodb.com/try/download/community)
  1. Follow the Installer to get it installed
    - NOTE: Select `complete` version to ensure everything you need is installed. 
    - Should install as a service
  1. For more detailed instructions follow the [MongoDB instructions](https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-windows/#install-mongodb-community-edition)
  1. After installing the MongoDB Server install the mongosh shell by following the [Mongosh Install Instructions](https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-windows/#install-mongodb-community-edition)

#### Running MongoDB
  1. Run `mongosh` from any terminal to connect to the default port of 27017
    - This is equivalent to running `mongosh "mongodb://localhost:27017"`

#### Disconnect from a Server

  To disconnect from a deployment and exit mongosh, you can:

  - Type .exit, exit, or exit().
  - Type quit or quit().
  - Press Ctrl + D.
  - Press Ctrl + C twice.

#### Stop MongoDB service

1. Open the Windows Services Manager on your Windows 11 or Windows 10 computer, do the following:
  - Right-click on the Start button to open the WinX Menu
  - Select Run
  - Type services.msc in the Run box which opens
  - Windows Services Manager will open.

1. Find the MongoDB Service in the console and `right click`
1. Select `Stop`

### Populate the Database

1. run `go run scripts/seed-database`
2. in mongo, run `show dbs`
3. check for the 'tower-of-babel' instance

## Testing

1. Tests go in the `tests` directory
1. Follow the direcotry structure of the file you are testing and add `_test.go` to the name of the file
  * ex) If you are testing `app/model/base` then the test directory should be `tests/model/base_test`
1. To run the tests run `go test ./tests/*`
  * Could run a specific test file by giving the direct path to test
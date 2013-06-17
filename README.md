# balance

Simple TCP/HTTP/HTTPS load balancer in Go

## Install

    go get github.com/darkhelmet/balance

## Usage

    # Simple tcp mode
    balance tcp -bind :4000 localhost:4001 localhost:4002

    # HTTP mode
    balance http -bind :4000 localhost:4001 localhost:4002    

## License

GNU AGPL, see `LICENSE` file.

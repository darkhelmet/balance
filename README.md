# balance

Simple TCP/HTTP load balancer in Go

## Usage

    # Simple tcp mode
    balance -bind :4000 localhost:4001 localhost:4002

    # HTTP mode
    balance -bind :4000 -mode http localhost:4001 localhost:4002    

## License

GNU AGPL, see `LICENSE` file.

# blog-engine

A primitive and opionated [gemini](https://gemini.circumlunar.space/) blog engine.

## Requirements

- [go](https://go.dev/)
- openssl or ready SSL certs for your blog domain

## Usage

The main idea is to compile blog content into the blog app binary, thanks to [go:embed](https://pkg.go.dev/embed)


1. Clone repo `git clone https://github.com/ninedraft/blog-engine.git && cd blog engine`;
2. Add documents to the `content` dir;
3. Compile with `go build`;
4. Deploy with your method of choice;
5. Run `blog-engine -addr ... -ca-cert ... etc`;

## Configuration

CLI flags
```
  -addr string
    	optional address to serve (default "localhost:1965")
  -ca-cert string
    	certificate file (default "cert.pem")
  -ca-key string
    	certificate key (default "key.pem")
  -host string
    	optional host
```

[Env variables](https://pkg.go.dev/runtime#hdr-Environment_Variables)

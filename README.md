# Whisperd

Single binary whisper application server. Pet project mostly used to practice Go and example app for DevOps and CloudOps learning.

Quick reference:

```sh
whisperd --version

whisperd --role frontend

whisperd --role backend

whisperd --role frontend --role backend

whisperd --role backend --addr 0.0.0.0:3000 --postgres-url postgres://bla/foo
```

## What Is This?

- Simple web application server that can serve backend and frontend from a single binary
- Modular roles, i.e. frontend or backend only for scalable deployment
- Embedded sqlite database vs. postgres database
- Use pure golang as much as possible, avoid frameworks, ORMs, etc.
- Test coverage
- Build toolchain (GHA)
- Docker image
- K8S deployment

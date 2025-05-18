# Quick start
> [!NOTE]
> Run all commands from the **root of the project**, unless otherwise specified.

> ## Start the Frontend

```bash

$ cd frontend

first ensure that you installed all dependencies, if not
$ npm i
if you have all dependencies installed, then proceed with
$ npm run dev
```

> ## Start the Mock Server

```bash
$ cd mockServer
$ npm run start
```

> ## Start the Signaling Server

```bash
$ cd signalingServer
$ go run main.go
```

> ## Start the CRUD Server

```bash
$ cd crud
$ # setup the DB and run migrations if needed
$ cargo run
```


# Prerequisites
> Before running the above, ensure you have:
> - Node.js (>= 16.x)
> - Go (>= 1.18)
> - Rust & Cargo (latest stable)
> - Any required .env files for each service

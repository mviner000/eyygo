# eyymi

To install dependencies:

WINDOWS: to start a new app

WINDOWS: to start the server

```bash
go build -o server ./cmd/server
```
then

```bash
./server start
```

to start a new app

```bash
go build -o manage.exe .\cmd\manage
```

```bash
.\manage.exe startapp newapp 
```

to run migrate


```bash
go build -o manage ./cmd/manage
```


```bash
 ./manage migrate
```


```bash
bun install
```

To run:

```bash
bun run index.ts
```

This project was created using `bun init` in bun v1.1.10. [Bun](https://bun.sh) is a fast all-in-one JavaScript runtime.

# Execute Gotouch

## Run the binary
Execute

```bash
gotouch
```

Gotouch will use [default properties yaml](https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/package.yaml) if
**`-f`**/**`--file`** argument is not provided. If you have custom properties yaml, execute

```bash
gotouch -f path-to-properties-yaml
```

## Run inside a docker container
Gotouch can be executed inside a docker container. Replace `$(pwd)` with the path to which you want to create the project.

```bash
docker run -it -v $(pwd):/out --rm ghcr.io/denizgursoy/gotouch:latest
```

If you want to use another properties YAML, you can execute with `-f` flag:
```bash
docker run -it -v $(pwd):/out --rm ghcr.io/denizgursoy/gotouch:latest -f url-of-properties-yaml
```

::: warning
When you use docker container, `-f` flag can only be URL. Local path is not **supported**.
:::
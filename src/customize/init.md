# Init files

You may want to execute some commands after project is created, Gotouch will execute files in the following table.

| File Name | OS |   
|---------|---|
| init.sh | Linux and MacOS|   
| init.bat|  Windows|

::: tip
Gotouch will execute init files after templating is completed. This means that you can take advantage of templating and using values in the init files. 
This allows you run some commands conditionally or change commands inside init files. As an example:

```sh
echo "Run some commands"

{{ if .isDocker }}
docker build .
{{ end }}

echo "Run some other commands"
```

Docker build command will only run if user selects docker.

:::

After the execution of the init file, Gotouch will delete both `init.sh` and `init.bat` on the root folder if they exist.

::: warning 
Make sure that init files are on the root folder. If they are not on the root folder, they will be ignored.
:::
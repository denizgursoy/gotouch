# Init files

You may want to execute some commands after project is created, Gotouch will execute files in the following table.

| File Name | OS |   
|---------|---|
| init.sh | Linux and MacOS|   
| init.bat|  Windows|

After the execution of init file, Gotouch will delete `init.sh` and `init.bat` on the root folder.

::: warning Make sure that init files on root folder. If they are not on the root folder, they will be ignored.   :::
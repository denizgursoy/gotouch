# Properties YAML

Properties yaml is a **list** of what we call [Project Structure](./project-structure).

```yaml
- name: Backend for Frontend # mandatory
  url: https://github.com/foo/bff.git # optional
- name: Microservice
  url: https://github.com/foo/microservice.git # can be a git repository
- name: Api Gateway
  url: https://foo.com/bar.tar.gz # can be a tar.gz archive file
  branch: empty 
  language: go # must be go or golang for go projects, otherwise omit the field
```

After creating your properties YAML, **you should validate your YAML** with [validate](../commands#validate) command to
check if it can be processed by Gotouch.


::: tip
Properties YAML can be stored on the root folder of the template project. Gotouch deletes file `properties.yaml` on the 
root folder. Thanks to this feature, you do not need another place to store `properties.yaml`.
:::


::: tip
Gotouch will use [default properties yaml](https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/package.yaml) if
**`-f`**/**`--file`** argument is not provided. If you always use another properties YAML, you can change default YAML
by executing `gotouch config set url path-to-new-url` command. See [config command](../commands) for more information.
:::
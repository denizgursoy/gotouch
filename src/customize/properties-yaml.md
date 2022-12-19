# Properties YAML

Properties yaml is a **list** of what we call [Project Structure](./project-structure).

```yaml
- name: Backend for Frontend # mandatory
  url: https://github.com/foo/bff.git # mandatory
- name: Microservice
  url: https://github.com/foo/microservice.git # can be a git repository
- name: Api Gateway
  url: https://foo.com/bar.tar.gz # can be a tar.gz archive file
  branch: empty 
  language: go # must be go or golang for go projects, otherwise omit the field
```

After creating your properties YAML, **you should validate your YAML** with [validate](../commands#validate) command to
check if it can be processed by Gotouch.
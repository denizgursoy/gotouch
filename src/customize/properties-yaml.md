# Properties YAML

Properties yaml is a **list** of what we call [Project Structure](./project-structure).

```yaml
- name: Backend for Frontend # mandatory
  url: https://github.com/microservice-project/microservice.git # mandatory
- name: Microservice
  url: https://github.com/bff-project/bff.git # can be a git repository
- name: Api Gateway
  url: https://github.com/denizgursoy/go-touch-projects.git # can be git repository or a tar.gz archive file
  branch: empty 
  language: go # must be go or golang for go projects, otherwise omit the field
```

After creating your properties YAML, **you should validate your YAML** with [validate](../commands#validate) command to
check if it can be processed by Gotouch.
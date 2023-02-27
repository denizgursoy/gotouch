# Dependency

Dependency can be in different formats depending on the language of selected project structure.

## Dependencies in Golang

If the language of Selected project structure is `go` or `golang` then dependencies should be a list of strings. Check
table to see which command is executed by gotouch

| depedency                          | command                                   |   
|------------------------------------|-------------------------------------------|
| github.com/labstack/echo/v4        | go get github.com/labstack/echo/v4@latest |   
| github.com/labstack/echo/v4@v4.9.1 | go get github.com/labstack/echo/v4@v4.9.1 |

## Dependencies in Other Languages

If the language of selected project structure is empty or any other value except `go` or `golang`, you can use any
format in dependencies. Gotouch merges dependencies of all selected choices and add them as an array to values
with `Dependencies` key so that you can template with dependencies. You can find example templates for different languages.

::: warning
Following examples are just suggestions. Do not forget that you can use any format in dependencies, and you can template
them however you want. If the language you use is not in the examples, in this case you can use the examples as a guide.
:::

### Java Maven

In a Maven project, you can define your dependencies as object as seen below.

<code-group>
<code-block title="Values">

```yaml
questions:
  - direction: Which DB do you want to use?
    choices:
      - choice: Postgres
        dependencies:
          - groupId: org.postgresql
            artifactId: postgresql
            version: 42.5.0
      - choice: MySQL
        dependencies:
          - groupId: mysql
            artifactId: mysql-connector-java
            version: 8.0.24
      - choice: Oracle
        dependencies:
          - groupId: com.oracle.jdbc
            artifactId: ojdbc8
            version: 12.2.0.1
```
</code-block>

<code-block title="pom.xml">

```xml
<dependencies>
    {{- range .Dependencies}}
    <dependency>
        <groupId>{{ .groupId }}</groupId>
        <artifactId>{{ .artifactId }}</artifactId>
        <version>{{ .version }}</version>
        {{ if .scope }}<scope>{{.scope}}</scope>{{- end }}
    </dependency>
    {{- end }}
</dependencies>
```
</code-block>

<code-block title="Result">

```xml
<dependencies>
    <dependency>
        <groupId>org.postgresql</groupId>
        <artifactId>postgresql</artifactId>
        <version>42.5.0</version>
    </dependency>
</dependencies>
```
</code-block>
</code-group>

### JS/Node.js
<code-group>
<code-block title="Values">

```yaml
questions:
  - direction: Which Test framework do you want to use?
    choices:
      - choice: Jest
        dependencies:
          - name: jest
            version: 29.3.1
            devDependency: true
      - choice: Mocha
        dependencies:
          - name: mocha
            version: 10.2.0
            devDependency: true
      - choice: Jasmine
        dependencies:
          - name: jasmine
            version: 4.5.0
            devDependency: true
```
</code-block>

<code-block title="package.json">

```json
{
  "dependencies": {
    {{ range .Dependencies }}
    {{- if not .devDependency -}}"{{.name}}": "{{.version}}"
    {{ end }}
    {{- end -}}
},
"devDependencies": {
    {{ range .Dependencies }}
    {{- if .devDependency -}}"{{.name}}": "{{.version}}"
    {{ end }}
    {{- end -}}
}
}
```
</code-block>

<code-block title="Result">

```json
{
  "dependencies": {
  },
  "devDependencies": {
    "jest": "29.3.1"
  }
}
```
</code-block>
</code-group>

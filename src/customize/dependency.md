# Dependency

Dependency can be in different formats depending on the language of selected project structure.

## Dependencies in Golang
If the language of Selected project structure is `go` or `golang` then
Dependencies should be a list of string. If version of a dependency is not written Gotouch will append `@latest` and
execute `go get`.

## Dependencies in other languages
If the language of selected project structure is empty or any other value except `go` or `golang`. For other languages,
 you can use any format in your dependency. For example, in a Maven project, you can define your dependencies as object as seen below.
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
Gotouch merges dependencies of all selected choices and add them as an array to values with `Dependencies` key so that you can template with dependencies.

```xml
<dependencies>
    {{- range .Dependencies}}
    <dependency>
        <groupId>{{ .groupId }}</groupId>
        <artifactId>{{ .artifactId }}</artifactId>
        <version>{{ .version }}</version>
    </dependency>
    {{- end }}
</dependencies>
```

If the user select Postgres choice, pom.xml will be generated like:

```xml
<dependencies>
    <dependency>
        <groupId>org.postgresql</groupId>
        <artifactId>postgresql</artifactId>
        <version>42.5.0</version>
    </dependency>
</dependencies>
```
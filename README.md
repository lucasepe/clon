# map

> Generate JSON or YAML using a more convenient syntax for command line _"mumbo-jumbo"_.

# Usage

```bash
$ map user = { name=foo age=:30 type=C address.zip=123 address.country=Italy }
```

generates...

```json
{
   "user": {
      "address": {
         "country": "Italy",
         "zip": "123"
      },
      "age": 30,
      "name": "foo",
      "type": "C"
   }
}
```

if you want a YAML output, use the `-f` flag:


```bash
$ map -f yaml user = { name=foo age=:30 type=C address.zip=123 address.country=Italy }
```
```yaml
user:
  address:
    country: Italy
    zip: "123"
  age: 30
  name: foo
  type: C
```

## Syntax

Syntax resembles that of JSON with a few caveats:

- uses the `=` sign instead of `:` to separate a field and its value (es. `name=foo`)
- strings are not quoted (unless they contain spaces)
- values are treated as strings by default
- other types (numbers, booleans, null) need to be prefixed with a `:` (es. `age=:30`)
- no commas required to separate elements of an object or array (es. `name=foo age=:30`)
- prefix the literals with `+` to read value from environment variables (es. `apiVersion=+API_VERSION`)

## fields

> A field is defined by: `IDENTIFIER = [:|+]VALUE` .

Example:

```bash
$ API_SECRET=abbracadabra map name=Pinco age=:30 secret=+API_SECRET
```

- the value of field `name` is a string
- the value of field `age` is a number (has a `:` prefix)
- the value of field `secret` is read from environment variables

output is:

```json
{
   "age": 30,
   "name": "Pinco",
   "secret": "abbracadabbra"
}
```

## nested fields

> A nested field is defined by: `PARENT.IDENTIFIER = [:|+]VALUE` .

Example:

```bash
$ map user.name=Pinco user.age=:30 user.address.zip=123 user.address.country=CA
```
```json
{
   "user": {
      "address": {
         "country": "CA",
         "zip": "123"
      },
      "age": 30,
      "name": "Pinco"
   }
}
```

## objects

> An object is defined by: `IDENTIFIER = { fields... }`.

Example:


```bash
$ map -f yaml user = { name=Pinco age=:30 address = { zip=123 country=CA } }
```

```yaml
user:
  address:
    country: CA
    zip: "123"
  age: 30
  name: Pinco
```

## arrays

> An array is defined by: `IDENTIFIER = [ fields...]`.

Example:

```bash
$ map tags = [ foo bar qix ]
```

```json
{
   "tags": [
      "foo",
      "bar",
      "qix"
   ]
}
```

## Why?

Let's consider the case in which we want to send a JSON payload to a ReST API via shell.

What is simpler and more readable? ...the standard way...

```bash
$ echo "{\"login\":\"my_login\",\"password\":\"my_password\"}" \
  | curl -H "Content-Type: application/json" \
         -X POST --data-binary @- \
         https://httpbin.org/anything
```

or something neat like this:

```bash
$ map login=my_login password=my_password \
  | curl -H "Content-Type: application/json" \
         -X POST --data-binary @- \
         https://httpbin.org/anything
```


## Use Cases

### Create `JSON` payload and POST it with `cURL`

```bash
$ map user = { name=Pinco age=:30 address = { zip=123 country=CA } } \
  | curl -H "Content-Type: application/json" \
         -X POST --data-binary @- \
         https://httpbin.org/anything
```

### Create a Kubernetes deployment

```bash
$ map -f yaml apiVersion=apps/v1 kind=Deployment \
      metadata = { name=nginx-deployment labels.app=nginx } \
      spec = { \
        replicas=:3 selector.matchLabels.app=nginx \
        template = { \
           metadata.labels.app=nginx \
           spec.containers = [ { name=nginx image=\"nginx:1.7.9\" ports=[containerPort=:80] } ] \
         }} \
  | curl -H "Content-Type: application/json" \
         -X POST --data-binary @- \
         http://localhost:8080/apis/apps/v1/namespaces/kube-system/deployments
```

# How to install?

In order to use the `map` command, compile it using the following command:

```bash
go get -u github.com/lucasepe/map
```

This will create the executable under your `$GOPATH/bin` directory.

## Ready-To-Use Releases 

If you don't want to compile the sourcecode yourself, [Here you can find the tool already compiled](https://github.com/lucasepe/map/releases/latest) for:

- MacOS
- Linux
- Windows


<br/>

### Credits

Thanks to [@jawher](https://github.com/jawher) for the amazing [jg](https://github.com/jawher/jg) idea and original implementation - this is a modified fork.

What I changed? 

- lexer and parser now can resolve attribute values ​​from environment variables
- additional YAML format encoding
- go modules support

# clon

> A convenient syntax to generate JSON (or YAML) for commandline _"mumbo-jumbo"_.

# Syntax Overview

Syntax resembles that of JSON with a few caveats:

- a field is a key/value pair
- fields are separated by space (one or more)
- curly braces hold objects
- square brackets hold arrays

## fields

> A field is defined by: `IDENTIFIER = VALUE` .

- field key/value pairs have a equal `=` between them as in `key = value` 
- each field is separated by space (one or more does not matter)

```sh
$ clon firstName = Scarlett lastName = Johansson
```

generates...

```json
{
   "firstName": "Scarlett",
   "lastName": "Johansson"
}
```

- values are treated as strings by default
- other types (numbers, booleans, null) need to be prefixed with a colon `:` 
  - es. `age = :30 customer = :true`
- to refer to an environment variable prefix the value with a `^` 
  - es. `bucket = ^S3_BUCKET`
  - by default the `.env` file in the current folder is used, use the `-env-file` flag to change it)


```sh
$ clon fullName = \"Scarlett Johansson\" age = :36 hot = :true
```

generates...

```json
{
   "age": 36,
   "fullName": "Scarlett Johansson",
   "hot": true
}
```

## objects

> An object is defined by: `IDENTIFIER = { fields... }`.

- begin a new object using the left curly brace `{`
- close the object with a right curly brace `}`

```sh
$ clon user = { name=foo age=:30 active=:true address = { zip=123 country=IT } }
```

generates...

```json
{
   "user": {
      "active": true,
      "address": {
         "country": "IT",
         "zip": "123"
      },
      "age": 30,
      "name": "foo"
   }
}
```

- you can also use dotted notation (and mix things)

```sh
$ clon user = { name=foo age=:30 active=:true address.zip=123 address.country=IT }
```

```sh
$ clon user.name=foo user.age=:30 user.active=:true user.address = {zip=123 country=IT}
```

```sh
$ clon user.name=foo user.age=:30 user.active=:true user.address.zip=123 user.address.country=IT
```

are all examples that generate the same JSON as above; it's up to you to find your way.


## arrays

> An array is defined by: `IDENTIFIER = [ fields...]`.

- begin a new array using the left square brace `[`
- end the array with a right quare brace `]`

```sh
$ clon tags = [ foo bar qix ]
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

You can create an array of object too:

```sh
$ clon pets = [ { name=Dash kind=cat age=:3 } {name=Harley kind=dog age=:4} ]
```

```json
{
   "pets": [
      {
         "age": 3,
         "kind": "cat",
         "name": "Dash"
      },
      {
         "age": 4,
         "kind": "dog",
         "name": "Harley"
      }
   ]
}
```

# Usage

```bash
$ clon user = { name=foo age=:30 type=C address.zip=123 address.country=Italy }
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

## Use Cases ?

### Create `JSON` payload and POST it with `cURL`

```sh
$ clon user = { name=Pinco age=:30 address = { zip=123 country=CA } } \
  | curl -H "Content-Type: application/json" \
         -X POST --data-binary @- \
         https://httpbin.org/anything
```

### Elasticsearch query string query

```sh
$ clon query.query_string.query = \"new york city\" \
  | curl -H "Content-Type: application/json" \
         --data-binary @- http://localhost:9200/_search
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
- dotenv files support
- additional YAML format encoding
- go modules support

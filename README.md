Library
=============

Simple library implementation to put in practice the [Google API design guide](https://cloud.google.com/apis/design).

# Examples

Start the API:

`make run/api`

## Create Book

### Success

#### Request

```sh
curl localhost:8081/v1/shelves/shelf1/books -d'{
  "name": "book1",
  "author": "Henrod",
}'
```

or

```sh
grpcurl -d '{
        "parent": "shelves/shelf1", 
        "book": {
            "name": "book1", 
            "author": "Henrod"
        }
    }' \
    -plaintext localhost:8080 api.v1.LibraryService/CreateBook
```

#### Response
```json
{
  "name": "book1",
  "author": "Henrod",
  "createTime": "2022-01-21T00:02:34.045823Z",
  "updateTime": "2022-01-21T00:02:34.045823Z"
}
```

### Already exists

#### Response

```json
{
  "code": 6,
  "message": "resource already exists",
  "details": [
    {
      "@type": "type.googleapis.com/google.rpc.ResourceInfo",
      "resourceType": "book",
      "resourceName": "book1",
      "owner": "shelves/shelf1",
      "description": "the book already exists in shelf"
    }
  ]
}
```


## Get Book

### Success

#### Request

```sh
curl localhost:8081/v1/shelves/shelf1/books/book1
```

or

```sh
grpcurl -d '{"name": "shelves/shelf1/books/book1"}' \
    api.v1.LibraryService/ListBooks
```

#### Response
```json
{
  "books": [
    {
      "name": "book1",
      "author": "Henrod",
      "createTime": "2022-01-20T11:01:42.327988Z",
      "updateTime": "2022-01-20T11:01:42.327988Z"
    }
  ],
  "nextPageToken": ""
}
```

### Success - Page 1

#### Request

```sh
curl localhost:8081/v1/shelves/shelf1/books\?page_size=1
```

or

```sh
grpcurl -d '{"parent": "shelves/shelf1", "page_size": 1}' \
    -plaintext localhost:8080 api.v1.LibraryService/ListBooks
```

#### Response
```json
{
  "books": [
    {
      "name": "book1",
      "author": "Henrod",
      "createTime": "2022-01-20T11:01:42.327988Z",
      "updateTime": "2022-01-20T11:01:42.327988Z"
    }
  ],
  "nextPageToken": "MQ=="
}
```

### Success - Page 2

#### Request

```sh
curl localhost:8081/v1/shelves/shelf1/books\?page_size=1\&page_token=MQ==
```

or

```sh
grpcurl -d '{
        "parent": "shelves/shelf1",
        "page_size": 1,
        "page_token": "MQ=="
    }' \
    -plaintext localhost:8080 api.v1.LibraryService/ListBooks
```

#### Response
```json
{
  "books": [
    {
      "name": "book2",
      "author": "Henrod",
      "createTime": "2022-01-21T00:02:29.938371Z",
      "updateTime": "2022-01-21T00:02:29.938371Z"
    }
  ],
  "nextPageToken": ""
}
```

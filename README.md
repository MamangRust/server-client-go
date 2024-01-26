## Crud Server Client gunakan golang

### items

```sh
curl -X GET http://localhost:8080/items

```

### create item

```sh
curl -X POST -H "Content-Type: application/json" -d '{"name":"New Item","description":"Description of new item","price":50}' http://localhost:8080/items
```

### update item

```sh
curl -X PUT -H "Content-Type: application/json" -d '{"name":"Updated Item","description":"Updated description","price":75}' http://localhost:8080/items/1

```

### delete item

```sh
curl -X DELETE http://localhost:8080/items/1

```

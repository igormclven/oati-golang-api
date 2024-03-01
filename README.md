# OATI / GO - GIN - MongoDB

Sencillo ejemplo de implementación de un API con persistencia en MongoDB.

Endpoint: https://oati-golang-api.up.railway.app/


1. Crear Asignatura

POST
```
curl --location '{{URL}}/saveSubject' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Fisica III",
    "code": "020-82"
}'
```


2. Guardar Nota

POST 
```
curl --location '{{URL}}/saveGrade' \
--header 'Content-Type: application/json' \
--data '{
    "subject": "Fisica III",
    "grade": 50
}'
```

3. Listar Notas

GET
```
curl --location 'http://localhost:8080/listAllGrades' \
--header 'Content-Type: application/json'
```




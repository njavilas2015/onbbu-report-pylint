# pylint-uploader
Esta API permite subir datos generados por el comando pylint en formato JSON a un servidor que utiliza MeiliSearch para la indexación y búsqueda de convenciones de código. La API proporciona dos endpoints: uno para subir datos y otro para obtener todos los documentos.

## Requisitos Previos
- MeiliSearch: Asegúrate de que MeiliSearch esté en funcionamiento. Puedes descargarlo y ejecutarlo desde su página oficial.
- MEILI_SEARCH_URL: URL de tu servidor de MeiliSearch (opcional, por defecto es http://127.0.0.1:7700).
- MEILI_API_KEY: La clave de API de MeiliSearch (opcional).
- DROP: Establecer a "true" si deseas eliminar el índice "pylint" en el inicio del servidor (opcional).
- PORT: El puerto en el que se ejecutará el servidor (opcional, por defecto es 9999).


## Como ejecutar 
```bash
 pylint <dir> --output-format=json | curl -X POST -H "Content-Type: application/json" --data-binary @- http://localhost:9999/
```

```yml
services:
  pylint-uploader:
    image: njavilas/pylint-uploader:latest
    ports:
      - "9999:9999"
      
    environment:
      MEILI_SEARCH_URL: "http://meilisearch:7700"  # URL del servicio de MeiliSearch
      MEILI_API_KEY: "masterKey"  # Cambia esto si tienes una clave de API
      DROP: "false"  # Cambia a "true" si deseas eliminar el índice al iniciar
      PORT: "9999"
    depends_on:
      - meilisearch

  meilisearch:
    image: getmeili/meilisearch:v1.11
    ports:
      - "7700:7700"
    environment:
      MEILI_MASTER_KEY: "masterKey"  # Cambia esto a una clave más segura
    volumes:
      - meili_data:/meili_data  # Persistencia de datos

volumes:
  meili_data:

```
```sh
docker ps
docker exec -it <nom_du_conteneur_postgres> psql -U postgres -d postgres
```
ou en local
```sh
docker exec -it $(docker ps --filter name=db -q) psql -U postgres -d postgres
```

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);
INSERT INTO users (email, password) VALUES ('alice@example.com', '1234');
```

```sh
curl -i -u alice@example.com:1234 http://localhost:8080/auth
```
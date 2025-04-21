
## К каждому таймстемпу надо дописать not null и now()

```sql
 CREATE TABLE IF NOT EXISTS pvzs (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    registration_date TIMESTAMP,
    city VARCHAR(255)
    );

CREATE TABLE IF NOT EXISTS receptions (
                                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date_time TIMESTAMP,
    pvz_id UUID,
    status VARCHAR(255),
    FOREIGN KEY (pvz_id) REFERENCES pvzs (id)
    );

CREATE TABLE IF NOT EXISTS products (
                                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date_time TIMESTAMP,
    type VARCHAR(255),
    reception_id UUID,
    FOREIGN KEY (reception_id) REFERENCES receptions (id)
    );

CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255)
    );
```
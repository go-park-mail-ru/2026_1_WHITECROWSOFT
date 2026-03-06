```mermaid
erDiagram
    user {
        id PK
        username UK
        password
        token_version
        created_at
        updated_at
    }
    note {
        id PK
        user_id FK
        title
        parent_id FK
        created_at
        updated_at
    }
    block {
        id PK
        note_id FK
        position
        block_type_id FK
        created_at
        updated_at
    }

    block_type {
        id PK
        name UK
    }

    block_state {
        id PK
        block_id FK
        formatting
        created_at
        updated_at
    }

    user ||--o{ note : "creates"
    note ||--o{ note : "has subnotes"
    note ||--o{ block : "contains"
    block }o--|| block_type : "is of type"
    block ||--o{ block_state : "has state of"
```

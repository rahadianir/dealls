CREATE TABLE IF NOT EXISTS "hr"."roles" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    "created_by" VARCHAR DEFAULT 'admin',
    "updated_by" VARCHAR
);

CREATE TABLE IF NOT EXISTS "hr"."user_roles_map" (
    "id" UUID PRIMARY KEY,
    "user_id" UUID NOT NULL,
    "role_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    "created_by" VARCHAR DEFAULT 'admin',
    "updated_by" VARCHAR,
    CONSTRAINT fk_role_user_id
        FOREIGN KEY (user_id)
        REFERENCES hr.users (id),
    CONSTRAINT fk_role_role_id
        FOREIGN KEY (role_id)
        REFERENCES hr.roles (id)
);
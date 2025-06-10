CREATE SCHEMA IF NOT EXISTS "hr";

CREATE TABLE IF NOT EXISTS "hr"."users" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "username" VARCHAR NOT NULL,
    "password" VARCHAR NOT NULL,
    "salary" DECIMAL(20,2) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    "created_by" VARCHAR DEFAULT 'admin',
    "updated_by" VARCHAR
);
CREATE TABLE IF NOT EXISTS "hr"."attendance_periods" (
    "id" UUID PRIMARY KEY,
    "start_date" TIMESTAMPTZ NOT NULL,
    "end_date" TIMESTAMPTZ NOT NULL,
    "active" BOOL DEFAULT true,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    "created_by" VARCHAR DEFAULT 'admin',
    "updated_by" VARCHAR,
    CONSTRAINT only_1_active UNIQUE (active)
);

CREATE TABLE IF NOT EXISTS "hr"."attendances" (
    "id" UUID PRIMARY KEY,
    "user_id" UUID NOT NULL,
    "attendance_time" TIMESTAMPTZ NOT NULL,
    "attendance_date" DATE NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    "created_by" VARCHAR DEFAULT 'admin',
    "updated_by" VARCHAR,
    CONSTRAINT fk_attendance_user_id
        FOREIGN KEY (user_id)
        REFERENCES hr.users (id)
);

CREATE TABLE IF NOT EXISTS "hr"."overtimes" (
    "id" UUID PRIMARY KEY,
    "user_id" UUID NOT NULL,
    "date" DATE NOT NULL,
    "hour_count" INTEGER NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    "created_by" VARCHAR DEFAULT 'admin',
    "updated_by" VARCHAR,
    CONSTRAINT fk_overtime_user_id
        FOREIGN KEY (user_id)
        REFERENCES hr.users (id)
);

CREATE TABLE IF NOT EXISTS "hr"."reimbursements" (
    "id" UUID PRIMARY KEY,
    "user_id" UUID NOT NULL,
    "amount" DECIMAL(12,2) NOT NULL,
    "description" VARCHAR,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    "created_by" VARCHAR DEFAULT 'admin',
    "updated_by" VARCHAR,
    CONSTRAINT fk_reimbursement_user_id
        FOREIGN KEY (user_id)
        REFERENCES hr.users (id)
);
CREATE TABLE IF NOT EXISTS "hr"."payrolls" (
    "id" UUID PRIMARY KEY,
    "start_date" DATE NOT NULL,
    "end_date" DATE NOT NULL,
    "total_work_days" INTEGER NOT NULL,
    "active" BOOL DEFAULT true,
    "total_salary_paid" DECIMAL(15,2) DEFAULT 0,
    "processed" BOOL DEFAULT false,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    "created_by" VARCHAR DEFAULT 'admin',
    "updated_by" VARCHAR,
    CONSTRAINT only_1_active UNIQUE (active)
    
);

CREATE TABLE IF NOT EXISTS "hr"."payslips" (
    "id" UUID PRIMARY KEY,
    "payroll_id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "base_salary" DECIMAL(12,2) NOT NULL,
    "attendance_days" INTEGER DEFAULT 0,
    "total_work_days" INTEGER DEFAULT 0,
    "overtime_hours" INTEGER DEFAULT 0,
    "overtime_bonus" DECIMAL(12,2) DEFAULT 0,
    "reimbursement_list" JSONB DEFAULT '{}',
    "total_reimbursement" DECIMAL(12,2) DEFAULT 0,
    "take_home_pay" DECIMAL(12,2) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    "created_by" VARCHAR DEFAULT 'admin',
    "updated_by" VARCHAR,
    CONSTRAINT fk_payslip_payroll_id
        FOREIGN KEY (payroll_id)
        REFERENCES hr.payrolls (id),
    CONSTRAINT fk_payslip_user_id
        FOREIGN KEY (user_id)
        REFERENCES hr.users (id)
);
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE GENDER AS ENUM (
  'MALE',
  'FEMALE'
);

CREATE TABLE IF NOT EXISTS profile (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "first_name" VARCHAR(255),
  "middle_name" VARCHAR(255),
  "last_name" VARCHAR(255),
  "gender" GENDER,
  "class" VARCHAR(255),
  "created_at" TIMESTAMP,
  "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS skill (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "profile_id" UUID,
  "skill" VARCHAR(255),
  "detail" TEXT,
  "created_at" TIMESTAMP,
  "updated_at" TIMESTAMP
);

ALTER TABLE skill
ADD CONSTRAINT fk_skill_profile
FOREIGN KEY ("profile_id")
REFERENCES "profile" ("id")
ON DELETE CASCADE
ON UPDATE CASCADE;

CREATE INDEX idx_skill_profile_id ON skill(profile_id);
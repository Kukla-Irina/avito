
-- -------------------------------------------------------------
-- TablePlus 5.4.0(504)
--
-- https://tableplus.com/
--
-- Database: postgres
-- Generation Time: 2023-08-31 14:42:21.7230
-- -------------------------------------------------------------


-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS segments_id_seq;

-- Table Definition
CREATE TABLE "public"."segments" (
    "id" int4 NOT NULL DEFAULT nextval('segments_id_seq'::regclass),
    "name" varchar(50) NOT NULL,
    "userid" int4 NOT NULL,
    PRIMARY KEY ("id")
);

INSERT INTO "public"."segments" ("id", "name", "userid") VALUES
(3, 'AVITO_VOICE_MESSAGES', 1002),
(5, 'AVITO_PERFORMANCE_VAS', 1000),
(6, 'AVITO_DISCOUNT_30', 1000),
(7, 'AVITO_VOICE_MESSAGES', 1000),
(8, 'AVITO_DISCOUNT_50', 1002);

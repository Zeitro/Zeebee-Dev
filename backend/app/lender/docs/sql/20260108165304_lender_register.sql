ALTER TABLE "public"."organizations"
    ADD COLUMN "metadata" jsonb NOT NULL DEFAULT '{}',
    ADD COLUMN "tagline" varchar(255) NOT NULL DEFAULT '';
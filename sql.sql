CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE "public"."users"
(
    "id"          uuid                                   NOT NULL DEFAULT uuid_generate_v4(),
    "username"    varchar COLLATE "pg_catalog"."default" NOT NULL,
    "password"    varchar COLLATE "pg_catalog"."default" NOT NULL,
    "email"       varchar COLLATE "pg_catalog"."default",
    "createdate"  timestamp(6),
    "updatedate"  timestamp(6),
    "deleteat"    bool                                            DEFAULT false,
    "lineID"      varchar(255) COLLATE "pg_catalog"."default",
    "admin"       bool                                   NOT NULL DEFAULT false,
    "noticeRange" int2                                   NOT NULL DEFAULT 0,
    CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "username_unique" UNIQUE ("username"),
    CONSTRAINT "email_unique" UNIQUE ("email"),
    CONSTRAINT "lineID_unique" UNIQUE ("lineID")
)
;

ALTER TABLE "public"."users"
    OWNER TO "wuchieh";

CREATE UNIQUE INDEX "idx_lineID" ON "public"."users" USING btree (
                                                                  "lineID" COLLATE "pg_catalog"."default"
                                                                  "pg_catalog"."text_ops" ASC NULLS LAST
    );

CREATE UNIQUE INDEX "idx_untitled" ON "public"."users" USING btree (
                                                                    "username" COLLATE "pg_catalog"."default"
                                                                    "pg_catalog"."text_ops" ASC NULLS LAST,
                                                                    "email" COLLATE "pg_catalog"."default"
                                                                    "pg_catalog"."text_ops" ASC NULLS LAST
    );

CREATE TABLE "public"."location"
(
    "id"        uuid                                       NOT NULL DEFAULT uuid_generate_v4(),
    "location"  point                                      NOT NULL,
    "nick_name" varchar(36) COLLATE "pg_catalog"."default" NOT NULL,
    "time"      timestamp(6)                               NOT NULL,
    "user_id"   uuid                                       NOT NULL,
    "range"     int2                                       NOT NULL DEFAULT 100,
    "delete_at" bool                                       NOT NULL DEFAULT false,
    CONSTRAINT "location_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "location_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION
)
;

ALTER TABLE "public"."location"
    OWNER TO "wuchieh";

CREATE TABLE "public"."aqi"
(
    "id"       serial,
    "location" point,
    "aqi"      numeric(5, 1),
    "time"     timestamp(6),
    CONSTRAINT "api_pkey" PRIMARY KEY ("id")
)
;

ALTER TABLE "public"."aqi"
    OWNER TO "wuchieh";

CREATE TABLE "public"."announcements"
(
    "id"         serial,
    "title"      varchar(50)[] COLLATE "pg_catalog"."default",
    "content"    text[] COLLATE "pg_catalog"."default",
    "createTime" timestamp(6),
    "hidden"     bool DEFAULT false,
    CONSTRAINT "Announcements_pkey" PRIMARY KEY ("id")
)
;

ALTER TABLE "public"."announcements"
    OWNER TO "wuchieh";

COMMENT ON TABLE "public"."announcements" IS 'SELECT SETVAL(''announcements_id_seq'', 1, false);';
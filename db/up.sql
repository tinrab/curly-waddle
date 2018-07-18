CREATE TABLE "users"
(
  "id"   CHAR(27)    NOT NULL,
  "name" VARCHAR(64) NOT NULL,

  PRIMARY KEY ("id")
);

CREATE TABLE "posts"
(
  "id"         CHAR(27)                 NOT NULL,
  "user_id"    CHAR(27)                 NOT NULL,
  "body"       TEXT                     NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,

  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
  ON DELETE CASCADE
);

DROP FUNCTION IF EXISTS read_user_posts( INT, INT );
CREATE OR REPLACE FUNCTION read_user_posts("user_ids" CHAR(27) [], "skip" INT, "take" INT)
  RETURNS TABLE(
    "id"         CHAR(27),
    "user_id"    CHAR(27),
    "body"       TEXT,
    "created_at" TIMESTAMP WITH TIME ZONE
  ) AS $$
BEGIN
  RETURN QUERY
  SELECT
    "up"."id",
    "up"."user_id",
    "up"."body",
    "up"."created_at"
  FROM (
         SELECT
           "p"."id",
           "p"."user_id",
           "p"."body",
           "p"."created_at",
           ROW_NUMBER()
           OVER (
             PARTITION BY "p"."user_id"
             ORDER BY "p"."id" ) AS "row_number"
         FROM "posts" "p"
         WHERE "p"."user_id" = ANY ("user_ids")
         ORDER BY "p"."id"
       ) "up"
  WHERE "up"."row_number" BETWEEN "skip" + 1 AND "take" + "skip";
END;
$$
LANGUAGE plpgsql;

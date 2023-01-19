CREATE TABLE "country" (
	"id" VARCHAR(100) NOT NULL,
	"name" VARCHAR(100),
	"created_at" TIMESTAMP without TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP without TIME ZONE NOT null,
	PRIMARY KEY("id")
);

CREATE TABLE "state" (
	"id" VARCHAR(100) NOT NULL,
	"name" VARCHAR(100),
	"country_id" VARCHAR(100),
	"created_at" TIMESTAMP without TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP without TIME ZONE NOT null,
	CONSTRAINT "state_country_id_fkey" FOREIGN KEY("country_id") REFERENCES "country"("id"),
	PRIMARY KEY("id")
);

CREATE TABLE "territory" (
	"id" VARCHAR(100) NOT NULL,
	"name" VARCHAR(100),
	"state_id" VARCHAR(100),
	"polygons" geography(polygon, 4326),
	"created_at" TIMESTAMP without TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP without TIME ZONE NOT null,
	CONSTRAINT "territory_state_id_fkey" FOREIGN KEY("state_id") REFERENCES "state"("id"),
	PRIMARY KEY("id")
);

CREATE TABLE "customer_tracking" (
	"id" VARCHAR(100) NOT NULL,
    "customer_id" VARCHAR(100) NOT NULL,
    "location" geography(Point,4326),
    "created_at" TIMESTAMP without TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP without TIME ZONE NOT null,
    "acknowledged" BOOL DEFAULT FALSE,
    PRIMARY KEY("id")
);

CREATE TABLE "escort_tracking_status" (
    "id" VARCHAR(100) NOT NULL,
    "name" VARCHAR(100),
    "created_at" TIMESTAMP without TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP without TIME ZONE NOT null,
    PRIMARY KEY("id")
);

CREATE TABLE "escort_tracking" (
	"id" VARCHAR(100) NOT NULL,
    "escort_id" VARCHAR(100) NOT NULL,
    "location" geography(Point,4326),
    "created_at" TIMESTAMP without TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP without TIME ZONE NOT null,
    "acknowledged" BOOL DEFAULT FALSE,
    "escort_tracking_status_id" VARCHAR(100),
    CONSTRAINT "escort_tracking_escort_tracking_status_id_fkey" FOREIGN KEY("escort_tracking_status_id") REFERENCES "escort_tracking_status"("id"),
    primary key("id")
);

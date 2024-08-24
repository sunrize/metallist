/* Application users. */
CREATE TABLE IF NOT EXISTS "users" (
	"id" INTEGER NOT NULL UNIQUE,
	"name" VARCHAR NOT NULL UNIQUE,
	PRIMARY KEY("id")	
);

/* Entities and their sub-entities, for example, book series, book series volumes, volume chapters. */
CREATE TABLE IF NOT EXISTS "entities" (
	"id" INTEGER NOT NULL UNIQUE,
	"name" VARCHAR NOT NULL,
	"description" TEXT,
	PRIMARY KEY("id")	
);

/* Lists of entities. */
CREATE TABLE IF NOT EXISTS "lists" (
	"id" INTEGER NOT NULL UNIQUE,
	"name" VARCHAR NOT NULL UNIQUE,
	"owner_id" INTEGER NOT NULL,
	PRIMARY KEY("id"),
	FOREIGN KEY ("owner_id") REFERENCES "users"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS "list_actions" (
	"id" VARCHAR NOT NULL UNIQUE,
	PRIMARY KEY("id")	
);

CREATE TABLE IF NOT EXISTS "entity_list_memberships" (
	"entity_id" INTEGER NOT NULL,
	"list_id" INTEGER NOT NULL,
	PRIMARY KEY("entity_id", "list_id"),
	FOREIGN KEY ("entity_id") REFERENCES "entities"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
	FOREIGN KEY ("list_id") REFERENCES "lists"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS "entity_children" (
	"parent_id" INTEGER NOT NULL,
	"child_id" INTEGER NOT NULL,
	PRIMARY KEY("parent_id", "child_id"),
	FOREIGN KEY ("parent_id") REFERENCES "entities"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
	FOREIGN KEY ("child_id") REFERENCES "entities"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

/* External services */
CREATE TABLE IF NOT EXISTS "services" (
	"id" INTEGER NOT NULL UNIQUE,
	"id_string" VARCHAR NOT NULL,
	"url" VARCHAR NOT NULL,
	"name" VARCHAR NOT NULL,
	PRIMARY KEY("id")	
);

CREATE TABLE IF NOT EXISTS "service_entities" (
	"service_id" INTEGER NOT NULL,
	"id" INTEGER NOT NULL,
	"last_sync" TIMESTAMP,
	"id_string" VARCHAR,
	PRIMARY KEY("service_id", "id"),
	FOREIGN KEY ("service_id") REFERENCES "services"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS "service_lists" (
	"service_id" INTEGER NOT NULL,
	"id" INTEGER NOT NULL,
	"last_sync" TIMESTAMP,
	"id_string" VARCHAR,
	PRIMARY KEY("service_id", "id"),
	FOREIGN KEY ("service_id") REFERENCES "services"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS "service_entity_list_memberships" (
	"entity_id" INTEGER NOT NULL,
	"list_id" INTEGER NOT NULL,
	PRIMARY KEY("entity_id", "list_id"),
	FOREIGN KEY ("entity_id") REFERENCES "service_entities"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
	FOREIGN KEY ("list_id") REFERENCES "service_lists"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS "service_entity_children" (
	"parent_id" INTEGER NOT NULL,
	"child_id" INTEGER NOT NULL,
	PRIMARY KEY("parent_id", "child_id"),
	FOREIGN KEY ("parent_id") REFERENCES "service_entities"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
	FOREIGN KEY ("child_id") REFERENCES "service_entities"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS "list_connections" (
	"list_id" INTEGER NOT NULL,
	"service_list_id" INTEGER NOT NULL,
	"main" BOOLEAN NOT NULL DEFAULT FALSE,
	"update" BOOLEAN NOT NULL DEFAULT FALSE,
	PRIMARY KEY("list_id", "service_list_id"),
	FOREIGN KEY ("list_id") REFERENCES "lists"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
	FOREIGN KEY ("service_list_id") REFERENCES "service_lists"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS "entity_connections" (
	"entity_id" INTEGER NOT NULL,
	"service_entity_id" INTEGER NOT NULL,
	PRIMARY KEY("entity_id", "service_entity_id"),
	FOREIGN KEY ("entity_id") REFERENCES "entities"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
	FOREIGN KEY ("service_entity_id") REFERENCES "service_entities"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

/* List history after sync. */
CREATE TABLE IF NOT EXISTS "service_list_history" (
	"id" INTEGER NOT NULL UNIQUE,
	"timestamp" TIMESTAMP NOT NULL,
	"list_id" INTEGER NOT NULL,
	"entity_id" INTEGER NOT NULL,
	"action_id" VARCHAR NOT NULL,
	PRIMARY KEY("id"),
	FOREIGN KEY ("list_id") REFERENCES "service_lists"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
	FOREIGN KEY ("entity_id") REFERENCES "service_entities"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
	FOREIGN KEY ("action_id") REFERENCES "list_actions"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

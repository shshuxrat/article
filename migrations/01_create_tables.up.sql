CREATE TABLE IF NOT EXISTS "author"(
	id SERIAl  PRIMARY KEY   NOT NULL,
	firstname VARCHAR (100)     NOT NULL,
	lastname VARCHAR (100)     NOT NULL,
	created_at TIMESTAMP DEFAULT(Now()),
	updated_at TIMESTAMP DEFAULT(Now())
 );
 
 
  CREATE TABLE IF NOT EXISTS "article"(
	id SERIAl  PRIMARY KEY   NOT NULL,
	"title" VARCHAR (250)     NOT NULL,
	"body" TEXT,
	"author_id" INT,
	created_at TIMESTAMP DEFAULT(Now()),
	updated_at TIMESTAMP DEFAULT(Now()),
	CONSTRAINT fk_author FOREIGN KEY(author_id) REFERENCES author(id)
 ); 
INSERT INTO "author" (firstname, lastname) VALUES ('Jason','Moiron') ON CONFLICT DO NOTHING;
INSERT INTO "author" (firstname, lastname) VALUES ('John','Doe') ON CONFLICT DO NOTHING;

INSERT INTO "article" (title, body,author_id) VALUES ('Lorem1','Lorem Ipsum1',1) ON CONFLICT DO NOTHING;
INSERT INTO "article" (title, body,author_id) VALUES ('Lorem2','Lorem Ipsum2',2) ON CONFLICT DO NOTHING;

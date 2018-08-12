DROP TABLE IF EXISTS News, Source;

CREATE TABLE News (
  app_id int,
  author varchar(255),
  title varchar(255),
  description varchar(800),
  url varchar(500),
  image_url varchar(500),
  published_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(app_id)
);

CREATE TABLE Source (
  news_id int REFERENCES News(app_id),
  id varchar(100),
  name varchar(255)
);

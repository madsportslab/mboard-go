create table if not exists formats(
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR NOT NULL,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP
);

create table if not exists media(
  id INTEGER NOT NULL PRIMARY KEY,
  format_id INTEGER,
  key VARCHAR NOT NULL UNIQUE,
  meta VARCHAR,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(format_id) REFERENCES formats(id)
);

INSERT INTO formats(name) VALUES('MP4');
INSERT INTO formats(name) VALUES('MOV');
INSERT INTO formats(name) VALUES('JPG');
INSERT INTO formats(name) VALUES('PNG');
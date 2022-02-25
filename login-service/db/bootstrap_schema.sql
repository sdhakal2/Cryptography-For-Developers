CREATE TABLE IF NOT EXISTS signuplogin (
  username TEXT NOT NULL,
  passwordhash TEXT NOT NULL,
  salt TEXT NOT NULL,
  PRIMARY KEY (username)
);

CREATE INDEX IF NOT EXISTS usernames on signuplogin
  (username);

INSERT INTO signuplogin VALUES
  ('defaultuser', '6507eaf0e1b9bd2db7f8aab0769a829ac2bba85801ea2d5722f1dc7503356861','ea4de1afac17f22dac11de9b7676cea817389c8549a187965e114691da921f0e');

--username:
--defaultuser
--Salt:
--ea4de1afac17f22dac11de9b7676cea817389c8549a187965e114691da921f0e
--password:
--defaultpassword
--Password Hash:
--6507eaf0e1b9bd2db7f8aab0769a829ac2bba85801ea2d5722f1dc7503356861
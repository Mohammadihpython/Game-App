
-- +migrate Up
-- you can also  use another or better way base on your mysql version
-- mysql8 set mysqluser for all values Dont  change the order of the Enum
Alter TABLE users ADD COLUMN role ENUM('mysqluser','admin') not null;




-- +migrate Down
Alter TABLE users drop column password ;

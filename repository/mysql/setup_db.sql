
Alter TABLE users add column password varchar(255) not null ;
-- insert object

insert into users(id,name) values(1,"hamed");

-- show records
select * from users;



--
select * from users where phone_number= "0912"


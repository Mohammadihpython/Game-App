CREATE TABLE users (
    id int primary key AUTO_INCREMENT ,
    name varchar(255) not null ,
    phone_number varchar(255) not null unique,
    password text not null ,
    created_at TIMESTAMP DEFAULT  CURRENT_TIMESTAMP
);

-- insert object

insert into users(id,name) values(1,"hamed");

-- show records
select * from users;



--
select * from users where phone_number= "0912"


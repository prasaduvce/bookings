create TABLE users (
id serial PRIMARY KEY NOT null, 
first_name varchar(255) DEFAULT '' NOT NULL,
last_name varchar(255) DEFAULT '' NOT NULL ,
email varchar(255) DEFAULT '' NOT NULL,
password varchar(255) NOT NULL,
access_level integer DEFAULT 1
);

alter table users 
add column created_at timestamp;

alter table users
add column updated_at timestamp; 

CREATE TABLE reservations (
id serial PRIMARY KEY NOT null, 
first_name varchar(255) DEFAULT '' NOT NULL,
last_name varchar(255) DEFAULT '' NOT NULL ,
email varchar(255) DEFAULT '' NOT NULL,
phone varchar(255) NOT NULL,
start_date date,
end_date date,
room_id integer,
created_at timestamp,
updated_at timestamp
);

create table rooms (
id serial PRIMARY KEY NOT null,
room_name varchar(255) DEFAULT '' NOT NULL
);

alter table rooms
add column created_at timestamp,
add column updated_at timestamp;

create table restrictions (
id serial PRIMARY KEY NOT null,
restriction_name varchar(255) DEFAULT '' NOT null,
created_at timestamp,
updated_at timestamp
);

create table room_restrictions (
id serial PRIMARY KEY NOT null,
start_date date,
end_date date,
room_id integer,
reservation_id integer,
restriction_id integer,
created_at timestamp,
updated_at timestamp
);

alter table reservations
add constraint 
fk_reservations_rooms 
FOREIGN key (room_id) references rooms(id) 
on delete cascade on update cascade;

alter table room_restrictions
add constraint 
fk_room_restrictions_rooms 
FOREIGN key (room_id) references rooms(id) 
on delete cascade on update cascade;

alter table room_restrictions
add constraint 
fk_room_restrictions_restrictions
FOREIGN key (restriction_id) references restrictions(id) 
on delete cascade on update cascade;

alter table users 
add constraint idx_users_email_unique unique (email);

create index idx_users_email 
on users(email);

create index idx_start_date_end_date
on room_restrictions(start_date, end_date);

create index idx_room_id
on room_restrictions(room_id);

create index idx_reservation_id
on room_restrictions(reservation_id);

create index idx_email
on reservations(email);

create index idx_last_name
on reservations(last_name);


alter table room_restrictions
add constraint 
fk_room_restrictions_reservations
FOREIGN key (reservation_id) references reservations(id) 
on delete cascade on update cascade;

ALTER TABLE room_restrictions ALTER COLUMN restriction_id SET NOT NULL;
ALTER TABLE room_restrictions ALTER COLUMN created_at SET NOT NULL;
ALTER TABLE room_restrictions ALTER COLUMN end_date SET NOT NULL;
ALTER TABLE room_restrictions ALTER COLUMN room_id SET NOT NULL;
ALTER TABLE room_restrictions ALTER COLUMN start_date SET NOT NULL;
ALTER TABLE room_restrictions ALTER COLUMN updated_at SET NOT NULL;


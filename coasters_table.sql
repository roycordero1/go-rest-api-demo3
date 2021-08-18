use go_demo;
CREATE TABLE coasters(
    id bigint unsigned not null primary key auto_increment,
    name VARCHAR(255) NOT NULL,
    manufacturer VARCHAR(255) NOT NULL,
    in_park VARCHAR(255) NOT NULL,
    height INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS voting_users (
    id serial NOT NULL,
    username varchar(255) NOT NULL primary key,
    password varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS candidates (
    id serial NOT NULL,
    candidate_name varchar(255) NOT NULL,
    party_name varchar(255) NOT NULL primary key
);

CREATE TABLE IF NOT EXISTS admins (
    public_key varchar(255) NOT NULL,
    username varchar(255) NOT NULL primary key,
    vote boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS user_sign (
    signature varchar(255) NOT NULL,
    username varchar(255) NOT NULL primary key,
    sign_hash varchar(255) NOT NULL
    );


DO $$ BEGIN
    CREATE TYPE USER_ROLE AS ENUM (
        'participant',
        'administrator'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS users (
    username VARCHAR (256) PRIMARY KEY,
    role     USER_ROLE     NOT NULL
);

DO $$ BEGIN
    CREATE TYPE SKILL_TYPE AS ENUM (
        'Инженерные',
        'Исследовательские',
        'Социальные',
        'Творческие',
        'Спортивные'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS characters (
    group_name VARCHAR (8)   PRIMARY KEY,
    username   VARCHAR (256) NOT NULL UNIQUE,
    started_at TIMESTAMP     NULL,

    CONSTRAINT fk_username
        FOREIGN KEY ( username )
            REFERENCES users ( username )
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS activities (
    name        VARCHAR (256) PRIMARY KEY,
    full_name   VARCHAR (256) NOT NULL,
    description TEXT          NULL,
    location    VARCHAR (256) NULL,
    skills      SKILL_TYPE[]  NOT NULL,
    max_points  INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS admins (
    activity_name VARCHAR (256) NOT NULL,
    username      VARCHAR (256) NOT NULL,

    PRIMARY KEY (activity_name, username),

    CONSTRAINT fk_activity_name
        FOREIGN KEY ( activity_name )
            REFERENCES activities ( name )
            ON DELETE CASCADE,

    CONSTRAINT fk_username
        FOREIGN KEY ( username )
            REFERENCES users ( username )
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS activity_slots (
    activity_name VARCHAR (256) NOT NULL,
    start         TIMESTAMP     NOT NULL,
    end_          TIMESTAMP     NOT NULL,
    group_name    VARCHAR (256) NULL,

    PRIMARY KEY ( activity_name, start ),

    CONSTRAINT fk_activity_name
        FOREIGN KEY ( activity_name )
            REFERENCES activities ( name )
            ON DELETE CASCADE,

    CONSTRAINT fk_group_name
        FOREIGN KEY ( group_name )
            REFERENCES characters ( group_name )
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS grades (
    group_name    VARCHAR (8)   NOT NULL,
    skill_type    SKILL_TYPE    NOT NULL,
    points        INTEGER       NOT NULL,
    activity_name VARCHAR (256) NOT NULL,
    time          TIMESTAMP     NOT NULL,

    CONSTRAINT fk_group_name
        FOREIGN KEY ( group_name )
            REFERENCES characters ( group_name )
            ON DELETE CASCADE,

    CONSTRAINT fk_activity_name
        FOREIGN KEY ( activity_name )
            REFERENCES activities ( name )
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS character_slots (
    group_name    VARCHAR (8)   NOT NULL,
    start         TIMESTAMP     NOT NULL,
    end_          TIMESTAMP     NOT NULL,
    activity_name VARCHAR (256) NULL,

    PRIMARY KEY ( group_name, start ),

    CONSTRAINT fk_activity_name
        FOREIGN KEY ( activity_name )
            REFERENCES activities ( name )
            ON DELETE CASCADE,

    CONSTRAINT fk_group_name
        FOREIGN KEY ( group_name )
            REFERENCES characters ( group_name )
            ON DELETE CASCADE
);

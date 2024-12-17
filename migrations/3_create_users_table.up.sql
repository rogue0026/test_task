CREATE TABLE IF NOT EXISTS registered_users (
    id UUID PRIMARY KEY,
    login VARCHAR(40) NOT NULL UNIQUE,
    name VARCHAR(40) NOT NULL,
    email VARCHAR(40) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL
);

CREATE OR REPLACE PROCEDURE register_user(
    _id UUID,
    _login VARCHAR(40),
    _name VARCHAR(40),
    _email VARCHAR(40),
    _password VARCHAR(128))
LANGUAGE plpgsql
AS
$$
    BEGIN
        INSERT INTO registered_users (id, login, name, email, password)
        VALUES (_id, _login, _name, _email, _password);
    END;
$$;

CREATE OR REPLACE FUNCTION get_registered_user_by_login(_login VARCHAR(40))
RETURNS SETOF registered_users
LANGUAGE plpgsql
AS
$$
    BEGIN
        RETURN QUERY
        SELECT
            id,
            login,
            name,
            email,
            password
        FROM registered_users
        WHERE login = _login;
    END;
$$;

CREATE OR REPLACE FUNCTION get_registered_user_by_email(_email VARCHAR(40))
RETURNS SETOF registered_users
LANGUAGE plpgsql
AS
$$
BEGIN
    RETURN QUERY
        SELECT
            id,
            login,
            name,
            email,
            password
        FROM registered_users
        WHERE login = _email;
END;
$$;
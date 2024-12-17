CREATE OR REPLACE FUNCTION create_unregistered_user(
    _login VARCHAR(40),
    _name VARCHAR(40),
    _email VARCHAR(40),
    _password VARCHAR(128),
    _verif_code VARCHAR(10),
    _verif_code_expires VARCHAR(50))
RETURNS uuid
LANGUAGE plpgsql
AS
$$
    DECLARE
        _unique_user_id uuid;
    BEGIN
        SELECT gen_random_uuid() into _unique_user_id;

        INSERT INTO unregistered_users (
                           id,
                           login,
                           name,
                           email,
                           password,
                           is_verified,
                           verification_code,
                           verification_code_expires)
        VALUES(_unique_user_id,
               _login,
               _name,
               _email,
               _password,
               FALSE,
               _verif_code,
               _verif_code_expires);
        RETURN _unique_user_id;
    END;
$$;

CREATE OR REPLACE FUNCTION get_unregistered_user_by_id(_id uuid)
RETURNS table (
    __id UUID,
    __login VARCHAR(40),
    __name VARCHAR(40),
    __email VARCHAR(40),
    __password VARCHAR(128),
    __is_verified BOOLEAN,
    __verification_code VARCHAR(10),
    __verification_code_expires VARCHAR(50)
)
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
               password,
               is_verified,
               verification_code,
               verification_code_expires
           FROM unregistered_users
        WHERE id::text = _id::text;
    END;
$$;
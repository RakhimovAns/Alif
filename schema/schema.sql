drop table wallets,customers,roles,actions;
CREATE TABLE roles(
                      id bigserial primary key ,
                      name text not null
);
CREATE TABLE customers(
                          id text primary key ,
                          name text not null ,
                          login text not null unique ,
                          password text not null ,
                          role_id integer references roles(id) ,
                          created timestamp not null default current_timestamp
);

CREATE TABLE wallets(
                        id bigserial primary key ,
                        balance integer not null check ( balance>=0 ),
                        customer_id text not null references customers(id),
                        role_id integer not null references roles(id)

);
CREATE OR REPLACE FUNCTION check_role_balance()
    RETURNS TRIGGER AS
$BODY$
BEGIN
    IF NEW.role_id = 1 AND NEW.balance > 100000 THEN
        RAISE EXCEPTION 'Role 1 cannot have a balance greater than 100000.';
    ELSIF NEW.role_id = 2 AND NEW.balance > 10000 THEN
        RAISE EXCEPTION 'Role 2 cannot have a balance greater than 10000.';
    END IF;

    RETURN NEW;
END;
$BODY$
    LANGUAGE plpgsql;

CREATE TRIGGER role_balance_trigger
    BEFORE INSERT OR UPDATE ON wallets
    FOR EACH ROW
EXECUTE FUNCTION check_role_balance();


CREATE TABLE actions(
 wallet_id integer not null  unique,
 amount integer not null default 0,
 sum integer not null default 0
);
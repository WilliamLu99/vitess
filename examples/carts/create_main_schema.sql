CREATE TABLE carts_multicol
(
    region SMALLINT UNSIGNED NOT NULL,
    uid VARCHAR(36) NOT NULL, /* Maybe BINARY(16) in reality */
    data VARCHAR(128) NOT NULL, /* Maybe a BLOB in reality */
    PRIMARY KEY (uid)
);

CREATE TABLE carts_placement
(
    region VARCHAR(128) NOT NULL,
    uid VARCHAR(36) NOT NULL, /* Maybe BINARY(16) in reality */
    data VARCHAR(128) NOT NULL, /* Maybe a BLOB in reality */
    PRIMARY KEY (uid)
);

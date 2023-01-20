use main/-40;
select * from carts_multicol as `carts_multicol/-40`;
select * from carts_placement as `carts_placement/-40`;

use main/40-80;
select * from carts_multicol as `carts_multicol/40-80`;
select * from carts_placement as `carts_placement/40-80`;

use main/80-c0;
select * from carts_multicol as `carts_multicol/80-c0`;
select * from carts_placement as `carts_placement/80-c0`;

use main/c0-;
select * from carts_multicol as `carts_multicol/c0-`;
select * from carts_placement as `carts_placement/c0-`;

use main;
select * from carts_multicol;
select * from carts_placement;

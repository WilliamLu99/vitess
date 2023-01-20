* Initialize the cluster: `101_initial_cluster.sh`
* Source `env.sh` to set up some aliases
* Load data: `mysql --verbose main < insert_carts.sql`
* Show data in each shard: `mysql --verbose main < show_data.sql`
* Teardown: `201_teardown.sh`

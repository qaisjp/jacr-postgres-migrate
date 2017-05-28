.PHONY: schema

init_schema::
	createdb -U postgres jacr_dev
	cat schema.sql | psql -U postgres jacr_dev

schema.sql::
	pg_dump -s -U postgres jacr_dev > schema.sql



# save a copy of dev database into dev_backup
checkpoint:
	mkdir -p dev_backup
	pg_dump -F c -U postgres jacr_dev > dev_backup/$$(date +%F_%H-%M-%S).dump

# restore latest dev backup
restore_checkpoint::
	dropdb -U postgres jacr_dev
	createdb -U postgres jacr_dev
	pg_restore -U postgres -d jacr_dev $$(find dev_backup | grep \.dump | sort | tail -n 1)

#!/usr/bin/env bash
docker run -v $(pwd)/config/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://travelhack:123@localhost:7557/travelhack?sslmode=disable" up

package databases

import (
	Applications "job-portal-lite/domain/applications/entities"
	Jobs "job-portal-lite/domain/jobs/entities"
	User "job-portal-lite/domain/user/entities"
)

func Migrate() error {
	conn := Connect()
	return conn.AutoMigrate(
		&User.User{},
		&Jobs.Jobs{},
		&Applications.Applications{},
		&Applications.ApplicationHistories{},
	)
}

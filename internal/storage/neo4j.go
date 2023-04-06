package storage

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/umalmyha/go-neo4j/internal/model"
)

const emailsDatabase = "emails"

type Neo4jStorage struct {
	drv neo4j.DriverWithContext
}

func NewNeo4jStorage(drv neo4j.DriverWithContext) *Neo4jStorage {
	return &Neo4jStorage{
		drv: drv,
	}
}

func (s *Neo4jStorage) CreateUser(ctx context.Context, u *model.User) (err error) {
	sess := s.drv.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: emailsDatabase,
	})
	defer func() {
		if closeErr := sess.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	const q = `CREATE (u:User {id:$id, email:$email})`
	params := map[string]any{
		"id":    u.ID.String(),
		"email": u.Email,
	}

	_, err = sess.Run(ctx, q, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Neo4jStorage) CreateEmail(ctx context.Context, email *model.Email) (err error) {
	sess := s.drv.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: emailsDatabase,
	})
	defer func() {
		if closeErr := sess.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	recipients := make([]string, 0, len(email.To)+len(email.Cc))
	for _, to := range email.To {
		recipients = append(recipients, "(email)")
	}

	const q = `
		MATCH (from:User {email:$from})
		CREATE (email:Email {content:$content}),
		(from)-[:SENT]->(email)
	`
}

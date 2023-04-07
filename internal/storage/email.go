package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/umalmyha/go-neo4j/internal/model"
)

type Neo4jEmailStorage struct {
	drv neo4j.DriverWithContext
}

func NewNeo4jEmailStorage(drv neo4j.DriverWithContext) *Neo4jEmailStorage {
	return &Neo4jEmailStorage{drv: drv}
}

func (s *Neo4jEmailStorage) CreateUser(ctx context.Context, u *model.User) (err error) {
	sess := s.drv.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		if closeErr := sess.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	const q = `CREATE (u:User {id:$id, email:$email})`
	params := map[string]any{
		"id":    u.ID,
		"email": u.Email,
	}

	_, err = sess.Run(ctx, q, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Neo4jEmailStorage) FindUserByID(ctx context.Context, id string) (u *model.User, err error) {
	sess := s.drv.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() {
		if closeErr := sess.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	const q = `
		MATCH (u:User)
		WHERE u.id = $id
		RETURN u.id as id, u.email as email
		LIMIT 1
	`

	_, err = sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, q, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}

		for res.Next(ctx) {
			rec := res.Record()
			u = &model.User{
				ID:    rec.Values[0].(string),
				Email: rec.Values[1].(string),
			}
		}

		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Neo4jEmailStorage) CreateEmail(ctx context.Context, email *model.Email) (err error) {
	sess := s.drv.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		if closeErr := sess.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}()

	matches := make([]string, 0, len(email.Cc)+2)
	matches = append(matches, "(from:User {email:$from})", "(to:User {email:$to})")

	create := make([]string, 0, len(email.Cc)+3)
	create = append(create, "(email:Email {content:$content})", "(from)-[:SENT]->(email)", "(email)-[:TO]->(to)")

	params := map[string]any{
		"from":    email.From,
		"to":      email.To,
		"content": email.Content,
	}

	for i, cc := range email.Cc {
		ccIdent := fmt.Sprintf("cc_%d", i)
		ccEmail := fmt.Sprintf("cc_email_%d", i)

		matches = append(matches, fmt.Sprintf("(%s:User {email:$%s})", ccIdent, ccEmail))
		create = append(create, fmt.Sprintf("(email)-[:CC]->(%s)", ccIdent))

		params[ccEmail] = cc
	}

	q := fmt.Sprintf("MATCH %s CREATE %s", strings.Join(matches, ","), strings.Join(create, ","))

	if _, err := sess.Run(ctx, q, params); err != nil {
		return err
	}

	return nil
}

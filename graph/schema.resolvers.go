package graph

import (
	"context"
	"pastebin/graph/model"
)

// CreateDocument is the resolver for the createDocument field.
func (r *mutationResolver) CreateDocument(ctx context.Context, input model.CreateDocumentInput) (*int, error) {

	var id, err = r.DocumentService.CreateDocument(input.Value, input.Title, input.AccessKey, *input.MaxViewCount, *input.TTLMs)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// DeleteDocument is the resolver for the deleteDocument field.
func (r *mutationResolver) DeleteDocument(ctx context.Context, id int, accessKey string) (*bool, error) {
	success, err := r.DocumentService.DeleteDocument(id, accessKey)
	if err != nil {
		return nil, err
	}
	return &success, nil
}

// GetDocument is the resolver for the getDocument field.
func (r *queryResolver) GetDocument(ctx context.Context, id int) (*model.Document, error) {
	return r.DocumentService.GetDocument(id)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

package repository

import "gorm.io/gorm"

type RepositoryOpts struct {
	Preloads map[string][]string
	Omit     map[string][]string
	Limit    int64
	Offset   int64
}

type Repository interface {
	// Will find the object by id and set it to the calling object.
	FindByID(tx *gorm.DB, id int, opts *RepositoryOpts) error
	// Will return a slice of objects to out with the result.
	Query(tx *gorm.DB, query map[string]any, opts *RepositoryOpts, out any) error
	// Will create the object in the db, and set the id, created,updated, and deleted
	// times.
	New(tx *gorm.DB) error
	// Updates an object from map, this ensures we don't have any stale state
	// in our object and instead just updates the fields that exist in body.
	Update(tx *gorm.DB, id int, body map[string]any) error
	// Will use the source object (the one calling the function) to delete.
	Delete(tx *gorm.DB) error
	// Will delete the object by id.
	DeleteById(tx *gorm.DB, id int) error
	// Returns all objects for a given source type.
	All(tx *gorm.DB, opts *RepositoryOpts, out any) error
	// Find associations will query the associated objects and return them to out.
	// Example usage: Query Author(s) by name for a specefic post.
	FindAssociations(tx *gorm.DB, assoc string, query map[string]any, out AssociationEntity) error
	// Deletes the assocated item(s) in del.
	DeleteAssociation(tx *gorm.DB, assoc string, del ...AssociationEntity) error
	// Adds the item(s) associations for m2m, o2m, replaces current for has one,
	// and belongs to.
	AddAssociations(tx *gorm.DB, assoc string, inc ...AssociationEntity) error
	// Counts the associations of a specific type on that object. For example,
	// on Post, we could count the number of `Authors` associated with that post.
	AssociationCount(tx *gorm.DB, assoc string, out AssociationEntity) (int64, error)
	// Removes ALL of the specific associations of type `assoc`, does not delete the actual
	// assocated object though.
	// So if assoc = Authors it will remove all author associations with the object.
	ClearAssociations(tx *gorm.DB, assoc string) error
	Raw(tx *gorm.DB, raw string, queryParams []any, opts *RepositoryOpts) error
}

type PostRepositoryInterface interface {
	Repository
	// It just made sense to have this because at least the public facing
	// detail page for the client is going to be referencing the slug.
	FindBySlug(tx *gorm.DB, slug string, opts *RepositoryOpts) error
}

// More of a safety interface in an attempt to write "idomatic" code.
// This gives us some sort of idea what type of object we're going to be passing
// for associations.
// It also calls the ValidAssociation safety method that prevents the user from
// accidentally passing the wrong Association string and object combo.
type AssociationEntity interface {
	// ValidAssociation is defined on the associated object. For example,
	// if a post has a m2m with User called Authors. The User object
	// implements the AssociationEntity interface.
	ValidAssociation(srcType any, assoc string) bool
}

package graph

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Resolver struct {
	DB *mongo.Client
}

func (r *Resolver) Post() PostResolver {
	return PostResolver{r}
}

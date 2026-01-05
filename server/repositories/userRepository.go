package repositories

import (
	"context"
	models "server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Database
}

func (r *UserRepository) FindUserByEmail(email string) (*models.UserRegister, error) {
	collection := r.db.Collection("users")
	var user models.UserRegister
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (r *UserRepository) CreateUser(user models.UserRegister) error {
	collection := r.db.Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	return err
}

func (r *UserRepository) UpdateUser(user models.UserRegister) error {
	collection := r.db.Collection("users")
	_, err := collection.UpdateOne(context.Background(), bson.M{"email": user.Email}, bson.M{"$set": user})
	return err
}

func InitUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}
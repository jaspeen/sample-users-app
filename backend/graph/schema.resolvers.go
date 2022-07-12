package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"

	"github.com/jaspeen/sample-users-app/db"
	"github.com/jaspeen/sample-users-app/graph/generated"
	"github.com/jaspeen/sample-users-app/graph/model"
	"github.com/jaspeen/sample-users-app/services"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*model.LoginPayload, error) {
	var user db.User
	userRes := r.DB.First(&user, "email = ?", username)
	if userRes.RowsAffected == 0 {
		return nil, services.Err_UNAUTHENTICATED
	}
	if userRes.Error != nil {
		return nil, userRes.Error
	}

	if !services.CheckPassword(user.Password, password) {
		return nil, services.Err_UNAUTHENTICATED
	}

	var accessToken string
	var refreshToken string
	var err error

	err = r.DB.Transaction(func(tx *gorm.DB) error {
		accessToken, refreshToken, err = services.GenerateRefreshAndAccessToken(tx, &user)
		log.Printf("AT: %v, RT: %v\n", accessToken, refreshToken)
		return err
	})

	var gUser model.User
	copier.Copy(&gUser, user)
	gUser.ID = strconv.Itoa(user.ID)

	res := &model.LoginPayload{Token: accessToken, RefreshToken: refreshToken, User: &gUser}
	return res, nil
}

func (r *mutationResolver) RenewToken(ctx context.Context, refreshToken string) (*model.RefreshPayload, error) {
	accessToken, err := services.RenewAccessToken(r.DB, refreshToken)
	return &model.RefreshPayload{Token: accessToken, RefreshToken: refreshToken}, err
}

func (r *mutationResolver) AddUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	var user db.User
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if r.DB.Take(&user, "email = ?", input.Email).RowsAffected > 0 {
			// user already exist
			return services.Err_ALREADY_EXIST
		}
		err := copier.Copy(&user, input)
		if err != nil {
			return err
		}
		user.Password, err = services.HashPasword(input.Password)
		if err != nil {
			return err
		}

		return r.DB.Create(&user).Error
	})
	log.Debug().Msgf("User created: %v", user.ID)

	return &model.User{ID: strconv.Itoa(int(user.ID)), FirstName: input.FirstName, LastName: input.LastName, Gender: input.Gender, Email: input.Email, Phone: input.Phone}, err
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UserUpdate) (*model.User, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, services.Err_NOT_FOUND
	}
	var user db.User
	err = r.DB.Transaction(func(tx *gorm.DB) error {
		dbRes := tx.Take(&user, "id = ?", intId)
		if err := dbRes.Error; err != nil {
			return err
		}
		if dbRes.RowsAffected == 0 {
			return services.Err_NOT_FOUND
		}

		copier.CopyWithOption(&user, input, copier.Option{IgnoreEmpty: true})
		user.ID = intId

		return tx.Save(&user).Error
	})

	return &model.User{ID: id, FirstName: user.FirstName, LastName: user.LastName, Gender: DbGenderToModelGender(user.Gender), Email: user.Email, Phone: input.Phone}, err
}

func (r *mutationResolver) RemoveUser(ctx context.Context, id string) (*string, error) {
	return &id, r.DB.Where("id=?", id).Delete(&model.User{}).Error
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	var users []*db.User
	err := r.DB.Order("id").Find(&users).Error
	if err != nil {
		return nil, err
	}
	var res []*model.User
	err = copier.CopyWithOption(&res, users, copier.Option{Converters: []copier.TypeConverter{services.Int2StringConverter, services.String2IntConverter}})
	return res, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func DbGenderToModelGender(dbGen *db.Gender) *model.Gender {
	if dbGen != nil {
		res := model.Gender(string(*dbGen))
		return &res
	}
	return nil
}

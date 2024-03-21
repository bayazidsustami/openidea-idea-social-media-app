package app

import (
	"openidea-idea-social-media-app/controller"
	"openidea-idea-social-media-app/repository"
	"openidea-idea-social-media-app/security"
	"openidea-idea-social-media-app/service"
	"openidea-idea-social-media-app/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoute(app *fiber.App, dbPool *pgxpool.Pool) {

	validator := validator.New()
	validation.RegisterValidation(validator)

	authService := service.NewAuthService()

	userRepository := repository.New()
	userService := service.NewUserService(userRepository, validator, dbPool, authService)
	userController := controller.New(userService, authService)

	friendRepository := repository.NewFriendRepository(dbPool)
	friendService := service.NewFriendsService(validator, friendRepository)
	friendController := controller.NewFriendsController(friendService, authService)

	postRepository := repository.NewPostRepository(dbPool)
	postService := service.NewPostService(validator, postRepository)
	postController := controller.NewPostController(postService, authService)

	commentRepository := repository.NewCommentRepository(dbPool)
	commentService := service.NewCommentService(validator, commentRepository)
	commentController := controller.NewCommentController(commentService, authService)

	imageService := service.NewImageService(security.GetAws3Session())
	imageController := controller.NewImageUploadController(authService, imageService)

	user := app.Group("/v1/user")
	user.Post("/register", userController.Register)
	user.Post("/login", userController.Login)

	app.Use(security.CheckTokenHeaderExist)
	app.Use(security.GetJwtTokenHandler())

	app.Post("/v1/user/link/email", userController.UpdateEmail)
	app.Post("/v1/user/link/phone", userController.UpdatePhone)
	app.Patch("/v1/user", userController.UpdateAccount)

	app.Post("v1/friend", friendController.AddFriend)
	app.Delete("v1/friend", friendController.RemoveFriends)
	app.Get("/v1/friend", friendController.GetAllFriends)

	app.Post("/v1/post", postController.Create)
	app.Get("/v1/post", postController.GetAll)

	app.Post("/v1/post/comment", commentController.Create)

	app.Post("/v1/image", imageController.UploadImage)
}

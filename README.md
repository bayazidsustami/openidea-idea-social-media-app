# Social Media REST API

## Overview

This project is a RESTful API built in GoLang using Fiber v2 framework to simulate functionalities similar to popular social media platforms like Instagram or Twitter. It includes features such as adding/removing friends, creating posts, and commenting on posts. This project is task #2 of ProjectSprint that initiate by [@nandapagi](https://twitter.com/nandapagi). 

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/bayazidsustami openidea-idea-social-media-app.git
    cd openidea-idea-social-media-app
    ```

2. **Install dependencies:**

    ```bash
    go mod tidy
    ```

3. **Set up environment variables:**

    ```bash
    export DB_NAME=
    export DB_PORT=
    export DB_HOST=
    export DB_USERNAME=
    export DB_PASSWORD=
    export DB_PARAMS="sslmode=disabled"
    export JWT_SECRET=
    export BCRYPT_SALT=8 
    export S3_ID=
    export S3_SECRET_KEY=
    export S3_BUCKET_NAME=
    export S3_REGION=
    ```

4. **Database Setup:**

    This project use Postgresql, so ensure you have postgres.

5. **Run the application:**

    For developement, you can use `make run` to running the project. But if you want to build the binary you can use `make build` and go to `/buidl/project_name`. (Important!, please update build target as needs)

## Usage

Once the application is running, you can access the API endpoints using tools like cURL, Postman, or any HTTP client.

### Endpoints

- `POST /v1/user/register`: User Register.
- `POST /v1/user/login`: User Login.
- `POST /v1/user/link` : Link Email to Account.
    - Headers
        - `Authorization : Bearer {accessToken}`
    - Request Body:
        ```json
        {
            "email": "mai@mail.com"
        }
        ```
- `POST /v1/user/link/phone` : Link Phone to Account.
    - Headers
        - `Authorization : Bearer {accessToken}`
    - Request Body:
        ```json
        {
            "phone": "+612312131323"
        }
        ```
- `PATCH /v1/user` : Update Account.
    - Headers
        - `Authorization : Bearer {accessToken}`
    - Request Body:
        ```json
        {
            "imageUrl": "http://image.jpg",
            "name": "firstName lastName",
        }
        ```
- `GET /v1/friend`: Get User Friends.
    - Headers
        - `Authorization : Bearer {accessToken}`
    - Params (all optional)
        - `limit` & `offset` (number) default `limit=5&offset=0` display how much data in single request
        - `sortBy` (”`friendCount`”|”`createdAt`”) default `createdAt` , display the information based on defined value
        - `orderBy` (”`asc`”|”`desc`”) default `desc`
        - `onlyFriend` (`true`|`false`) default `false`, show only the user’s friend
        - `search` (string) , display information that contains the name of search
- `POST /v1/friend`: Add Friends
    - Headers
        - `Authorization : Bearer {accessToken}`
    - Request Body 
        ```json
        {
            "userId": "" //should be a valid user id
        }
        ```
- `DELETE /v1/friend`: Remove Friends
    - Headers
        - `Authorization : Bearer {accessToken}`
    - Request Body 
        ```json
        {
            "userId": "" //should be a valid user id
        }
        ```
- `POST /v1/post`: Create Post.
    - Headers
        - `Authorization : Bearer {accessToken}`
    - Request Body
        ```json
        {
            "postInHtml": "", 
            "tags": [""] 
        }
        ```
- `GET /v1/post`: Get Posts.
    - Headers
        - `Authorization : Bearer {accessToken}`
    - Params (all optional)
        - `limit` & `offset` (number) default `limit=5&offset=0` display how much data in single request
        - `search` (string) , display information that contains the post of search
        - `searchTag` ([]string) , search by tags
- `POST /v1/post/comment`: Create comment post.
    - Headers
        - `Authorization : Bearer {accessToken}`
    - Request Body
        ```json
        {
            "postInHtml": "", 
            "tags": [""] 
        }
        ```    
- `POST /v1/image`: Upload Image.
    - Headers
        - `Authorization : Bearer {accessToken}`
    - Request 
        - `Content-Type` : `multipart/form-data`
        - `file` : `image.jpeg` | `image.jpg`

**Note:** Make sure to include appropriate authorization headers when accessing protected endpoints.

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

- [Fiber](https://github.com/gofiber/fiber): Fast HTTP framework for Go.
- [Golang Validation](https://github.com/go-playground/validator): Go Struct and Field validation.
- [PGX](https://github.com/jackc/pgx): PostgreSQL driver and toolkit for Go.
- [Viper](https://github.com/spf13/viper): Go configuration with fangs!
- [Goccy JSON](https://github.com/goccy/go-json): JSON parser for Go.

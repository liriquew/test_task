import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";


using Http;
using TypeSpec.OpenAPI;
@doc("""
  Simple User service. Basic auth is used for authorization.
  - all /users/{userId} require admin permissions.
  - endpoint GET /users/ can be used by all users
  """)
@service(#{
  title: "user_service"
})
@server("http://localhost:8080", "server")
namespace user_service;

@doc("uuid defenition required in user model")
@format("uuid")
scalar uuid extends string;

/* models */
@doc("""
  User model all fields isn't required
    - `id`: the uuid
    - `username`: the user's name
    - `password`: the user's password, returned like a base64 string
    - `email`: the user's email
    - `is_admin`: define user permissions
  """)
@example(AdminUser, #{title: "1", description: "2"})
model User {
  @extension("x-oapi-codegen-extra-tags", #{db: "id"})
  id?: uuid;
  
  @extension("x-oapi-codegen-extra-tags", #{db: "username"})
  username?: string;

  @extension("x-oapi-codegen-extra-tags", #{db: "password"})
  password?: string;

  @extension("x-oapi-codegen-extra-tags", #{db: "email"})
  email?: string;

  @extension("x-oapi-codegen-extra-tags", #{db: "is_admin"})
  is_admin?: boolean;
}

/* Response models */
model UserResponse {
  ...OkResponse;
  ...Body<User>;
}

model UserCreatedResponse {
  ...CreatedResponse;
  ...Body<User>;
}

model UserListResponse {
  ...OkResponse;
  ...Body<User[]>;
} 


model ValidationErrorResponse {
  ...ValidationError
}

model NotFoundResponse {
  ...NotFoundError
}

model AlreadyExistsResponse {
  ...AlreadyExistsError
}

model ForbiddenResponse {
  ...ForbiddenError
}

model InternalErrorResponse {
  ...InternalServerError;
}

/* handlers */
@route("/health")
@get
op health(): OkResponse;

@route("/users/")
namespace Service {
  @tag("Users")
  @doc("Returns a list of all users")
  @get
  @useAuth(BasicAuth)
  op listUsers(): 
    | UserListResponse 
    | InternalErrorResponse;

  @tag("Users")
  @doc("Returns a User if user with provided userId exists, 404 otherwise")
  @get
  @useAuth(BasicAuth)
  op getUser(@path userId: uuid): 
    | UserResponse 
    | ValidationErrorResponse
    | ForbiddenResponse
    | NotFoundResponse
    | InternalErrorResponse;

  @tag("Users")
  @doc("""
    Create a user 
    - all fields must be provided, 400 otherwise
    - admin permission required
  """)
  @post
  @useAuth(BasicAuth)
  op createUser(@body user: User):
    | UserCreatedResponse
    | ValidationErrorResponse
    | ForbiddenResponse 
    | AlreadyExistsResponse 
    | InternalErrorResponse;

  @tag("Users")
  @doc("""
    Patch User
    - one of the fields must be provided, except id
    - admin permission required 
  """)
  @patch
  @useAuth(BasicAuth)
  op patchUser(@path userId: uuid, @body user: User):
    | OkResponse
    | ValidationErrorResponse
    | ForbiddenResponse
    | AlreadyExistsResponse
    | InternalErrorResponse;

  @tag("Users")
  @doc("""
    Put a new User params 
    - all fields must be provided, except id
    - admin permission required
  """)
  @put
  @useAuth(BasicAuth)
  op putUser(@path userId: uuid, @body user: User): 
    | OkResponse
    | ValidationErrorResponse
    | ForbiddenResponse
    | AlreadyExistsResponse
    | InternalErrorResponse;
    
  @tag("Users")
  @doc("""
    Delete User
    - admin permission required 
  """)
  @delete
  @useAuth(BasicAuth)
  op deleteUser(@path userId: uuid): 
    | OkResponse
    | ValidationErrorResponse
    | ForbiddenResponse
    | InternalErrorResponse;
}

/* errors */
@error
model InternalServerError {
  @statusCode code: 500;
  message: "internal server error";
}

@error
model ValidationError {
  @statusCode code: 400;
  message: "bad params";
}

@error
model AlreadyExistsError {
  @statusCode code: 409;
  message: "already exists, username taken" | "already exists, email taken";
}

@error
model NotFoundError {
  @statusCode code: 404;
  message: "not found";
}

@error
model ForbiddenError {
  @statusCode code: 403;
  message: "forbidden, admin permission required";
}

/* example */
const AdminUser: User = #{
  id: "ac63a680-bddb-4102-b7a3-9fdc6ee53df2",
  username: "admin", 
  password: "admin",
  email: "admin@admin.ru",
  is_admin: true,
};
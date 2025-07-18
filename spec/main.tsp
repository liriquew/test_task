import "@typespec/http";
import "@typespec/openapi";


using Http;
using TypeSpec.OpenAPI;
@service(#{
  title: "user_service"
})
namespace user_service;

@format("uuid")
@minLength(36)
@maxLength(36)
scalar uuid extends string;

enum OptBool {
  admin: "true",
  user: "false",
  unSet: "unset",
}

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

@route("/ping")
@get
op health(): {
  @statusCode statusCode: 200;
};

@route("/users/")
namespace Service {

  @get
  @useAuth(BasicAuth)
  op listUsers(): {
    @statusCode statusCode: 200;
    @body users: User[];
  } | {
    @statusCode statusCode: 500;
    @body error: InternalServerError;
  };

  @get
  @useAuth(BasicAuth)
  op getUser(@path userId: uuid): {
    @statusCode statusCode: 200;
    @body user: User;
  } | {
    @statusCode statusCode: 400;
    @body error: ValidationError;
  } | {
    @statusCode statusCode: 403;
    @body error: Forbidden;
  } | {
    @statusCode statusCode: 404;
    @body error: NotFound;
  } | {
    @statusCode statusCode: 500;
    @body error: InternalServerError;
  };

  @post
  @useAuth(BasicAuth)
  op createUser(@body user: User): {
    @statusCode statusCode: 201;
    @body newUser: User;
  } | {
    @statusCode statusCode: 400;
    @body error: ValidationError;
  } | {
    @statusCode statusCode: 403;
    @body error: Forbidden;
  } | {
    @statusCode statusCode: 409;
    @body error: AlreadyExists;
  } | {
    @statusCode statusCode: 500;
    @body error: InternalServerError;
  };

  @patch
  @useAuth(BasicAuth)
  op patchUser(@path userId: uuid, @body user: User): {
    @statusCode statusCode: 200;
  } | {
    @statusCode statusCode: 400;
    @body error: ValidationError;
  } | {
    @statusCode statusCode: 403;
    @body error: Forbidden;
  } | {
    @statusCode statusCode: 409;
    @body error: AlreadyExists;
  } | {
    @statusCode statusCode: 500;
    @body error: InternalServerError
  };

  @put
  @useAuth(BasicAuth)
  op putUser(@path userId: uuid, @body user: User): {
    @statusCode statusCode: 200;
  } | {
    @statusCode statusCode: 400;
    @body error: ValidationError;  
  } | {
    @statusCode statusCode: 403;
    @body error: Forbidden;
  } | {
    @statusCode statusCode: 409;
    @body error: AlreadyExists;
  } | {
    @statusCode statusCode: 500;
    @body error: InternalServerError;
  };

  @delete
  @useAuth(BasicAuth)
  op deleteUser(@path userId: uuid): {
    @statusCode statusCode: 200;
  } | {
    @statusCode statusCode: 400;
    @body error: ValidationError;
  } | {
    @statusCode statusCode: 403;
    @body error: Forbidden;
  } | {
    @statusCode statusCode: 500;
    @body error: InternalServerError;
  };
}

@error
model InternalServerError {
  code: 500;
  message: "internal server error";
}

@error
model ValidationError {
  code: 400;
  message: "bad params";
  details: string[];
}

@error
model AlreadyExists {
  code: 409;
  message: "already exists";
  details: string[];
}

@error
model NotFound {
  code: 404;
  message: "not found";
}

@error
model Forbidden {
  code: 403;
  message: "forbidden, admin permission required";
}



// model Widget {
//   id: string;
//   weight: int32;
//   color: "red" | "blue";
// }

// model WidgetList {
//   items: Widget[];
// }

// @error
// model Error {
//   code: int32;
//   message: string;
// }

// model AnalyzeResult {
//   id: string;
//   analysis: string;
// }

// @route("/widgets")
// @tag("Widgets")
// interface Widgets {
//   /** List widgets */
//   @get list(): WidgetList | Error;
//   /** Read widgets */
//   @get read(@path id: string): Widget | Error;
//   /** Create a widget */
//   @post create(@body body: Widget): Widget | Error;
//   /** Update a widget */
//   @patch update(@path id: string, @body body: MergePatchUpdate<Widget>): Widget | Error;
//   /** Delete a widget */
//   @delete delete(@path id: string): void | Error;

//   /** Analyze a widget */
//   @route("{id}/analyze") @post analyze(@path id: string): AnalyzeResult | Error;
// }

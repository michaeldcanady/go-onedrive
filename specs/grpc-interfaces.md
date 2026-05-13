# gRPC Interface Specification

## Storage Service
The `StorageService` defines the standard contract for storage backends to interact with the core filesystem. All requests contain an `options` map for injecting session-specific data like authentication tokens and target drive identifiers.

### RPC Methods
- `List(ListRequest) -> ListResponse`: Retrieves a list of nodes at a given path.
- `Stat(StatRequest) -> StatResponse`: Retrieves metadata for a single node.
- `Mkdir(MkdirRequest) -> MkdirResponse`: Creates a new directory.
- `Read(ReadRequest) -> stream ReadResponse`: Reads the contents of a file as a stream of data chunks.
- `Write(stream WriteRequest) -> WriteResponse`: Writes content to a specified path using a stream of data chunks. The first message must contain the path and options.
- `Delete(DeleteRequest) -> DeleteResponse`: Deletes a node at a specified path.
- `Move(MoveRequest) -> MoveResponse`: Moves or renames a node within the backend.
- `ListDrives(ListDrivesRequest) -> ListDrivesResponse`: Discovers available storage containers (drives/libraries) for the provided identity.
- `GetDrive(GetDriveRequest) -> GetDriveResponse`: Retrieves details for a specific drive by ID.

### Key Data Structures
- `Node`: Represents a file or directory with ID, name, path, type, size, and modification time.
- `NodeType`: Enum for `FILE` or `DIRECTORY`.
- `Drive`: Represents a storage container with ID, name, and type (e.g., personal, business).

---

## Identity Plugin
The `IdentityPlugin` manages the interaction with specific authentication providers (e.g., Azure AD, Google). It is strictly **stateless**; the plugin is responsible for performing authentication and refresh flows, returning results to the host for persistence.

### RPC Methods
- `Login(LoginRequest) -> LoginResponse`: Initiates the authentication flow (Interactive, Device Code, etc.). Returns the `AccessToken` and the `Identity` metadata.
- `Refresh(RefreshRequest) -> RefreshResponse`: Performs a token refresh using a provided `refresh_token`. Returns a new `AccessToken`.
- `ListIdentities(ListIdentitiesRequest) -> ListIdentitiesResponse`: Returns identities known to the local provider environment.
- `Logout(LogoutRequest) -> LogoutResponse`: (Optional) Informs the provider to invalidate the remote session.

### Key Data Structures
- `AccessToken`: Contains the raw access token, refresh token, expiry timestamp, and granted scopes.
- `Identity`: Represents a user with a unique ID, display name, email, and provider name.
- `LoginRequest`: Contains an `options` map for parameters like `client_id`, `tenant_id`, `redirect_uri`, and `method`.

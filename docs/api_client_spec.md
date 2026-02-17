# Workbrew API Client Specification

## API Characteristics & Design Implications

### Authentication
- **Pattern**: Static Bearer token in HTTP header
- **Header**: `Authorization: Bearer {API_KEY}`
- **Token Lifetime**: No expiration
- **Token Refresh**: Not required
- **Impact on Client**:
  - Simple authentication: API key set once during client initialization via `SetAuthScheme("Bearer")` and `SetAuthToken()`
  - No `TokenManager` needed
  - No token refresh logic or middleware
  - Thread-safe by nature (read-only static value)
  - Contrast with Nexthink: No OAuth2 flow, no separate auth endpoint
  - Contrast with VirusTotal: Uses `Authorization: Bearer` instead of `x-apikey` header

### API Versioning
- **Pattern**: Version specified via HTTP header
- **Header**: `X-Workbrew-API-Version: {version}`
- **Default Version**: `v0`
- **Impact on Client**:
  - API version set as global header during client initialization
  - No version in URL path (unlike many REST APIs)
  - Version header included in all requests
  - Allows server-side version negotiation

### Workspace Scoping
- **Pattern**: All API endpoints are workspace-scoped
- **URL Structure**: `https://console.workbrew.com/workspaces/{workspace_name}/{resource}.{format}`
- **Workspace**: Customer-specific workspace identifier
- **Impact on Client**:
  - Workspace name stored in client configuration
  - Base URL construction: `BaseURL + "/workspaces/" + workspace`
  - Workspace injected into all endpoint paths
  - Validation ensures workspace is provided during initialization
  - Contrast with Nexthink: Nexthink uses instance+region, VirusTotal has fixed base URL

### Response Formats
- **Supported**: JSON, CSV
- **Format Selection**: Via URL extension (`.json`, `.csv`)
- **Examples**:
  - JSON: `/workspaces/{workspace}/devices.json`
  - CSV: `/workspaces/{workspace}/devices.csv`
  
- **Impact on Client**:
  - Separate methods per format (`ListDevices`, `ListDevicesCSV`)
  - CSV methods return `[]byte` raw data
  - JSON methods return typed structs
  - CSV parsing left to caller
  - Contrast with Nexthink: Nexthink uses `Accept` header for format negotiation

### Pagination
- **Pattern**: None
- **Data Delivery**: Full datasets returned in single response
- **Impact on Client**:
  - No pagination helpers needed
  - No cursor/offset parameters
  - Single API call returns complete list
  - Contrast with VirusTotal: VirusTotal uses cursor-based pagination
  - Contrast with Nexthink: Nexthink uses async export for large datasets

### Resource Model
- **Pattern**: Simple REST with list-oriented operations
- **Resources**: Devices, Device Groups, Formulae, Casks, Brew Taps, Brew Configurations, Brew Commands, Brewfiles, Licenses, Events, Vulnerabilities, Analytics
- **Operations**: Primarily read-only (List, Get)
- **Impact on Client**:
  - Service-per-resource structure
  - Mostly GET operations
  - No complex relationships or navigation
  - Contrast with VirusTotal: VirusTotal has relationship navigation between resources

### Query Capabilities
- **Pattern**: None
- **Filtering**: No query parameters for filtering/searching
- **Impact on Client**:
  - No query builder needed
  - No search functionality
  - Client-side filtering if needed
  - Contrast with Nexthink: Nexthink uses pre-configured NQL queries
  - Contrast with VirusTotal: VirusTotal has advanced search with modifiers

### Device Identification
- **Format**: Serial numbers
- **Fields**: Serial number, groups, MDM name, last seen, device type, OS version, Homebrew/Workbrew versions, package counts
- **Impact on Client**:
  - Simple string identifiers (no UUID validation required)
  - No multiple ID type support
  - Contrast with Nexthink: Nexthink has internal IDs (Collector IDs, SIDs) vs external IDs

### Rate Limiting
- **Pattern**: Standard HTTP 429 (assumed, not explicitly documented)
- **Impact on Client**:
  - Standard retry logic with exponential backoff
  - Response wrapper exposes headers for rate limit info

### Async Operations
- **Pattern**: None
- **All Operations**: Synchronous
- **Impact on Client**:
  - No polling mechanisms needed
  - All responses immediate
  - Contrast with Nexthink: Nexthink has async NQL exports
  - Contrast with VirusTotal: VirusTotal has async file analysis

### Validation Requirements
- **Minimal**: No complex format validation (UUIDs, SIDs, etc.)
- **Required Fields**: Workspace name, API key
- **Impact on Client**:
  - Simple validation (non-empty strings)
  - No regex validators needed
  - Contrast with Nexthink: Nexthink requires extensive format validation

### Data Types
- **Timestamps**: Last seen, command last run times
- **Versions**: Homebrew version, Workbrew version, OS version strings
- **Counts**: Installed formulae count, casks count
- **Groups**: Array of group names
- **Impact on Client**:
  - Standard Go types (time.Time, string, int)
  - No special type handling required

### Error Handling
- **Structured Errors**: JSON error responses (assumed standard REST)
- **Status Codes**: Standard HTTP status codes
  - 400: Invalid parameters
  - 401: Invalid API key
  - 403: Permission denied
  - 404: Resource not found
  - 429: Rate limit exceeded
  
- **Impact on Client**:
  - Standard error types
  - Response object returned even on error
  - Simple error messages

### Base URL
- **Pattern**: `https://console.workbrew.com`
- **Fixed**: No dynamic construction (similar to VirusTotal)
- **Workspace Appended**: `/workspaces/{workspace_name}` added to all paths
- **Impact on Client**:
  - Constant base URL
  - Workspace configuration required at initialization
  - Optional override for testing/proxies

### CSV Format Details
- **Headers**: Yes (first row contains column names)
- **Columns**: serial_number, groups, mdm_user_or_device_name, last_seen_at, command_last_run_at, device_type, os_version, homebrew_prefix, homebrew_version, workbrew_version, formulae_count, casks_count
- **Delimiter**: Comma
- **Impact on Client**:
  - Raw byte stream returned
  - No CSV parsing in SDK
  - Caller responsible for parsing
  - Use standard library `encoding/csv` or third-party parser

### API Philosophy
- **Design**: Focused, opinionated REST API for Homebrew management
- **Scope**: Read-only operations for monitoring and reporting
- **Simplicity**: No complex queries, relationships, or async operations
- **Impact on Client**:
  - Lightweight client implementation
  - Minimal configuration required
  - Simple, predictable API calls
  - Contrast: Much simpler than Nexthink or VirusTotal clients

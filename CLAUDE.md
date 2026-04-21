# Clipper

Clipper is a cross-platform clipboard synchronization application with user authentication.

## Project Overview

- **Purpose**: Sync clipboard content (text, images) across multiple devices
- **Architecture**: Go backend + Vue 3 frontend
- **Key Features**: User auth (Session Cookie), SQLite storage, WebSocket real-time sync, device naming, image support

## Project Structure

```
Clipper/
├── cmd/server/main.go          # Entry point + CORS middleware + static file routes
├── config.yaml                 # Configuration
├── clipper.exe                 # Compiled executable
├── internal/
│   ├── config/config.go        # Viper config loader
│   ├── model/                  # Data models (user, clip, device)
│   ├── handler/                # HTTP handlers (auth, clip, device)
│   ├── service/                # Business logic (auth, clip, device)
│   ├── repository/             # SQLite data access (db, user, clip, device)
│   ├── ws/                     # WebSocket (hub, client, message)
│   ├── middleware/auth.go      # Session authentication (SameSite=Lax)
│   ├── pkg/crypto/password.go  # Password hashing (bcrypt)
│   └── pkg/storage/storage.go  # Image storage & thumbnail generation
├── uploads/                    # Uploaded images (user-isolated directories)
│   └── {user_id}/              # Per-user directory
│       ├── {clip_id}_orig.{ext}    # Original image
│       └── {clip_id}_thumb.jpg    # Thumbnail (300px width)
├── web/                        # Vue 3 frontend
│   ├── dist/                   # Build output
│   │   └── favicon.svg         # Browser tab icon
│   ├── public/                 # Static assets (copied to dist)
│   │   └── favicon.svg         # Clipper SVG icon
│   └── src/
│       ├── views/              # Login, Register, Dashboard, Admin
│       ├── stores/             # Pinia stores (auth, clip, device)
│       ├── composables/        # useWebSocket (auto reconnect with device name)
│       ├── types/index.ts      # TypeScript types
│       └── router/index.ts     # Vue Router with auth guards
└── data/                       # Database directory
```

## Implementation Status

**Completed:**
- [x] Full backend implementation (Repository/Service/Handler/WebSocket/Middleware)
- [x] Frontend auth flow (register/login/logout + state persistence)
- [x] Text clip creation, listing, deletion
- [x] WebSocket real-time sync (user-isolated broadcast)
- [x] Session Cookie auth (SameSite=Lax)
- [x] iOS Safari clipboard compatibility (execCommand)
- [x] Register page password confirmation validation
- [x] Password show/hide toggle (eye icon)
- [x] Page layout optimization (gradient backgrounds, rounded cards)
- [x] Dashboard two-column layout (create left / list right)
- [x] Tab switching animation (horizontal slide + elastic curve)
- [x] Pagination (5/10/20 items)
- [x] Right side list independent scrolling (fixed max height)
- [x] Image clipboard upload, preview, download functionality
- [x] Image thumbnail generation (300px width, JPEG 85% quality)
- [x] User-isolated image storage (`uploads/{user_id}/`)
- [x] 20MB image size limit, format validation (JPEG/PNG/GIF/WebP)
- [x] Device naming (login page + dashboard editable)
- [x] Device online/offline notifications with device name
- [x] Username duplicate validation (toast + red input)
- [x] Removed email field (username-only registration)
- [x] Mobile responsive layout
- [x] Image preview dialog (copy/download/close buttons)
- [x] Clip card layout fix (no compression, list container scroll)
- [x] Browser favicon (Clipper SVG icon)
- [x] Global Ctrl+V paste detection (auto-detect text/image, auto-switch tab)
- [x] Session error fix (invalid old cookie doesn't block login/logout)

**Pending:**
- [ ] Clip search/filter
- [ ] Device management page (online device list)
- [ ] Encrypted clip storage

## Quick Start

```bash
# Production mode (recommended)
cd web && npm run build
go run ./cmd/server/main.go

# Access at http://localhost:8080
```

## API Endpoints

### Auth
- `POST /api/v1/auth/register` - Register user (username + password only)
- `POST /api/v1/auth/login` - Login (sets session cookie)
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/me` - Get current user

### Clips
- `GET /api/v1/clips` - List user's clips
- `POST /api/v1/clips?device_id=xxx` - Create text clip
- `POST /api/v1/clips/upload?device_id=xxx` - Upload image (multipart/form-data)
- `GET /api/v1/clips/:id` - Get clip
- `DELETE /api/v1/clips/:id?device_id=xxx` - Delete clip (also deletes image files)

### Devices
- `GET /api/v1/devices` - List user's devices
- `GET /api/v1/devices/ws?device_id=xxx&device_name=xxx` - WebSocket (requires auth)

### Static Files
- `GET /uploads/{user_id}/{clip_id}_orig.{ext}` - Original image
- `GET /uploads/{user_id}/{clip_id}_thumb.jpg` - Thumbnail image

## WebSocket Message Types

| Type | Description |
|------|-------------|
| `clip_created` | New clip created notification (includes clip metadata) |
| `clip_deleted` | Clip deleted notification |
| `device_online` | Device connected (shows device_name) |
| `device_offline` | Device disconnected (shows device_name) |

## Fixed Issues

1. Login state lost on refresh - Router guard now waits for auth check
2. Login/logout not redirecting - Reset auth check flag before navigation
3. iOS clipboard copy fails - execCommand + textarea fallback
4. Cross-origin cookie not sent - Dynamic CORS origin + SameSite=Lax
5. Tab switching visual leak - overflow:hidden + horizontal slide animation
6. WebSocket shows "Unknown Device" - Now passes device_name in URL params
7. Clip card content truncated - List container scroll, cards not compressed
8. Preview dialog white plate - Global style override with transparent dialog background
9. No browser favicon - Added SVG favicon in public/, served via main.go
10. Session error blocks login - Invalid old cookie now creates new session instead of returning error

## Database Schema

```sql
-- Users (no email field)
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME,
    updated_at DATETIME
);

-- Clips
CREATE TABLE clips (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    device_id TEXT NOT NULL,
    type TEXT NOT NULL,        -- 'text' or 'image'
    content TEXT NOT NULL,     -- text content or image path
    meta TEXT,                 -- JSON: {width, height, size, format, thumb_path}
    created_at DATETIME,
    updated_at DATETIME
);

-- Devices
CREATE TABLE devices (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    ip TEXT,
    last_seen DATETIME,
    is_online BOOLEAN
);
```

## Notes for New Session

1. All core components are implemented and tested
2. Use `http://localhost:8080` (production build) for best experience
3. SQLite database auto-migrates on first run
4. WebSocket broadcasts are user-isolated
5. Image uploads stored in `uploads/{user_id}/` with thumbnails
6. Device names stored in localStorage and passed via WebSocket URL
7. Next priority: clip search/filter or device management page
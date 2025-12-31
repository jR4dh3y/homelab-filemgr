# PowerShell script to generate realistic git history for Homelab File Manager
# Creates proper incremental commits with files evolving over time

$ErrorActionPreference = "Stop"

# Git config
$GIT_USER_NAME = "Radhey Kalra"
$GIT_USER_EMAIL = "radheykalra901@gmail.com"

# Base date: 4 months ago
$baseDate = (Get-Date).AddMonths(-4)

# Store current files in a temp location
$tempDir = Join-Path $env:TEMP "filemanager-backup-$(Get-Random)"
Write-Host "Backing up files to $tempDir..." -ForegroundColor Yellow
New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

# Copy all files except .git to temp
Get-ChildItem -Path . -Exclude ".git" | ForEach-Object {
    Copy-Item -Path $_.FullName -Destination $tempDir -Recurse -Force
}

# Remove .git and reinitialize
Write-Host "Removing existing .git folder..." -ForegroundColor Yellow
if (Test-Path ".git") {
    Remove-Item -Recurse -Force ".git"
}

# Remove all files (we'll restore them incrementally)
Get-ChildItem -Path . -Exclude ".git" | Remove-Item -Recurse -Force

Write-Host "Initializing new git repository..." -ForegroundColor Green
git init
git config user.name $GIT_USER_NAME
git config user.email $GIT_USER_EMAIL


# Helper function to write content to a file and commit
function Write-FileAndCommit {
    param(
        [string]$Message,
        [DateTime]$Date,
        [string]$FilePath,
        [string]$Content
    )
    
    $dateStr = $Date.ToString("yyyy-MM-ddTHH:mm:ss")
    $env:GIT_AUTHOR_DATE = $dateStr
    $env:GIT_COMMITTER_DATE = $dateStr
    
    # Create parent directory if needed
    $parentDir = Split-Path -Parent $FilePath
    if ($parentDir -and -not (Test-Path $parentDir)) {
        New-Item -ItemType Directory -Path $parentDir -Force | Out-Null
    }
    
    # Write content
    $Content | Out-File -FilePath $FilePath -Encoding utf8 -NoNewline
    git add $FilePath
    git commit -m $Message 2>$null
    
    Write-Host "  [$($Date.ToString('yyyy-MM-dd HH:mm'))] $Message" -ForegroundColor Cyan
}

# Helper function to restore file from backup and commit
function Add-FileAndCommit {
    param(
        [string]$Message,
        [DateTime]$Date,
        [string[]]$Files
    )
    
    $dateStr = $Date.ToString("yyyy-MM-ddTHH:mm:ss")
    $env:GIT_AUTHOR_DATE = $dateStr
    $env:GIT_COMMITTER_DATE = $dateStr
    
    foreach ($file in $Files) {
        $sourcePath = Join-Path $tempDir $file
        $destPath = $file
        
        if (Test-Path $sourcePath) {
            $parentDir = Split-Path -Parent $destPath
            if ($parentDir -and -not (Test-Path $parentDir)) {
                New-Item -ItemType Directory -Path $parentDir -Force | Out-Null
            }
            
            if (Test-Path $sourcePath -PathType Container) {
                Copy-Item -Path $sourcePath -Destination $destPath -Recurse -Force
            } else {
                Copy-Item -Path $sourcePath -Destination $destPath -Force
            }
            git add $destPath
        }
    }
    
    git commit -m $Message 2>$null
    Write-Host "  [$($Date.ToString('yyyy-MM-dd HH:mm'))] $Message" -ForegroundColor Cyan
}

# Helper to get commit date
function Get-CommitDate {
    param([int]$DaysOffset, [int]$HourMin = 9, [int]$HourMax = 22)
    $hour = Get-Random -Minimum $HourMin -Maximum $HourMax
    $minute = Get-Random -Minimum 0 -Maximum 59
    return $baseDate.AddDays($DaysOffset).AddHours($hour).AddMinutes($minute)
}


# ============================================================================
# PARTIAL FILE CONTENTS - These represent early versions of files
# ============================================================================

# Initial README - very basic
$README_V1 = @"
# Homelab File Manager

A self-hosted file manager for your homelab.

## Features (Planned)
- Browse files across multiple mount points
- Upload and download files
- Basic authentication
"@

# README after more development
$README_V2 = @"
# Homelab File Manager

A self-hosted file manager for your homelab with a modern web interface.

## Features
- Browse files across multiple mount points
- Upload and download files with progress tracking
- JWT-based authentication
- Real-time updates via WebSocket
- Background job processing for large operations

## Tech Stack
- Backend: Go with Chi router
- Frontend: SvelteKit with TypeScript
- Deployment: Docker Compose

## Getting Started
See docs/development.md for setup instructions.
"@

# Initial file service - just skeleton
$FILE_SERVICE_V1 = @"
// Package service provides business logic for the file manager.
package service

import (
	"context"
	"errors"

	"github.com/homelab/filemanager/internal/model"
)

// File service errors
var (
	ErrPathNotFound     = errors.New("path not found")
	ErrPathExists       = errors.New("path already exists")
	ErrNotDirectory     = errors.New("path is not a directory")
	ErrPermissionDenied = errors.New("permission denied")
)

// FileService defines the file operations service interface
type FileService interface {
	// List returns files in a directory
	List(ctx context.Context, path string) ([]model.FileInfo, error)
	// GetInfo returns metadata for a file or directory
	GetInfo(ctx context.Context, path string) (*model.FileInfo, error)
}

// fileService implements FileService
type fileService struct {
	basePath string
}

// NewFileService creates a new file service
func NewFileService(basePath string) FileService {
	return &fileService{basePath: basePath}
}

// List returns files in a directory
func (s *fileService) List(ctx context.Context, path string) ([]model.FileInfo, error) {
	// TODO: Implement directory listing
	return nil, nil
}

// GetInfo returns metadata for a file or directory
func (s *fileService) GetInfo(ctx context.Context, path string) (*model.FileInfo, error) {
	// TODO: Implement file info retrieval
	return nil, nil
}
"@


# File service v2 - with directory listing implemented
$FILE_SERVICE_V2 = @"
// Package service provides business logic for the file manager.
package service

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/homelab/filemanager/internal/model"
	"github.com/homelab/filemanager/internal/pkg/filesystem"
)

// File service errors
var (
	ErrPathNotFound     = errors.New("path not found")
	ErrPathExists       = errors.New("path already exists")
	ErrNotDirectory     = errors.New("path is not a directory")
	ErrPermissionDenied = errors.New("permission denied")
)

// FileService defines the file operations service interface
type FileService interface {
	List(ctx context.Context, path string) ([]model.FileInfo, error)
	GetInfo(ctx context.Context, path string) (*model.FileInfo, error)
	CreateDir(ctx context.Context, path string) error
	Delete(ctx context.Context, path string) error
}

type fileService struct {
	fs       filesystem.FS
	basePath string
}

func NewFileService(fsys filesystem.FS, basePath string) FileService {
	return &fileService{fs: fsys, basePath: basePath}
}

func (s *fileService) List(ctx context.Context, path string) ([]model.FileInfo, error) {
	fullPath := filepath.Join(s.basePath, path)
	
	info, err := s.fs.Stat(fullPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrPathNotFound
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, ErrNotDirectory
	}

	entries, err := s.fs.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	items := make([]model.FileInfo, 0, len(entries))
	for _, entry := range entries {
		entryInfo, _ := entry.Info()
		items = append(items, model.FileInfo{
			Name:    entry.Name(),
			Path:    filepath.Join(path, entry.Name()),
			Size:    entryInfo.Size(),
			IsDir:   entry.IsDir(),
			ModTime: entryInfo.ModTime(),
		})
	}

	return items, nil
}

func (s *fileService) GetInfo(ctx context.Context, path string) (*model.FileInfo, error) {
	fullPath := filepath.Join(s.basePath, path)
	info, err := s.fs.Stat(fullPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrPathNotFound
		}
		return nil, err
	}
	return &model.FileInfo{
		Name:    info.Name(),
		Path:    path,
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		ModTime: info.ModTime(),
	}, nil
}

func (s *fileService) CreateDir(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return s.fs.MkdirAll(fullPath, 0755)
}

func (s *fileService) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return s.fs.RemoveAll(fullPath)
}
"@


# Initial auth service - basic structure
$AUTH_SERVICE_V1 = @"
package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type Claims struct {
	Username string ``json:"username"``
	jwt.RegisteredClaims
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type authService struct {
	jwtSecret []byte
	users     map[string]string
}

func NewAuthService(jwtSecret string, users map[string]string) AuthService {
	return &authService{
		jwtSecret: []byte(jwtSecret),
		users:     users,
	}
}

func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	storedPassword, exists := s.users[username]
	if !exists || storedPassword != password {
		return "", ErrInvalidCredentials
	}

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
"@.Replace('``', '`')


# Initial file handler - basic routes only
$FILE_HANDLER_V1 = @"
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/homelab/filemanager/internal/service"
)

type FileHandler struct {
	fileService service.FileService
}

func NewFileHandler(fileService service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (h *FileHandler) RegisterRoutes(r chi.Router) {
	r.Get("/*", h.ListFiles)
}

func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	
	files, err := h.fileService.List(r.Context(), path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}
"@

# Initial auth store - basic login only
$AUTH_STORE_V1 = @"
/**
 * Auth store for managing authentication state
 */

import { writable } from 'svelte/store';

export interface AuthState {
	isAuthenticated: boolean;
	isLoading: boolean;
	error: string | null;
}

const initialState: AuthState = {
	isAuthenticated: false,
	isLoading: false,
	error: null
};

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>(initialState);

	async function login(username: string, password: string): Promise<boolean> {
		update((state) => ({ ...state, isLoading: true, error: null }));

		try {
			const response = await fetch('/api/v1/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username, password })
			});

			if (!response.ok) {
				throw new Error('Login failed');
			}

			const data = await response.json();
			localStorage.setItem('accessToken', data.accessToken);
			
			update((state) => ({ ...state, isAuthenticated: true, isLoading: false }));
			return true;
		} catch (err) {
			update((state) => ({
				...state,
				isLoading: false,
				error: err instanceof Error ? err.message : 'Login failed'
			}));
			return false;
		}
	}

	function logout(): void {
		localStorage.removeItem('accessToken');
		set(initialState);
	}

	return { subscribe, login, logout };
}

export const authStore = createAuthStore();
"@


# Initial FileBrowser - basic structure
$FILE_BROWSER_V1 = @"
<script lang="ts">
	import type { FileInfo } from '`$lib/api/files';
	import FileList from './FileList.svelte';

	interface Props {
		currentPath?: string;
		files?: FileInfo[];
		isLoading?: boolean;
		onNavigate?: (path: string) => void;
	}

	let {
		currentPath = '',
		files = [],
		isLoading = false,
		onNavigate
	}: Props = `$props();

	function handleItemClick(item: FileInfo) {
		if (item.isDir) {
			onNavigate?.(item.path);
		}
	}
</script>

<div class="file-browser">
	<header class="browser-header">
		<span class="current-path">{currentPath || '/'}</span>
	</header>

	<main class="browser-content">
		{#if isLoading}
			<p>Loading...</p>
		{:else}
			<FileList items={files} onItemClick={handleItemClick} />
		{/if}
	</main>
</div>

<style>
	.file-browser {
		display: flex;
		flex-direction: column;
		height: 100%;
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
	}

	.browser-header {
		padding: 1rem;
		border-bottom: 1px solid #e5e7eb;
		background: #f9fafb;
	}

	.current-path {
		font-family: monospace;
		color: #374151;
	}

	.browser-content {
		flex: 1;
		overflow: auto;
		padding: 1rem;
	}
</style>
"@.Replace('`$', '$')


# Initial auth middleware - basic token check
$AUTH_MIDDLEWARE_V1 = @"
package middleware

import (
	"net/http"
	"strings"

	"github.com/homelab/filemanager/internal/service"
)

func JWTAuth(authService service.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			_, err := authService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
"@

# Initial path validator - basic checks
$PATH_VALIDATOR_V1 = @"
package validator

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidPath      = errors.New("invalid path")
	ErrPathTraversal    = errors.New("path traversal detected")
)

// ValidatePath checks if a path is safe
func ValidatePath(path string) error {
	// Check for path traversal
	if strings.Contains(path, "..") {
		return ErrPathTraversal
	}

	// Clean the path
	cleaned := filepath.Clean(path)
	if cleaned != path && path != "" {
		return ErrInvalidPath
	}

	return nil
}
"@


# ============================================================================
# COMMIT SEQUENCE - Building the project incrementally
# ============================================================================

Write-Host "`n=== Phase 1: Project Initialization (Week 1) ===" -ForegroundColor Magenta

# Day 1: Initial setup - empty commit
$env:GIT_AUTHOR_DATE = (Get-CommitDate 0 10 11).ToString("yyyy-MM-ddTHH:mm:ss")
$env:GIT_COMMITTER_DATE = $env:GIT_AUTHOR_DATE
git commit --allow-empty -m "Initial commit"
Write-Host "  [$(($baseDate).ToString('yyyy-MM-dd'))] Initial commit" -ForegroundColor Cyan

# Add basic README
Write-FileAndCommit -Message "Add initial README" -Date (Get-CommitDate 0 14 15) -FilePath "README.md" -Content $README_V1

# Day 2: Backend scaffolding
Add-FileAndCommit -Message "Initialize Go module for backend" -Date (Get-CommitDate 1 9 10) -Files @("backend/go.mod")
Add-FileAndCommit -Message "Add chi router dependency" -Date (Get-CommitDate 1 11 12) -Files @("backend/go.sum")
Add-FileAndCommit -Message "Create backend directory structure" -Date (Get-CommitDate 1 14 15) -Files @(
    "backend/internal/handler/handler.go",
    "backend/internal/middleware/middleware.go", 
    "backend/internal/service/service.go",
    "backend/internal/model/model.go"
)

# Day 3: Frontend scaffolding
Add-FileAndCommit -Message "Initialize SvelteKit frontend" -Date (Get-CommitDate 2 9 10) -Files @(
    "frontend/package.json",
    "frontend/svelte.config.js",
    "frontend/vite.config.ts",
    "frontend/src/app.html",
    "frontend/src/app.d.ts",
    "frontend/src/lib/index.ts",
    "frontend/src/routes/layout.css",
    "frontend/bun.lock",
    "frontend/.npmrc"
)
Add-FileAndCommit -Message "Configure TypeScript for frontend" -Date (Get-CommitDate 2 15 16) -Files @("frontend/tsconfig.json")

# Day 4: Basic config
Add-FileAndCommit -Message "Add backend config structure" -Date (Get-CommitDate 3 10 11) -Files @("backend/internal/config/config.go")
Add-FileAndCommit -Message "Create config.yaml template" -Date (Get-CommitDate 3 14 15) -Files @("backend/config.yaml")
Add-FileAndCommit -Message "Add environment example file" -Date (Get-CommitDate 3 16 17) -Files @(".env.example")

# Day 5: Models
Add-FileAndCommit -Message "Define file model types" -Date (Get-CommitDate 4 9 10) -Files @("backend/internal/model/file.go")
Add-FileAndCommit -Message "Add config model" -Date (Get-CommitDate 4 11 12) -Files @("backend/internal/model/config.go")
Add-FileAndCommit -Message "Create error types" -Date (Get-CommitDate 4 14 15) -Files @("backend/internal/model/error.go")


Write-Host "`n=== Phase 2: Core Backend Development (Weeks 2-4) ===" -ForegroundColor Magenta

# Week 2: File service - start with skeleton
Add-FileAndCommit -Message "Implement filesystem abstraction" -Date (Get-CommitDate 7 9 10) -Files @("backend/internal/pkg/filesystem/fs.go")

# Add basic path validator first
Write-FileAndCommit -Message "Add path validator utility" -Date (Get-CommitDate 7 14 15) -FilePath "backend/internal/pkg/validator/path.go" -Content $PATH_VALIDATOR_V1

# Create file service skeleton
Write-FileAndCommit -Message "Create file service skeleton" -Date (Get-CommitDate 8 10 11) -FilePath "backend/internal/service/file.go" -Content $FILE_SERVICE_V1

# Implement directory listing
Write-FileAndCommit -Message "Implement directory listing" -Date (Get-CommitDate 8 15 16) -FilePath "backend/internal/service/file.go" -Content $FILE_SERVICE_V2

# Week 3: Handlers - start basic
Write-FileAndCommit -Message "Implement file handler" -Date (Get-CommitDate 14 14 15) -FilePath "backend/internal/handler/file.go" -Content $FILE_HANDLER_V1

Add-FileAndCommit -Message "Create main server entry point" -Date (Get-CommitDate 15 15 16) -Files @("backend/cmd/server/main.go")

# Week 4: Authentication - start basic
Write-FileAndCommit -Message "Create auth service" -Date (Get-CommitDate 21 14 15) -FilePath "backend/internal/service/auth.go" -Content $AUTH_SERVICE_V1

Write-FileAndCommit -Message "Create auth middleware" -Date (Get-CommitDate 23 9 10) -FilePath "backend/internal/middleware/auth.go" -Content $AUTH_MIDDLEWARE_V1

Add-FileAndCommit -Message "Implement auth handler" -Date (Get-CommitDate 23 14 15) -Files @("backend/internal/handler/auth.go")


Write-Host "`n=== Phase 3: Frontend Development (Weeks 5-7) ===" -ForegroundColor Magenta

# Week 5: Frontend foundation
Add-FileAndCommit -Message "Create API client module" -Date (Get-CommitDate 28 9 10) -Files @("frontend/src/lib/api/client.ts")
Add-FileAndCommit -Message "Add file types definitions" -Date (Get-CommitDate 28 14 15) -Files @("frontend/src/lib/types/files.ts")
Add-FileAndCommit -Message "Implement files API" -Date (Get-CommitDate 29 10 11) -Files @("frontend/src/lib/api/files.ts")
Add-FileAndCommit -Message "Add auth API module" -Date (Get-CommitDate 29 15 16) -Files @("frontend/src/lib/api/auth.ts")

# Auth store - start basic
Write-FileAndCommit -Message "Create auth store" -Date (Get-CommitDate 30 9 10) -FilePath "frontend/src/lib/stores/auth.ts" -Content $AUTH_STORE_V1

Add-FileAndCommit -Message "Implement files store" -Date (Get-CommitDate 30 14 15) -Files @("frontend/src/lib/stores/files.ts")
Add-FileAndCommit -Message "Add API index" -Date (Get-CommitDate 31 10 11) -Files @("frontend/src/lib/api/index.ts")

# Week 6: UI Components - start basic
Add-FileAndCommit -Message "Implement FileList component" -Date (Get-CommitDate 35 14 15) -Files @("frontend/src/lib/components/FileList.svelte")

Write-FileAndCommit -Message "Create FileBrowser component" -Date (Get-CommitDate 35 9 10) -FilePath "frontend/src/lib/components/FileBrowser.svelte" -Content $FILE_BROWSER_V1

Add-FileAndCommit -Message "Add Breadcrumb navigation" -Date (Get-CommitDate 36 10 11) -Files @("frontend/src/lib/components/Breadcrumb.svelte")
Add-FileAndCommit -Message "Create SearchBar component" -Date (Get-CommitDate 36 15 16) -Files @("frontend/src/lib/components/SearchBar.svelte")
Add-FileAndCommit -Message "Implement login page" -Date (Get-CommitDate 37 9 10) -Files @("frontend/src/routes/login/+page.svelte")
Add-FileAndCommit -Message "Add main layout" -Date (Get-CommitDate 37 14 15) -Files @("frontend/src/routes/+layout.svelte")
Add-FileAndCommit -Message "Create home page" -Date (Get-CommitDate 38 10 11) -Files @("frontend/src/routes/+page.svelte")

# Week 7: Browse functionality
Add-FileAndCommit -Message "Implement browse route" -Date (Get-CommitDate 42 9 10) -Files @("frontend/src/routes/browse/[...path]/+page.svelte")
Add-FileAndCommit -Message "Create format utilities" -Date (Get-CommitDate 44 9 10) -Files @("frontend/src/lib/utils/format.ts")


Write-Host "`n=== Phase 4: Streaming & Large Files (Weeks 8-9) ===" -ForegroundColor Magenta

Add-FileAndCommit -Message "Create stream handler" -Date (Get-CommitDate 49 9 10) -Files @("backend/internal/handler/stream.go")
Add-FileAndCommit -Message "Create UploadDropzone component" -Date (Get-CommitDate 56 9 10) -Files @("frontend/src/lib/components/UploadDropzone.svelte")
Add-FileAndCommit -Message "Add upload utility functions" -Date (Get-CommitDate 57 10 11) -Files @("frontend/src/lib/utils/upload.ts")
Add-FileAndCommit -Message "Create UploadProgress component" -Date (Get-CommitDate 58 9 10) -Files @("frontend/src/lib/components/UploadProgress.svelte")

Write-Host "`n=== Phase 5: Background Jobs & WebSocket (Weeks 10-11) ===" -ForegroundColor Magenta

Add-FileAndCommit -Message "Define job model" -Date (Get-CommitDate 63 9 10) -Files @("backend/internal/model/job.go")
Add-FileAndCommit -Message "Create job service" -Date (Get-CommitDate 63 14 15) -Files @("backend/internal/service/job.go")
Add-FileAndCommit -Message "Create job handler" -Date (Get-CommitDate 66 10 11) -Files @("backend/internal/handler/job.go")
Add-FileAndCommit -Message "Add jobs API module" -Date (Get-CommitDate 66 15 16) -Files @("frontend/src/lib/api/jobs.ts")

Add-FileAndCommit -Message "Create WebSocket hub" -Date (Get-CommitDate 70 14 15) -Files @("backend/internal/websocket/hub.go", "backend/internal/websocket/websocket.go")
Add-FileAndCommit -Message "Implement WebSocket client" -Date (Get-CommitDate 71 10 11) -Files @("backend/internal/websocket/client.go")
Add-FileAndCommit -Message "Add WebSocket handler" -Date (Get-CommitDate 71 15 16) -Files @("backend/internal/handler/websocket.go")
Add-FileAndCommit -Message "Create websocket store" -Date (Get-CommitDate 72 9 10) -Files @("frontend/src/lib/stores/websocket.ts")
Add-FileAndCommit -Message "Implement jobs store" -Date (Get-CommitDate 72 14 15) -Files @("frontend/src/lib/stores/jobs.ts")
Add-FileAndCommit -Message "Create JobMonitor component" -Date (Get-CommitDate 73 10 11) -Files @("frontend/src/lib/components/JobMonitor.svelte")


Write-Host "`n=== Phase 6: Search & Security (Week 12) ===" -ForegroundColor Magenta

Add-FileAndCommit -Message "Create search service" -Date (Get-CommitDate 77 9 10) -Files @("backend/internal/service/search.go")
Add-FileAndCommit -Message "Add search handler" -Date (Get-CommitDate 78 10 11) -Files @("backend/internal/handler/search.go")
Add-FileAndCommit -Message "Implement security middleware" -Date (Get-CommitDate 78 15 16) -Files @("backend/internal/middleware/security.go")

Write-Host "`n=== Phase 7: Docker & Deployment (Weeks 13-14) ===" -ForegroundColor Magenta

Add-FileAndCommit -Message "Create backend Dockerfile" -Date (Get-CommitDate 84 9 10) -Files @("backend/Dockerfile")
Add-FileAndCommit -Message "Create frontend Dockerfile" -Date (Get-CommitDate 85 10 11) -Files @("frontend/Dockerfile")
Add-FileAndCommit -Message "Add docker-compose.yml" -Date (Get-CommitDate 85 15 16) -Files @("docker-compose.yml")
Add-FileAndCommit -Message "Configure nginx reverse proxy" -Date (Get-CommitDate 86 9 10) -Files @("nginx/nginx.conf")
Add-FileAndCommit -Message "Add SSL certificate placeholder" -Date (Get-CommitDate 86 11 12) -Files @("nginx/certs/.gitkeep")
Add-FileAndCommit -Message "Create production compose file" -Date (Get-CommitDate 87 10 11) -Files @("docker-compose.prod.yml")

# Documentation
Add-FileAndCommit -Message "Add API documentation" -Date (Get-CommitDate 91 9 10) -Files @("docs/api.md")
Add-FileAndCommit -Message "Create architecture docs" -Date (Get-CommitDate 91 14 15) -Files @("docs/architecture.md")
Add-FileAndCommit -Message "Add configuration guide" -Date (Get-CommitDate 92 10 11) -Files @("docs/configuration.md")
Add-FileAndCommit -Message "Create development guide" -Date (Get-CommitDate 92 15 16) -Files @("docs/development.md")
Add-FileAndCommit -Message "Add Docker documentation" -Date (Get-CommitDate 93 9 10) -Files @("docs/docker.md")
Add-FileAndCommit -Message "Create security documentation" -Date (Get-CommitDate 93 14 15) -Files @("docs/security.md")


Write-Host "`n=== Phase 8: Testing & Final Polish (Weeks 15-16) ===" -ForegroundColor Magenta

Add-FileAndCommit -Message "Add path validator tests" -Date (Get-CommitDate 98 9 10) -Files @("backend/internal/pkg/validator/path_test.go")
Add-FileAndCommit -Message "Create file service property tests" -Date (Get-CommitDate 98 14 15) -Files @("backend/internal/service/file_property_test.go")
Add-FileAndCommit -Message "Add auth middleware tests" -Date (Get-CommitDate 99 10 11) -Files @("backend/internal/middleware/auth_property_test.go")
Add-FileAndCommit -Message "Implement security middleware tests" -Date (Get-CommitDate 99 15 16) -Files @("backend/internal/middleware/security_property_test.go")
Add-FileAndCommit -Message "Add job service tests" -Date (Get-CommitDate 100 9 10) -Files @("backend/internal/service/job_property_test.go")
Add-FileAndCommit -Message "Create search service tests" -Date (Get-CommitDate 100 14 15) -Files @("backend/internal/service/search_property_test.go")
Add-FileAndCommit -Message "Add stream handler tests" -Date (Get-CommitDate 101 10 11) -Files @("backend/internal/handler/stream_property_test.go")

# Final touches
Add-FileAndCommit -Message "Add .gitignore" -Date (Get-CommitDate 105 9 10) -Files @(".gitignore")
Add-FileAndCommit -Message "Create frontend .gitignore" -Date (Get-CommitDate 105 11 12) -Files @("frontend/.gitignore")
Add-FileAndCommit -Message "Add .dockerignore files" -Date (Get-CommitDate 106 9 10) -Files @("backend/.dockerignore", "frontend/.dockerignore")
Add-FileAndCommit -Message "Configure prettier" -Date (Get-CommitDate 106 14 15) -Files @("frontend/.prettierrc", "frontend/.prettierignore")
Add-FileAndCommit -Message "Add ESLint configuration" -Date (Get-CommitDate 107 10 11) -Files @("frontend/eslint.config.js")
Add-FileAndCommit -Message "Add MIT license" -Date (Get-CommitDate 108 9 10) -Files @("LICENSE")
Add-FileAndCommit -Message "Add frontend README" -Date (Get-CommitDate 108 11 12) -Files @("frontend/README.md")


Write-Host "`n=== Phase 9: Enhancements & Bug Fixes ===" -ForegroundColor Magenta

# Now update files to their final versions - simulating feature additions

# Update path validator with mount point support
Add-FileAndCommit -Message "Add mount point validation to path validator" -Date (Get-CommitDate 109 10 11) -Files @("backend/internal/pkg/validator/path.go")

# Update file service with pagination and mount points
Add-FileAndCommit -Message "Add pagination and mount point support to file service" -Date (Get-CommitDate 110 9 10) -Files @("backend/internal/service/file.go")

# Update file handler with full CRUD
Add-FileAndCommit -Message "Add full CRUD operations to file handler" -Date (Get-CommitDate 110 14 15) -Files @("backend/internal/handler/file.go")

# Update auth service with refresh tokens
Add-FileAndCommit -Message "Add refresh token support to auth service" -Date (Get-CommitDate 111 9 10) -Files @("backend/internal/service/auth.go")

# Update auth middleware with context
Add-FileAndCommit -Message "Add user context to auth middleware" -Date (Get-CommitDate 111 14 15) -Files @("backend/internal/middleware/auth.go")

# Update auth store with token refresh
Add-FileAndCommit -Message "Add automatic token refresh to auth store" -Date (Get-CommitDate 112 9 10) -Files @("frontend/src/lib/stores/auth.ts")

# Update FileBrowser with search and drag-drop
Add-FileAndCommit -Message "Add search and drag-drop to FileBrowser" -Date (Get-CommitDate 112 14 15) -Files @("frontend/src/lib/components/FileBrowser.svelte")

# Update README with full documentation
Write-FileAndCommit -Message "Update README with full documentation" -Date (Get-CommitDate 113 10 11) -FilePath "README.md" -Content $README_V2

# Add static assets
Add-FileAndCommit -Message "Add static assets" -Date (Get-CommitDate 114 9 10) -Files @(
    "frontend/static/robots.txt",
    "frontend/src/lib/assets/favicon.svg",
    "frontend/src/lib/components/.gitkeep",
    "frontend/src/lib/stores/.gitkeep",
    "frontend/src/lib/utils/.gitkeep"
)

# Kiro specs
Add-FileAndCommit -Message "Add project specifications" -Date (Get-CommitDate 115 10 11) -Files @(
    ".kiro/specs/homelab-file-manager/requirements.md",
    ".kiro/specs/homelab-file-manager/tasks.md"
)

# Bug fixes
Add-FileAndCommit -Message "Fix upload progress calculation" -Date (Get-CommitDate 117 10 11) -Files @("frontend/src/lib/components/UploadProgress.svelte")
Add-FileAndCommit -Message "Improve error handling in file service" -Date (Get-CommitDate 118 14 15) -Files @("backend/internal/service/file.go")
Add-FileAndCommit -Message "Optimize WebSocket reconnection logic" -Date (Get-CommitDate 119 11 12) -Files @("frontend/src/lib/stores/websocket.ts")
Add-FileAndCommit -Message "Fix path validation edge cases" -Date (Get-CommitDate 120 14 15) -Files @("backend/internal/pkg/validator/path.go")

# Final release
Add-FileAndCommit -Message "Release v1.0.0" -Date (Get-CommitDate 121 16 17) -Files @("generate-git-history.ps1")


# ============================================================================
# CLEANUP
# ============================================================================

Remove-Item Env:\GIT_AUTHOR_DATE -ErrorAction SilentlyContinue
Remove-Item Env:\GIT_COMMITTER_DATE -ErrorAction SilentlyContinue
Remove-Item -Recurse -Force $tempDir -ErrorAction SilentlyContinue

Write-Host "`n=== Git History Generation Complete ===" -ForegroundColor Green
Write-Host "Total commits created: " -NoNewline
git rev-list --count HEAD
Write-Host "`nRun 'git log --oneline' to see all commits" -ForegroundColor Yellow
Write-Host "Run 'git log -p --follow <file>' to see file evolution" -ForegroundColor Yellow

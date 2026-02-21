// Package file provides the infrastructure layer for interacting with files and folders
// in Microsoft OneDrive.
//
// It implements the domain-level file services and repositories by communicating
// with the Microsoft Graph API. The package handles path normalization,
// error mapping from Graph-specific errors to domain errors, and provides
// caching mechanisms for both file metadata and content to improve performance
// and reduce API calls.
package file

package service

import "time"

// getDemoObjects returns static demo objects for demo users
// These are displayed when BB_ENABLE_DEMO_LOGIN is enabled and the user is a demo user
func getDemoObjects(bucketName, prefix string) []BucketObject {
	// Normalize prefix
	if prefix != "" && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	allObjects := map[string][]BucketObject{
		"product-images": {
			// Root level files
			{
				Key:          "product-hero-banner.jpg",
				Name:         "product-hero-banner.jpg",
				Kind:         "file",
				Size:         "2.4 MB",
				LastModified: time.Now().Add(-72 * time.Hour),
				Icon:         "image",
				IconColor:    "text-blue-500",
			},
			// Thumbnails folder
			{
				Key:          "thumbnails/",
				Name:         "thumbnails",
				Kind:         "folder",
				Size:         "-",
				LastModified: time.Now().Add(-48 * time.Hour),
				Icon:         "folder",
				IconColor:    "text-amber-500",
			},
			// Files inside thumbnails/
			{
				Key:          "thumbnails/product-thumbnail-001.jpg",
				Name:         "product-thumbnail-001.jpg",
				Kind:         "file",
				Size:         "156 KB",
				LastModified: time.Now().Add(-48 * time.Hour),
				Icon:         "image",
				IconColor:    "text-blue-500",
			},
			{
				Key:          "thumbnails/product-thumbnail-002.jpg",
				Name:         "product-thumbnail-002.jpg",
				Kind:         "file",
				Size:         "142 KB",
				LastModified: time.Now().Add(-36 * time.Hour),
				Icon:         "image",
				IconColor:    "text-blue-500",
			},
			{
				Key:          "thumbnails/product-thumbnail-003.jpg",
				Name:         "product-thumbnail-003.jpg",
				Kind:         "file",
				Size:         "168 KB",
				LastModified: time.Now().Add(-24 * time.Hour),
				Icon:         "image",
				IconColor:    "text-blue-500",
			},
		},
		"user-uploads": {
			// Avatars folder
			{
				Key:          "avatars/",
				Name:         "avatars",
				Kind:         "folder",
				Size:         "-",
				LastModified: time.Now().Add(-120 * time.Hour),
				Icon:         "folder",
				IconColor:    "text-amber-500",
			},
			// Documents folder
			{
				Key:          "documents/",
				Name:         "documents",
				Kind:         "folder",
				Size:         "-",
				LastModified: time.Now().Add(-96 * time.Hour),
				Icon:         "folder",
				IconColor:    "text-amber-500",
			},
			// Files inside avatars/
			{
				Key:          "avatars/user-avatar-john.png",
				Name:         "user-avatar-john.png",
				Kind:         "file",
				Size:         "45 KB",
				LastModified: time.Now().Add(-120 * time.Hour),
				Icon:         "image",
				IconColor:    "text-blue-500",
			},
			{
				Key:          "avatars/user-avatar-sarah.png",
				Name:         "user-avatar-sarah.png",
				Kind:         "file",
				Size:         "52 KB",
				LastModified: time.Now().Add(-100 * time.Hour),
				Icon:         "image",
				IconColor:    "text-blue-500",
			},
			// Files inside documents/
			{
				Key:          "documents/report-q4-2024.pdf",
				Name:         "report-q4-2024.pdf",
				Kind:         "file",
				Size:         "1.2 MB",
				LastModified: time.Now().Add(-96 * time.Hour),
				Icon:         "description",
				IconColor:    "text-red-500",
			},
			{
				Key:          "documents/presentation-slides.pptx",
				Name:         "presentation-slides.pptx",
				Kind:         "file",
				Size:         "3.5 MB",
				LastModified: time.Now().Add(-84 * time.Hour),
				Icon:         "slideshow",
				IconColor:    "text-orange-500",
			},
			{
				Key:          "documents/meeting-notes.docx",
				Name:         "meeting-notes.docx",
				Kind:         "file",
				Size:         "156 KB",
				LastModified: time.Now().Add(-72 * time.Hour),
				Icon:         "description",
				IconColor:    "text-blue-600",
			},
		},
		"dev-testing": {
			// Root level files
			{
				Key:          "test-data.json",
				Name:         "test-data.json",
				Kind:         "file",
				Size:         "2.1 KB",
				LastModified: time.Now().Add(-12 * time.Hour),
				Icon:         "code",
				IconColor:    "text-green-500",
			},
			{
				Key:          "README.md",
				Name:         "README.md",
				Kind:         "file",
				Size:         "1.8 KB",
				LastModified: time.Now().Add(-24 * time.Hour),
				Icon:         "description",
				IconColor:    "text-slate-500",
			},
			{
				Key:          "config.yaml",
				Name:         "config.yaml",
				Kind:         "file",
				Size:         "892 bytes",
				LastModified: time.Now().Add(-6 * time.Hour),
				Icon:         "settings",
				IconColor:    "text-purple-500",
			},
			// Scripts folder
			{
				Key:          "scripts/",
				Name:         "scripts",
				Kind:         "folder",
				Size:         "-",
				LastModified: time.Now().Add(-48 * time.Hour),
				Icon:         "folder",
				IconColor:    "text-amber-500",
			},
			// Files inside scripts/
			{
				Key:          "scripts/deploy.sh",
				Name:         "deploy.sh",
				Kind:         "file",
				Size:         "4.2 KB",
				LastModified: time.Now().Add(-48 * time.Hour),
				Icon:         "terminal",
				IconColor:    "text-slate-600",
			},
			{
				Key:          "scripts/backup.sh",
				Name:         "backup.sh",
				Kind:         "file",
				Size:         "2.8 KB",
				LastModified: time.Now().Add(-36 * time.Hour),
				Icon:         "terminal",
				IconColor:    "text-slate-600",
			},
		},
	}

	objects, exists := allObjects[bucketName]
	if !exists {
		return []BucketObject{}
	}

	// Filter by prefix and construct proper result
	var result []BucketObject
	folderMap := make(map[string]bool)

	for _, obj := range objects {
		// If no prefix, only show root level items
		if prefix == "" {
			// Root files (no slash in key)
			if obj.Kind == "file" && !containsChar(obj.Key, '/') {
				result = append(result, obj)
			}
			// Root folders (ends with / and only one segment)
			if obj.Kind == "folder" {
				result = append(result, obj)
			}
			continue
		}

		// With prefix, filter objects that start with prefix
		if !startsWithPrefix(obj.Key, prefix) {
			continue
		}

		// Remove prefix from key for display
		relativeKey := obj.Key[len(prefix):]

		// If it's a direct file (no slash after prefix)
		if obj.Kind == "file" && !containsChar(relativeKey, '/') {
			displayObj := obj
			displayObj.Key = obj.Key
			displayObj.Name = relativeKey
			result = append(result, displayObj)
			continue
		}

		// If it's a subfolder
		if containsChar(relativeKey, '/') {
			// Extract immediate subfolder name
			parts := splitFirst(relativeKey, '/')
			folderName := parts[0]
			folderKey := prefix + folderName + "/"

			// Only add folder once
			if !folderMap[folderKey] {
				folderMap[folderKey] = true
				result = append(result, BucketObject{
					Key:          folderKey,
					Name:         folderName,
					Kind:         "folder",
					Size:         "-",
					LastModified: time.Now().Add(-24 * time.Hour),
					Icon:         "folder",
					IconColor:    "text-amber-500",
				})
			}
		}
	}

	return result
}

// Helper functions
func containsChar(s string, c rune) bool {
	for _, ch := range s {
		if ch == c {
			return true
		}
	}
	return false
}

func startsWithPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}

func splitFirst(s string, delimiter rune) []string {
	for i, ch := range s {
		if ch == delimiter {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}

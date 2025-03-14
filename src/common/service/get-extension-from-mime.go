package common_service

// Helper function to get the file extension from the MIME type
func GetExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "application/pdf":
		return "pdf"
	// Add more MIME types and their corresponding extensions as needed
	default:
		return "bin"
	}
}

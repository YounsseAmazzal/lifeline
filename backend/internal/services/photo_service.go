package services

import (
	"context"
	"errors"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// PhotoService interface: باش نسهلو التيست من بعد
type PhotoService interface {
	UploadImage(file *multipart.FileHeader) (string, string, error) // Returns: URL, PublicID, Error
	DeleteImage(publicID string) error
}

type photoService struct {
	cld *cloudinary.Cloudinary
}

// NewPhotoService constructor
func NewPhotoService() PhotoService {
	url := os.Getenv("CLOUDINARY_URL")
	if url == "" {
		// Fallback for dev only, or panic
		// panic("CLOUDINARY_URL is missing in .env")
		return &photoService{} // Return empty service to avoid crash if not configured yet
	}

	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		panic("Failed to initialize Cloudinary Service ❌: " + err.Error())
	}
	
	return &photoService{cld: cld}
}

// UploadImage: رفع صورة جديدة
func (s *photoService) UploadImage(file *multipart.FileHeader) (string, string, error) {
	if s.cld == nil {
		return "", "", errors.New("cloudinary not configured")
	}

	ctx := context.Background()

	// فتح الملف للقراءة
	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	// إعدادات التحويل (Resize & Crop)
	uploadParams := uploader.UploadParams{
		Folder:         "lifeline-app", // السمية ديال الدوسي فـ Cloudinary
		Transformation: "w_500,h_500,c_fill,g_face", // مربع 500x500 مركز على الوجه
	}

	// Upload API Call
	resp, err := s.cld.Upload.Upload(ctx, src, uploadParams)
	if err != nil {
		return "", "", err
	}

	return resp.SecureURL, resp.PublicID, nil
}

// DeleteImage: مسح صورة قديمة
func (s *photoService) DeleteImage(publicID string) error {
	if s.cld == nil {
		return errors.New("cloudinary not configured")
	}
	if publicID == "" {
		return nil
	}

	ctx := context.Background()
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	
	return err
}
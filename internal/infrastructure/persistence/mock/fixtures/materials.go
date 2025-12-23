package fixtures

import (
	"time"

	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

var (
	// MaterialGuidaSumasID UUIDs fijos para materiales de prueba - coherentes con api-administracion
	MaterialGuidaSumasID = uuid.MustParse("f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	MaterialGuiaRestasID = uuid.MustParse("f2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22")
	MaterialLasPlantasID = uuid.MustParse("f3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33")
	MaterialCicloAguaID  = uuid.MustParse("f4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44")

	// IDs de referencia
	schoolPrimariaID = uuid.MustParse("b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	teacherMathID    = uuid.MustParse("a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22")
	teacherScienceID = uuid.MustParse("a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33")
	seccionPrimer1A  = uuid.MustParse("c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44")
	seccionPrimer1B  = uuid.MustParse("c5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55")

	materialsBaseTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
)

// GetDefaultMaterials retorna mapa de materiales educativos de prueba
func GetDefaultMaterials() map[uuid.UUID]*pgentities.Material {
	materials := make(map[uuid.UUID]*pgentities.Material)

	sp := func(s string) *string { return &s }
	up := func(u uuid.UUID) *uuid.UUID { return &u }

	materials[MaterialGuidaSumasID] = &pgentities.Material{
		ID:                  MaterialGuidaSumasID,
		SchoolID:            schoolPrimariaID,
		UploadedByTeacherID: teacherMathID,
		AcademicUnitID:      up(seccionPrimer1A),
		Title:               "Guía de Sumas",
		Description:         sp("Material educativo sobre sumas básicas"),
		Subject:             sp("Matemáticas"),
		Grade:               sp("Primaria"),
		FileURL:             "https://s3.example.com/materials/math/suma.pdf",
		FileType:            "application/pdf",
		FileSizeBytes:       1048576,
		Status:              "ready",
		IsPublic:            true,
		CreatedAt:           materialsBaseTime,
		UpdatedAt:           materialsBaseTime,
	}

	materials[MaterialGuiaRestasID] = &pgentities.Material{
		ID:                  MaterialGuiaRestasID,
		SchoolID:            schoolPrimariaID,
		UploadedByTeacherID: teacherMathID,
		AcademicUnitID:      up(seccionPrimer1A),
		Title:               "Guía de Restas",
		Description:         sp("Material educativo sobre restas básicas"),
		Subject:             sp("Matemáticas"),
		Grade:               sp("Primaria"),
		FileURL:             "https://s3.example.com/materials/math/resta.pdf",
		FileType:            "application/pdf",
		FileSizeBytes:       950000,
		Status:              "ready",
		IsPublic:            true,
		CreatedAt:           materialsBaseTime,
		UpdatedAt:           materialsBaseTime,
	}

	materials[MaterialLasPlantasID] = &pgentities.Material{
		ID:                  MaterialLasPlantasID,
		SchoolID:            schoolPrimariaID,
		UploadedByTeacherID: teacherScienceID,
		AcademicUnitID:      up(seccionPrimer1B),
		Title:               "Las Plantas",
		Description:         sp("Video educativo sobre plantas"),
		Subject:             sp("Ciencias Naturales"),
		Grade:               sp("Primaria"),
		FileURL:             "https://s3.example.com/materials/science/plantas.mp4",
		FileType:            "video/mp4",
		FileSizeBytes:       52428800,
		Status:              "ready",
		IsPublic:            true,
		CreatedAt:           materialsBaseTime,
		UpdatedAt:           materialsBaseTime,
	}

	materials[MaterialCicloAguaID] = &pgentities.Material{
		ID:                  MaterialCicloAguaID,
		SchoolID:            schoolPrimariaID,
		UploadedByTeacherID: teacherScienceID,
		AcademicUnitID:      up(seccionPrimer1B),
		Title:               "El Ciclo del Agua",
		Description:         sp("Presentación sobre el ciclo del agua"),
		Subject:             sp("Ciencias Naturales"),
		Grade:               sp("Primaria"),
		FileURL:             "https://s3.example.com/materials/science/agua.pptx",
		FileType:            "application/vnd.ms-powerpoint",
		FileSizeBytes:       2097152,
		Status:              "ready",
		IsPublic:            true,
		CreatedAt:           materialsBaseTime,
		UpdatedAt:           materialsBaseTime,
	}

	return materials
}

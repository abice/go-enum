//go:generate ../bin/go-enum -b example

package example

// @marshal:true @sql:false @prefix:"My"
// ENUM(pending, running, completed, failed)
type AnnotationStatus string

// @noprefix @nocase
// ENUM(annotation_red, annotation_green, annotation_blue)
type AnnotationColor string

// @marshal @sql @marshal
// ENUM(one, two, three)
type AnnotationNumber int

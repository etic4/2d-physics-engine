package physics

import (
	"fmt"
	"math"
)

//Vec2 Représente un vecteur 2D
type Vec2 struct {
	X float64
	Y float64
}

//NewVector crée un nouveau vecteur
func NewVector(x, y float64) Vec2 {
	return Vec2{x, y}
}

//Add additionne deux vecteurs
func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v.X + v2.X, v.Y + v2.Y}
}

//AddScalar ajoute 's' à chaque composante
func (v Vec2) AddScalar(s float64) Vec2 {
	return Vec2{v.X + s, v.Y + s}
}

//Sub soustrait
func (v Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v.X - v2.X, v.Y - v2.Y}
}

//SubScalar soustrait s à chaque composante
func (v Vec2) SubScalar(s float64) Vec2 {
	return Vec2{v.X - s, v.Y - s}
}

//DotProduct produit scalaire de deux vecteurs
func (v Vec2) DotProduct(v2 Vec2) float64 {
	return v.X*v2.X + v.Y*v2.Y
}

//Mult multiplie le vecteur par un scalaire
func (v Vec2) Mult(s float64) Vec2 {
	return Vec2{v.X * s, v.Y * s}
}

//Div divise le vecteur par un scalaire
func (v Vec2) Div(s float64) Vec2 {
	return Vec2{v.X / s, v.Y / s}
}

//Neg négation d'un vecteur (180°)
func (v Vec2) Neg() Vec2 {
	return Vec2{-1 * v.X, -1 * v.Y}
}

//Length norme du vecteur
func (v Vec2) Length() float64 {
	return math.Sqrt(v.DotProduct(v))
}

//Normalize retourne vecteur unitaire (normalisé)
func (v Vec2) Normalize() Vec2 {
	return v.Div(v.Length())
}

//Distance distance entre deux vecteurs
func (v Vec2) Distance(v2 Vec2) float64 {
	// dx := v2.X - v.X
	// dy := v2.Y - v.Y
	// d := v2.Sub(v)

	// return math.Sqrt(dx*dx + dy*dy)
	return v2.Sub(v).Length()
}

//DistanceCarree pour la facilité
func (v Vec2) DistanceCarree(v2 Vec2) float64 {
	dx := v2.X - v.X
	dy := v2.Y - v.Y

	return dx*dx + dy*dy
}

//Round arrondis les deux composantes du vecteur à la nième decimale
func (v Vec2) Round(decimals int) Vec2 {
	prop := math.Pow10(decimals)
	return Vec2{math.Round(v.X*prop) / prop, math.Round(v.Y*prop) / prop}
}

//Retourne imprimable
func (v Vec2) String() string {
	return fmt.Sprintf("%v,%v", v.X, v.Y)
}

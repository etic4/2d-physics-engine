// Massivement inspiré (sinon copié) de
// https://gamedevelopment.tutsplus.com/tutorials/how-to-create-a-custom-2d-physics-engine-the-basics-and-impulse-resolution--gamedev-6331
// et https://github.com/vonWolfehaus/von-physics/blob/master/src/physics/Physics2.js

package physique

import (
	"fmt"
)

const (
	velocityTolerance float64 = 0.001
)

//CollisionInfo Informations sur une collision ou son absence
type CollisionInfo struct {
	first       Shape
	second      Shape
	penetration float64
	normal      Vec2
	resolved    bool
}

//IsColliding retourne true s'il y a collision
func (i *CollisionInfo) IsColliding() bool {
	return i.second != nil
}

//Resolved retourne true si collision déjà résolue
func (i *CollisionInfo) Resolved() bool {
	return i.resolved
}

//SetResolved marque la collision comme résolue
func (i *CollisionInfo) SetResolved(b bool) {
	i.resolved = b
}

//First retourne la première des deux formes impliquées dans la collision
func (i *CollisionInfo) First() Shape {
	return i.first
}

//Second retourne la seconde des deux formes impliquées dans la collision
func (i *CollisionInfo) Second() Shape {
	return i.second
}

//GetShapeForTag retourne les formes impliquées dans la collision
// qui possèdent ce tag
func (i *CollisionInfo) GetShapeForTag(t string) ([]Shape, error) {
	shapes := []Shape{}
	if stringListContains(i.first.Tags(), t) {
		shapes = append(shapes, i.first)
	}

	if stringListContains(i.second.Tags(), t) {
		shapes = append(shapes, i.second)
	}
	err := error(nil)

	if len(shapes) == 0 {
		err = fmt.Errorf("Pas de forme ayant le tag %v", t)
	}
	return shapes, err
}

//GetShapeForName retourne la forme qui a ce nom ou rien et une erreur
func (i *CollisionInfo) GetShapeForName(n string) (Shape, error) {
	if i.first.Name() == n {
		return i.first, nil
	} else if i.second.Name() == n {
		return i.second, nil
	} else {
		return nil, fmt.Errorf("Aucune des formes n'a ce nom: %v", n)
	}
}

//String retourne une version imprimable
func (i *CollisionInfo) String() string {
	return fmt.Sprintf("first: %v  second: %v  penetration: %v  normal: %v",
		i.first.Name(), i.second.Name(), i.penetration, i.normal)
}

//Resolv résoud la collision en déterminant la rectification de position
//	et l'impulsion à donner aux objets
func (i *CollisionInfo) Resolv() {
	// Ne résoud pas plusieurs fois
	if i.resolved {
		return
	}

	//Pas besoin de "séparer" les objets dès lors que la collision est
	// résolue par l'application d'une impulsion
	// i.Separate()

	first := i.first
	second := i.second

	// Détermination de l'impulsion a appliquer aux objets sur la normale de la collision

	// vitesse relative
	f2S := second.Velocity().Sub(first.Velocity())

	// vitesse relative sur la normale de la collision
	vRelAlongNorm := f2S.DotProduct(i.normal)

	// résoud pas si s'éloignent
	if vRelAlongNorm > 0 {
		return
	}

	e := Min(first.Elasticity(), second.Elasticity())

	// calcule impulsion scalaire
	j := -(1 + e) * vRelAlongNorm
	j /= first.InvMass() + second.InvMass()

	impulse := i.normal.Mult(j)

	// appliquer l'impulsion
	first.SetVelocity(first.Velocity().Sub(impulse.Mult(first.InvMass())))
	second.SetVelocity(second.Velocity().Add(impulse.Mult(second.InvMass())))

	//Applique friction
	first.SetVelocity(first.Velocity().Mult(1 - first.Friction()))
	second.SetVelocity(second.Velocity().Mult(1 - second.Friction()))

	// Correction naufrage ("sinking"), "causé par le fait que "la résultante des vitesses
	// est insuffisante pour pousser l'objet hors d'une collision, quand un objet est stationnaire"
	percent := 0.5 // habituellement 20% à 80%
	corr := i.normal.Mult(i.penetration / (first.InvMass() + second.InvMass()) * percent)

	first.SetPos(first.Pos().Sub(corr.Mult(first.InvMass())))
	second.SetPos(second.Pos().Add(corr.Mult(second.InvMass())))

	//met resolved à true, pour ne pas pouvoir résoudre deux fois la même collision
	i.SetResolved(true)
}

//Separate sépare deux objets en revenant à une position pré-collision
func (i *CollisionInfo) Separate() {
	// Ne résoud pas plusieurs fois
	if i.resolved {
		return
	}

	first := i.first
	second := i.second

	// Repositionnement des objets qui s'interpénètrent
	totInvMass := first.InvMass() + second.InvMass()
	first.SetPos(first.Pos().Sub(i.normal.Mult(i.penetration * first.InvMass() / totInvMass)))
	second.SetPos(second.Pos().Add(i.normal.Mult(i.penetration * second.InvMass() / totInvMass)))

	//met resolved à true, pour ne jamais pouvoir résoudre deux fois la même collision
	i.SetResolved(true)

}

//AABBvsAABB Détermine l'ajustement des coordonées de first et second si entrent en collision
// retourne une réponse qui contient la pénétration et la normale de la face
//  sur laquelle il y a collision
func AABBvsAABB(first *Rectangle, second *Rectangle) *CollisionInfo {
	info := &CollisionInfo{}
	info.first = first

	if first.getMax().X < second.Pos().X || first.Pos().X > second.getMax().X {
		return info
	}

	if first.getMax().Y < second.Pos().Y || first.Pos().Y > second.getMax().Y {
		return info
	}

	// il y a collision
	info.second = second

	// différence entre distance de centre à centre et sommes des demi-côtés, sur les deux axes
	// donne la profondeur de pénétration d'une forme dans l'autre
	distance := second.Center().Sub(first.Center())
	px := (first.Width()+second.Width())/2 - Abs(distance.X)
	py := (first.Height()+second.Height())/2 - Abs(distance.Y)

	// choix de l'axe de moindre pénétration
	if px < py {
		sx := Sign(distance.X)
		info.normal.X = sx
		info.penetration = px
	} else {
		sy := Sign(distance.Y)
		info.normal.Y = sy
		info.penetration = py
	}


	// grounded si normale.Y == -1 et corps a une masse
	if info.normal.Y == -1 {
		if first.mass != 0 {
			first.SetGrounded(true)
		}
		if second.mass != 0 {
			second.SetGrounded(true)
		}
	}
	return info
}

//CirclevsCircle génère CollisionInfo pour collisions de cercles
func CirclevsCircle(first *Circle, second *Circle) *CollisionInfo {
	info := &CollisionInfo{}
	info.first = first

	firstCenter, secCenter := first.Center(), second.Center()

	FirstSec := secCenter.Sub(firstCenter)
	dist := FirstSec.Length()

	radSum := first.radius + second.radius

	// La distance entre le deux centres est plus grande que la somme des rayons
	if dist >= radSum {
		return info
	}

	// il y a collision
	info.second = second

	if dist != 0 {
		info.penetration = radSum - dist
		info.normal = FirstSec.Div(dist) // dist déjà calculée

	} else { // les deux cercles on même centre
		info.penetration = first.radius
		info.normal = Vec2{1, 0}
	}

	//TODO: c'est quoi ça?
	// grounded si normale.Y == -1 et corps a une masse
	if info.normal.Y == -1 {
		if first.mass != 0 {
			first.SetGrounded(true)
		}
		if second.mass != 0 {
			second.SetGrounded(true)
		}
	}

	return info

}

//AABBvsCircle génère CollisionInfo pour collisions rectangle/cercle
// Copié de https://gamedevelopment.tutsplus.com/tutorials/how-to-create-a-custom-2d-physics-engine-the-basics-and-impulse-resolution--gamedev-6331
func AABBvsCircle(first *Rectangle, second *Circle) *CollisionInfo {
	info := &CollisionInfo{}
	info.first = first

	n := second.Center().Sub(first.Center())

	closest := n

	xExtent := first.Width() / 2
	yExtent := first.Height() / 2

	closest.X = Clamp(closest.X, -xExtent, xExtent)
	closest.Y = Clamp(closest.Y, -yExtent, yExtent)

	closest = first.Center().Add(closest)

	inside := false

	// Cercle est dans AABB
	if n == closest {
		inside = true

		// Axe le plus proche
		if Abs(n.X) > Abs(n.X) {
			if closest.X > 0 {
				closest.X = xExtent
			} else {
				closest.X = -xExtent
			}
		} else {
			if closest.Y > 0 {
				closest.Y = yExtent
			} else {
				closest.Y = -yExtent
			}
		}
	}

	dist := second.Center().Distance(closest)
	if dist > second.Radius() && !inside {
		return info
	}

	info.second = second
	normale := second.Center().Sub(closest).Normalize()

	if inside {
		info.normal = normale.Neg()
	} else {
		info.normal = normale
	}

	info.penetration = second.Radius() - dist

	// grounded si normale.Y == -1 et corps a une masse
	if info.normal.Y == -1 {
		if first.mass != 0 {
			first.SetGrounded(true)
		}
		if second.mass != 0 {
			second.SetGrounded(true)
		}

	}

	return info
}

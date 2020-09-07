package physique

import "fmt"

//Space Contient toutes les shape
type Space struct {
	shapesList []Shape
	collisions *InfoList
	dt         float64
	gravity    Vec2
}

//Update met l'espace à jour: positions et collisions
func (s *Space) Update() {
	s.updatePositions()
	s.checkCollisions()
}

//Collisions retourne la liste des collisions
func (s *Space) Collisions() *InfoList {
	return s.collisions
}

//Shapes retourne la liste des shapes de l'espace
func (s *Space) Shapes() []Shape {
	return s.shapesList
}

//AddShape Ajoute une shape à l'espace
func (s *Space) AddShape(obj Shape) {
	s.shapesList = append(s.shapesList, obj)
}

//RemoveShape Supprime une shape de l'espace
func (s *Space) RemoveShape(obj Shape) {
	i := 0

	for i < len(s.shapesList) && s.shapesList[i].Name() != obj.Name() {
		i++
	}

	if i < len(s.shapesList) {
		if i < len(s.shapesList)-1 {
			copy(s.shapesList[i:], s.shapesList[i+1:])
		}
		s.shapesList[len(s.shapesList)-1] = nil
		s.shapesList = s.shapesList[:len(s.shapesList)-1]
	}
}

//updatePositions met à jour les positions des formes
func (s *Space) updatePositions() {
	for _, s := range s.shapesList {
		if !s.IsStatic() {
			s.UpdatePos()
		}
	}
}

//SetGravity mets la gravité de l'espace à g
func (s *Space) SetGravity(g Vec2) {
	s.gravity = g
}

//ApplyGravity met la gravité de tous les éléments de l'espace
// à la gravité de l'espace
func (s *Space) ApplyGravity() {
	for _, shape := range s.shapesList {
		shape.SetGravity(s.gravity)
	}
}

//checkCollisions retourne la liste de toutes les collisions
func (s *Space) checkCollisions() {
	collisions := newInfoList()

	for i := 0; i < len(s.shapesList)-1; i++ {
		if s.shapesList[i].IsSolid() { // Ne check pas les formes qui ne collisionnent pas
			for j := i + 1; j < len(s.shapesList); j++ {
				if s.shapesList[j].IsSolid() {
					info := s.dispatchCollisionCheck(s.shapesList[i], s.shapesList[j])
					if info.IsColliding() {
						collisions.Add(info)
					}
				}

			}
		}

	}
	s.collisions = collisions
}

// Dispatch le check de la collision aux fonctions appropriées
func (s *Space) dispatchCollisionCheck(obj1 Shape, obj2 Shape) *CollisionInfo {
	info := &CollisionInfo{}

	switch obj1.ShapeName() {
	case "Circle":
		shape1 := obj1.(*Circle)

		switch obj2.ShapeName() {
		case "Circle":
			shape2 := obj2.(*Circle)
			info = CirclevsCircle(shape1, shape2)
		case "Rectangle":
			shape2 := obj2.(*Rectangle)
			info = AABBvsCircle(shape2, shape1)
		default:
			panic(fmt.Sprintf("Pas de support pour %T\n", obj2.ShapeName()))
		}
	case "Rectangle":
		shape1 := obj1.(*Rectangle)

		switch obj2.ShapeName() {
		case "Circle":
			shape2 := obj2.(*Circle)
			info = AABBvsCircle(shape1, shape2)
		case "Rectangle":
			shape2 := obj2.(*Rectangle)
			info = AABBvsAABB(shape1, shape2)
		default:
			panic(fmt.Sprintf("Pas de support pour %T\n", obj2.ShapeName()))
		}
	}
	return info
}

//InfoList contient les infos de collisions
type InfoList struct {
	infoList []*CollisionInfo
	infoMap  map[string][]*CollisionInfo
}

func newInfoList() *InfoList {
	il := &InfoList{}
	il.infoList = []*CollisionInfo{}
	il.infoMap = map[string][]*CollisionInfo{}
	return il
}

//Add Ajoute une info à la liste
func (l *InfoList) Add(info *CollisionInfo) {
	l.infoList = append(l.infoList, info)

	tagList := info.first.Tags()
	for _, t := range info.second.Tags() {
		if !stringListContains(info.first.Tags(), t) {
			tagList = append(tagList, t)
		}
	}
	l.addInfoFor(info, tagList)
}

//GetAll retourne la liste de collisions non résolues. Si un tag est passé, retourne
// les collisions non résolues relatives à ce tag
func (l *InfoList) GetAll(tags ...string) []*CollisionInfo {
	if len(tags) > 1 {
		panic("Filtrer sur plusieurs tags n'est pas supporté")
	}

	notResolved := []*CollisionInfo{}

	if len(tags) == 0 {
		return l.infoList
	}

	for _, tag := range tags {
		for _, info := range l.infoMap[tag] {
			if !info.Resolved() {
				notResolved = append(notResolved, info)
			}
		}
	}
	return notResolved
}

//addInfoFor ajoute cette collision aux tags de s
// ajoute la clé de map si n'existe pas
func (l *InfoList) addInfoFor(info *CollisionInfo, tags []string) {
	for _, tag := range tags {
		if _, exists := l.infoMap[tag]; !exists {
			l.infoMap[tag] = []*CollisionInfo{info}
		} else {
			l.infoMap[tag] = append(l.infoMap[tag], info)
		}
	}
}

//Reset reset une InfoList avant nouvel usage
func (l *InfoList) Reset() {
	l.infoList = nil //TODO: vérifier que bien nécessaire
	l.infoMap = nil
	l.infoList = []*CollisionInfo{}
	l.infoMap = map[string][]*CollisionInfo{}

}

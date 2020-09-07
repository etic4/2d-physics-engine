package physique

//Shape type de la forme
type Shape interface {
	Pos() Vec2
	SetPos(Vec2)
	UpdatePos()
	Width() float64
	Height() float64
	Center() Vec2
	SetCenter(Vec2)
	Velocity() Vec2
	SetVelocity(Vec2)
	MaxVel() Vec2
	SetMaxVel(Vec2)
	Accel() Vec2
	SetAccel(Vec2)
	MaxAccel() Vec2
	SetMaxAccel(Vec2)
	Gravity() Vec2
	SetGravity(Vec2)
	IsGrounded() bool
	SetGrounded(bool)
	IsStatic() bool
	SetStatic(bool)
	IsSolid() bool
	SetSolid(bool)
	Friction() float64
	SetFriction(float64)
	InvMass() float64
	Elasticity() float64
	SetElasticity(float64)
	ShapeName() string
	Name() string
	Tags() []string
	SetTags([]string)
}

// BasicShape Une forme générique
type BasicShape struct {
	Kind       Shape
	pos        Vec2
	velocity   Vec2
	accel      Vec2
	gravity    Vec2
	maxVel     Vec2 //vitesse maximale
	maxAccel   Vec2 //accélération maximale
	grounded   bool
	static     bool
	solid      bool
	mass       float64
	invMass    float64
	elasticity float64
	friction   float64
	name       string
	tags       []string
}

//Pos retourne la position de la forme
func (s *BasicShape) Pos() Vec2 {
	return s.pos
}

//UpdatePos met à jour position avec la vitesse
// La fonction SetPos est définie sur les shape parce que
// Circle doit mettre à jour le centre
func (s *BasicShape) UpdatePos() {
	// clamp accel
	s.clampAccel()

	//ajout accélération et gravité
	s.velocity = s.velocity.Add(s.accel).Add(s.gravity)

	s.clampVelocity()

	s.Kind.SetPos(s.Pos().Add(s.Velocity()))

	//reset ground state
	s.SetGrounded(false)

}

//clampVelocity restreint la vitesse à -maxVel, maxVel
// Si maxVel est à 0,0 la vitesse n'est pas restreinte
func (s *BasicShape) clampVelocity() {
	if s.maxVel.X > 0 && s.maxVel.Y > 0 {
		s.velocity.X = Max(-s.maxVel.X, Min(s.velocity.X, s.maxVel.X))
		s.velocity.Y = Max(-s.maxVel.Y, Min(s.velocity.Y, s.maxVel.Y))
	}
}

//Velocity retourne la vitesse de la shape
func (s *BasicShape) Velocity() Vec2 {
	return s.velocity
}

//SetVelocity mets la vitesse à v
func (s *BasicShape) SetVelocity(v Vec2) {
	s.velocity = v
}

//MaxVel retourne la vitesse maximum de la shape
func (s *BasicShape) MaxVel() Vec2 {
	return s.maxVel
}

//SetMaxVel mets la vitesse maximum à v
func (s *BasicShape) SetMaxVel(v Vec2) {
	s.maxVel = v
}

// //AddForces ajoute un vecteur (Vec2.Add()) au vecteur existant (initialisé à 0.0)
// // le vecteur est remis à 0 à chaque thick
// func (s *BasicShape) AddForce(v Vec2) {
// 	s.force.Add(v)
// }

//clampAccel restreint l'accélération à -maxAccel, maxAccel
// Si maxAccel est à 0,0 l'accélération n'est pas restreinte
func (s *BasicShape) clampAccel() {
	if s.maxAccel.X > 0 && s.maxAccel.Y > 0 {
		s.accel.X = Max(-s.maxAccel.X, Min(s.accel.X, s.maxAccel.X))
		s.accel.Y = Max(-s.maxAccel.Y, Min(s.accel.Y, s.maxAccel.Y))
	}
}

//Accel retourne l'accélération
func (s *BasicShape) Accel() Vec2 {
	return s.accel
}

//SetAccel mets l'accélération à 'a'
func (s *BasicShape) SetAccel(a Vec2) {
	s.accel = a
}

//MaxAccel retourne l'accélération maximale
func (s *BasicShape) MaxAccel() Vec2 {
	return s.maxAccel
}

//SetMaxAccel mets l'accélération maximale à 'a'
func (s *BasicShape) SetMaxAccel(a Vec2) {
	s.maxAccel = a
}

//Gravity retourne la gravité
func (s *BasicShape) Gravity() Vec2 {
	return s.gravity
}

//IsGrounded retourne true si la forme est au sol
func (s *BasicShape) IsGrounded() bool {
	return s.grounded
}

//SetGrounded mets le statut au sol à true ou false
func (s *BasicShape) SetGrounded(b bool) {
	s.grounded = b
}

//SetGravity mets la gravité à 'g'
func (s *BasicShape) SetGravity(g Vec2) {
	s.gravity = g
}

//SetMass mets la masse à m
func (s *BasicShape) SetMass(mass float64) {
	s.mass = mass

	if mass <= 0 {
		s.invMass = 0
	} else {
		s.invMass = 1 / mass
	}
}

//InvMass retourne la masse inverse
func (s *BasicShape) InvMass() float64 {
	return s.invMass
}

//Elasticity retourne l'elasticité (le 'coefficient de restitution')
func (s *BasicShape) Elasticity() float64 {
	return s.elasticity
}

//SetElasticity met l'élasticité à e
func (s *BasicShape) SetElasticity(e float64) {
	s.elasticity = e
}

//Friction return friction
func (s *BasicShape) Friction() float64 {
	return s.friction
}

//SetFriction mets friction à f
func (s *BasicShape) SetFriction(f float64) {
	s.friction = f
}

//IsStatic retourne true si la forme est de type statique
func (s *BasicShape) IsStatic() bool {
	return s.static
}

//SetStatic mets la forme à statique ou non
func (s *BasicShape) SetStatic(b bool) {
	s.static = b
}

//IsSolid retourne true si la forme est de type solide, false sinon
func (s *BasicShape) IsSolid() bool {
	return s.solid
}

//SetSolid mets la forme à solide ou non
func (s *BasicShape) SetSolid(b bool) {
	s.solid = b
}

//SetTags attribue la liste des tags
func (s *BasicShape) SetTags(tags []string) {
	s.tags = tags
}

//Tags retourne la liste des tags
func (s *BasicShape) Tags() []string {
	return s.tags
}

//HasTag retourne true si tag est contenu dans la liste des tags
func (s *BasicShape) HasTag(tag string) bool {
	for _, t := range s.tags {
		if t == tag {
			return true
		}
	}
	return false
}

// Rectangle type
type Rectangle struct {
	*BasicShape
	width  float64
	height float64
}

//Width retourne la largeur
func (r *Rectangle) Width() float64 {
	return r.width
}

//Height retourne la hauteur
func (r *Rectangle) Height() float64 {
	return r.height
}

//getMax retourne le max
func (r *Rectangle) getMax() Vec2 {
	return Vec2{r.Pos().X + r.Width(), r.Pos().Y + r.Height()}
}

//Center retourne les coordonées du centre du Rectangle
func (r *Rectangle) Center() Vec2 {
	return Vec2{r.Pos().X + r.Width()/2, r.Pos().Y + r.Height()/2}

}

//SetCenter Positionne un rectangle par son centre
func (r *Rectangle) SetCenter(c Vec2) {
	r.SetPos(Vec2{c.X - r.Width()/2, c.Y - r.Height()/2})
}

//SetPos mets la position de l'objet à p
func (r *Rectangle) SetPos(p Vec2) {
	r.BasicShape.pos = p
}

//ShapeName retourne le nom de la forme
func (r *Rectangle) ShapeName() string {
	return "Rectangle"
}

//SetName met name à n
func (r *Rectangle) SetName(n string) {
	r.name = n
}

//Name retourne le nom de la forme
func (r *Rectangle) Name() string {
	return r.name
}

//NewRectangle Crée un rectangle
func NewRectangle(pos Vec2, width float64, height float64) *Rectangle {
	rect := &Rectangle{width: width, height: height}
	rect.BasicShape = &BasicShape{Kind: rect, pos: pos}
	rect.SetName(UUID())
	rect.SetSolid(true)
	return rect
}

//Circle un cercle
type Circle struct {
	*BasicShape
	center Vec2
	radius float64
}

//ShapeName retourne le nom de la forme
func (s *Circle) ShapeName() string {
	return "Circle"
}

//SetName met name à n
func (s *Circle) SetName(n string) {
	s.name = n
}

//Name retourne le nom de la forme
func (s *Circle) Name() string {
	return s.name
}

//Center retourne les coordonées du centre du Rectangle
func (s *Circle) Center() Vec2 {
	return s.center
}

//SetCenter positionne un cercle par son centre
func (s *Circle) SetCenter(c Vec2) {
	s.center = c
	s.SetPos(s.center.SubScalar(s.radius))
}

//SetPos mets la position à p
func (s *Circle) SetPos(p Vec2) {
	s.BasicShape.pos = p
	s.center = p.AddScalar(s.radius)
}

//getMax retourne le max
func (s *Circle) getMax() Vec2 {
	return s.center.AddScalar(s.radius)
}

//Radius retourne le rayon du cercle
func (s *Circle) Radius() float64 {
	return s.radius
}

//Width retourne la largeur
func (s *Circle) Width() float64 {
	return s.radius * 2
}

//Height retourne la hauteur
func (s *Circle) Height() float64 {
	return s.Width()
}

//NewCircle créé un nouveau cercle
func NewCircle(center Vec2, radius float64) *Circle {
	circ := &Circle{radius: radius}
	circ.BasicShape = &BasicShape{Kind: circ, pos: center.SubScalar(radius)}
	circ.SetName(UUID())
	circ.SetSolid(true)
	return circ
}

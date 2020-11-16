package shapes


// AABB (axis-aligned bounding box)
type Box3 struct {
	BoundA, BoundB Vec3
}

func NewBox3(vmin, vmax Vec3) Box3  {
	return Box3{BoundA: vmin, BoundB: vmax}
}

func (b Box3) Hit(r *Ray, tMin float64, tMax float64) (bool, *HitRecord)  {
	return false, nil
}
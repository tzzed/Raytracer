package shapes

import (
	"math"
	"math/rand"
	"time"
)

// AABB (axis-aligned bounding box)
type AABB struct {
	Min, Max Vec3
}

func NewAABB(box0, box1 AABB) AABB {
	min := Vec3{
		X: math.Min(box0.Min.X, box1.Min.X),
		Y: math.Min(box0.Min.Y, box1.Min.Y),
		Z: math.Min(box0.Min.Z, box1.Min.Z),
	}
	
	max := Vec3{
		X: math.Max(box0.Min.X, box1.Min.X),
		Y: math.Max(box0.Min.Y, box1.Min.Y),
		Z: math.Max(box0.Min.Z, box1.Min.Z),
	}
	
	return AABB{Min: min, Max: max}
}

func (ab AABB) Hit(r *Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	for a := 0; a < 3; a++ {
		invD := 1.0 / r.Dir.GetAxis(a)
		oa := r.Origin.Vec3().GetAxis(a)
		t0 := (ab.Min.GetAxis(a) - oa) * invD
		t1 := (ab.Max.GetAxis(a) - oa) * invD
		if invD < 0.0 {
			// swap
			t0, t1 = t1, t0
		}
		min := math.Max(t0, tMin)
		max := math.Min(t1, tMax)
		if max <= min {
			return false, nil
		}
	}
	
	return true, nil
}

type BhvNode struct {
	left, right HitTable
	box         AABB
}

func NewBVHNode(start, end int32, tm0, tm1 float64, hl *HitTableList) BhvNode {
	var bn BhvNode
	rand.Seed(time.Now().UnixNano())
	axis := int(rand.Int63n(2))
	
	os := start - end
	byAxis := func(a, b *HitTable) bool {
		aa := *a
		bb := *b
		_, boxA := aa.BoundingBox(0, 0)
		_, boxB := bb.BoundingBox(0, 0)
		return boxA.Min.GetAxis(axis) < boxB.Min.GetAxis(axis)
	}
	
	switch os {
	case 1:
		bn.left = hl.Hits[start]
		bn.right = hl.Hits[start]
	case 2:
		bn.left = hl.Hits[start]
		bn.right = hl.Hits[start+1]
		if hl.BoxCompare(bn.left, bn.right, axis) {
			bn.left = hl.Hits[start+1]
			bn.right = hl.Hits[start]
		}
	default:
		OrderedBy(byAxis).Sort(hl.Hits)
		mid := start + (os / 2)
		bn.left = NewBVHNode(start, mid, tm0, tm1, hl)
		bn.right = NewBVHNode(mid, end, tm0, tm1, hl)
		
	}
	
	_, bl := bn.left.BoundingBox(tm0, tm1)
	_, br := bn.right.BoundingBox(tm0, tm1)
	
	bn.box = NewAABB(*bl, *br)
	return bn
}


func (bn BhvNode) BoundingBox(tm0, tm1  float64) (bool, *AABB) {
	return true, &bn.box
}

func (bn BhvNode) Hit(r *Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	if hit, _ := bn.box.Hit(r, tMin, tMax); !hit {
		return false, nil
	}
	
	hitLeft, hrLeft := bn.left.Hit(r, tMin, tMax)
	hitRight, hrRight := bn.right.Hit(r, tMin, tMax)
	if !hitLeft && !hitRight {
		return false, nil
	}
	
	if hitLeft && hitRight {
		if hrRight.T > hrLeft.T {
			return true, hrLeft
		}
		
		return true, hrRight
	}
	
	if hitLeft {
		return true, hrLeft
	}
	
	if hitRight {
		return true, hrRight
	}
	
	return false, nil
}

package strategy

import (
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/commons/Units"
	"math"
)

// RegionHeight defines the height of a region
const RegionHeight = Units.CourtHeight / 4

// RegionWidth defines the width of a region
const RegionWidth = Units.CourtWidth / 8

// PlayerRegion defines a region based on the left-bottom and top-right coordinates of the region
type PlayerRegion struct {
	CornerA Physics.Point
	CornerB Physics.Point
}

// RegionCode is a cartesian coordinate based on the region division of the court
// The region code is always based on the home team perspective, so the coordinates must be mirrored to apply the
// same location in the field to the away team
type RegionCode struct {
	X int
	Y int
}

// Center finds the central coordinates of the region
func (r RegionCode) Center(place Units.TeamPlace) Physics.Point {
	center := Physics.Point{
		PosX: (r.X * RegionWidth) + (RegionWidth / 2),
		PosY: (r.Y * RegionHeight) + (RegionHeight / 2),
	}
	if place == Units.AwayTeam {
		center = MirrorCoordsToAway(center)
	}
	return center
}

// ForwardRightCorner finds the point of the region at the right edge that is closer to the attack field
func (r RegionCode) ForwardRightCorner(place Units.TeamPlace) Physics.Point {
	fr := Physics.Point{
		PosX: (r.X + 1) * RegionWidth,
		PosY: r.Y * RegionHeight,
	}
	if place == Units.AwayTeam {
		fr = MirrorCoordsToAway(fr)
	}
	return fr
}

// ForwardLeftCorner finds the point of the region at the left edge that is closer to the attack field
func (r RegionCode) ForwardLeftCorner(place Units.TeamPlace) Physics.Point {
	fl := Physics.Point{
		PosX: (r.X + 1) * RegionWidth,
		PosY: (r.Y + 1) * RegionHeight,
	}
	if place == Units.AwayTeam {
		fl = MirrorCoordsToAway(fl)
	}
	return fl
}

// Forwards finds the next region towards to the attack field, or return itself when there is no region in front of it
func (r RegionCode) Forwards() RegionCode {
	if r.X == 7 {
		return r
	}
	return RegionCode{
		X: r.X + 1,
		Y: r.Y,
	}
}

// Backwards finds the next region towards to the defense field, or return itself when there is no region in behind of it
func (r RegionCode) Backwards() RegionCode {
	if r.X == 0 {
		return r
	}
	return RegionCode{
		X: r.X - 1,
		Y: r.Y,
	}
}

// Left finds the region in the left side of this region, or return itself when there is no region there
func (r RegionCode) Left() RegionCode {
	if r.Y == 3 {
		return r
	}
	return RegionCode{
		X: r.X,
		Y: r.Y + 1,
	}
}

// Right finds the region in the right side of this region, or return itself when there is no region there
func (r RegionCode) Right() RegionCode {
	if r.Y == 0 {
		return r
	}
	return RegionCode{
		X: r.X,
		Y: r.Y - 1,
	}
}

// ChessDistanceTo calculates what is the chess distance (steps towards any direction, even diagonal) between these regions
func (r RegionCode) ChessDistanceTo(b RegionCode) int {
	return int(math.Max(
		math.Abs(float64(r.X-b.X)),
		math.Abs(float64(r.Y-b.Y)),
	))
}

func GetRegionCode(a Physics.Point, place Units.TeamPlace) RegionCode {
	if place == Units.AwayTeam {
		a = MirrorCoordsToAway(a)
	}
	cx := float64(a.PosX / RegionWidth)
	cy := float64(a.PosY / RegionHeight)
	return RegionCode{
		X: int(math.Min(cx, 7)),
		Y: int(math.Min(cy, 3)),
	}
}

// Invert the coords X and Y as in a mirror to found out the same position seen from the away team field
// Keep in mind that all coords in the field is based on the bottom left corner!
func MirrorCoordsToAway(coords Physics.Point) Physics.Point {
	return Physics.Point{
		PosX: Units.CourtWidth - coords.PosX,
		PosY: Units.CourtHeight - coords.PosY,
	}
}

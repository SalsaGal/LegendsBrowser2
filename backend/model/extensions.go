package model

import "fmt"

func (e *Entity) Position(id int) *EntityPosition {
	for _, p := range e.EntityPosition {
		if p.Id_ == id {
			return p
		}
	}
	return &EntityPosition{Name_: "UNKNOWN POSITION"}
}

func (p *EntityPosition) GenderName(hf *HistoricalFigure) string {
	if hf.Female() && p.NameFemale != "" {
		return p.NameFemale
	} else if hf.Male() && p.NameMale != "" {
		return p.NameMale
	} else {
		return p.Name_
	}
}

func (hf *HistoricalFigure) Female() bool {
	return hf.Sex == 0 || hf.Caste == "FEMALE"
}

func (hf *HistoricalFigure) Male() bool {
	return hf.Sex == 1 || hf.Caste == "MALE"
}

type HistoricalEventDetails interface {
	RelatedToEntity(int) bool
	RelatedToHf(int) bool
	Html() string
	Type() string
}

type HistoricalEventCollectionDetails interface {
}

func containsInt(list []int, id int) bool {
	for _, v := range list {
		if v == id {
			return true
		}
	}
	return false
}

var world *DfWorld

func artifact(id int) string {
	if x, ok := world.Artifacts[id]; ok {
		return fmt.Sprintf(`<a href="/artifact/%d">%s</a>`, x.Id(), x.Name())
	}
	return "UNKNOWN ARTIFACT"
}

func entity(id int) string {
	if x, ok := world.Entities[id]; ok {
		return fmt.Sprintf(`<a href="/entity/%d">%s</a>`, x.Id(), x.Name())
	}
	return "UNKNOWN ENTITY"
}

func hf(id int) string {
	if x, ok := world.HistoricalFigures[id]; ok {
		return fmt.Sprintf(`<a href="/hf/%d">%s</a>`, x.Id(), x.Name())
	}
	return "UNKNOWN HISTORICAL FIGURE"
}

func site(id int, prefix string) string {
	if x, ok := world.Sites[id]; ok {
		return fmt.Sprintf(`%s <a href="/site/%d">%s</a>`, prefix, x.Id(), x.Name())
	}
	return "UNKNOWN SITE"
}

func structure(siteId, structureId int) string {
	if x, ok := world.Sites[siteId]; ok {
		if y, ok := x.Structures[structureId]; ok {
			return fmt.Sprintf(`<a href="/site/%d/structure/%d">%s</a>`, siteId, structureId, y.Name())
		}
	}
	return "UNKNOWN STRUCTURE"
}

func region(id int) string {
	if x, ok := world.Regions[id]; ok {
		return fmt.Sprintf(`<a href="/region/%d">%s</a>`, x.Id(), x.Name())
	}
	return "UNKNOWN REGION"
}

package mapping

func (m *Mapping) NodeTagFilter() *TagFilter {
	mappings := make(map[string]map[string][]string)
	m.mappings("point", mappings)
	tags := make(map[string]bool)
	m.extraTags("point", tags)
	return &TagFilter{mappings, tags}
}

func (m *Mapping) WayTagFilter() *TagFilter {
	mappings := make(map[string]map[string][]string)
	m.mappings("linestring", mappings)
	m.mappings("polygon", mappings)
	tags := make(map[string]bool)
	m.extraTags("linestring", tags)
	m.extraTags("polygon", tags)
	return &TagFilter{mappings, tags}
}

func (m *Mapping) RelationTagFilter() *RelationTagFilter {
	mappings := make(map[string]map[string][]string)
	m.mappings("linestring", mappings)
	m.mappings("polygon", mappings)
	tags := make(map[string]bool)
	m.extraTags("linestring", tags)
	m.extraTags("polygon", tags)
	tags["type"] = true // do not filter out type tag
	return &RelationTagFilter{TagFilter{mappings, tags}}
}

type TagFilter struct {
	mappings  map[string]map[string][]string
	extraTags map[string]bool
}

type RelationTagFilter struct {
	TagFilter
}

func (f *TagFilter) Filter(tags map[string]string) bool {
	foundMapping := false
	for k, v := range tags {
		values, ok := f.mappings[k]
		if ok {
			if _, ok := values["__any__"]; ok {
				foundMapping = true
				continue
			} else if _, ok := values[v]; ok {
				foundMapping = true
				continue
			} else if _, ok := f.extraTags[k]; !ok {
				delete(tags, k)
			}
		} else if _, ok := f.extraTags[k]; !ok {
			delete(tags, k)
		}
	}
	if foundMapping {
		return true
	} else {
		return false
	}
}

func (f *RelationTagFilter) Filter(tags map[string]string) bool {
	if t, ok := tags["type"]; ok {
		if t != "multipolygon" && t != "boundary" && t != "land_area" {
			return false
		}
	} else {
		return false
	}
	f.TagFilter.Filter(tags)
	// always return true here since we found a matching type
	return true
}

package commands

// splitTags separates the provided tags into two slices: one for tags without the '!' prefix (with),
// and another for tags that have the '!' prefix (without).
func splitTags(tags []string) ([]string, []string) {
	var with, without []string
	for _, tag := range tags {
		if tag[0] == '!' {
			without = append(without, tag[1:])
		} else {
			with = append(with, tag)
		}
	}
	return with, without
}

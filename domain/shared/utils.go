package shared

func RemoveDuplicateMembers(members []int) []int {
	chatMembersMap := make(map[int]int)
	for _, member := range members {
		chatMembersMap[member] = 1
	}
	members = make([]int, 0, len(chatMembersMap))
	for member := range chatMembersMap {
		members = append(members, member)
	}
	return members
}

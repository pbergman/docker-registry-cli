package api

import "sort"

type RepositoryList []*RepositoryInfo

func (l *RepositoryList) Get(name string) *RepositoryInfo {
	for _, info := range *l {
		if info.Name == name {
			return info
		}
	}
	info := &RepositoryInfo{Name: name}
	*l = append(*l, info)
	return info
}

func (l RepositoryList) Sort() {
	sort.Sort(l)
	for _, info := range l {
		sort.Strings(info.Tags)
	}
}

func (l RepositoryList) Len() int {
	return len(l)
}

func (l RepositoryList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l RepositoryList) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

type RepositoryInfo struct {
	Name string
	Tags []string
}

func (l *RepositoryInfo) AddTag(tag string) {
	l.Tags = append(l.Tags, tag)
}

func GetList() *RepositoryList {

	list := new(RepositoryList)

	for _, repository := range GetRepositories().Images {
		for _, tag := range GetTags(repository).Tags {
			if manifest := GetManifest(repository, tag, false); manifest != nil {
				info := list.Get(repository)
				info.AddTag(tag)
			}
		}
	}

	return list
}

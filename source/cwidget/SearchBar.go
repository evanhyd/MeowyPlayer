package cwidget

import (
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Titler interface {
	Title() string
}

type SearchBar[T Titler] struct {
	widget.BaseWidget
	dropDown *DropDown
	entry    *widget.Entry
	tools    []fyne.CanvasObject
	cmp      func(T, T) int
	filter   func(Titler) bool
	onUpdate func()
}

func NewSearchBar[T Titler](onUpdate func(), tools ...fyne.CanvasObject) *SearchBar[T] {
	s := &SearchBar[T]{
		dropDown: NewDropDown(theme.MenuIcon()),
		entry:    widget.NewEntry(),
		tools:    tools,
		cmp:      func(T, T) int { return -1 },
		filter:   func(Titler) bool { return true },
		onUpdate: onUpdate,
	}
	s.entry.OnChanged = s.setFilter

	s.ExtendBaseWidget(s)
	return s
}

func (s *SearchBar[T]) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, s.dropDown, container.NewHBox(s.tools...), s.entry))
}

func (s *SearchBar[T]) AddComparator(title string, icon fyne.Resource, cmp func(T, T) int) {
	s.dropDown.Add(title, icon, func() { s.setComparator(cmp) })
}

func (s *SearchBar[T]) Select(i int) {
	s.dropDown.Select(i)
}

func (s *SearchBar[T]) setComparator(cmp func(T, T) int) {
	s.cmp = cmp
	s.onUpdate()
}

func (s *SearchBar[T]) setFilter(title string) {
	s.filter = func(t Titler) bool { return strings.Contains(strings.ToLower(t.Title()), strings.ToLower(title)) }
	s.onUpdate()
}

func (s *SearchBar[T]) Query(data []T) []T {
	//lazy removal, swapping to the back
	i, j := 0, len(data)
	for i < j {
		if s.filter(data[i]) {
			i++
		} else {
			j--
			data[i], data[j] = data[j], data[i]
		}
	}

	slices.SortStableFunc(data[:j], s.cmp)
	return data[:j]
}

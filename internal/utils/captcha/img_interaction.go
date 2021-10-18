package captcha

import (
	"errors"
	"fmt"
	"github.com/LittleBenx86/Benlog/internal/utils/collection"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Item struct {
	Font          string     `json:"font"`
	Coordinate    Coordinate `json:"coordinate"`
	Selected      bool       `json:"selected"`
	ValidateOrder int        `json:"validate_order"`
}

type Group struct {
	ValidateCount int    `json:"validate_count"` // <= ValidateMode * ValidateMode
	ValidateMode  int    `json:"validate_mode"`  // 3: 3*3, 4: 4*4
	Items         []Item `json:"items"`
	ImageSize     int    `json:"image_size"` // image width and height, 64 * 64, this is a reference value. If it is zero, the caller can select other sizes
}

type ImageInteraction struct {
	AccessKey      string `json:"access_key"`
	ImageInfoGroup Group  `json:"image_info_group"`
	Checksum       string `json:"-"`
}

func (i *Item) SetAsSelected(order int) {
	i.Selected = true
	i.ValidateOrder = order
}

func (i *Item) IsValid(mode int) bool {
	if len(i.Font) > 1 {
		return false
	}

	if i.Coordinate.X > (mode-1) ||
		i.Coordinate.Y > (mode-1) {
		return false
	}
	return true
}

func (g *Group) IsValid() bool {

	if g.ValidateCount > g.ValidateMode*g.ValidateMode {
		return false
	}

	if len(g.Items) < g.ValidateCount {
		return false
	}

	for _, i := range g.Items {
		if !i.IsValid(g.ValidateMode) {
			return false
		}
	}

	return true
}

var (
	RandomFontLib []string
)

func getRandomFonts(indexes []int) []string {
	s := make([]string, 0)
	for _, i := range indexes {
		s = append(s, RandomFontLib[i])
	}
	return s
}

func newCoordinates(validateMode int, shuffleFn func(interface{}) error) ([]Coordinate, error) {
	c := make([]Coordinate, 0)
	for i := 0; i < validateMode; i++ {
		for j := 0; j < validateMode; j++ {
			c = append(c, Coordinate{X: i, Y: j})
		}
	}

	if err := shuffleFn(c); err != nil {
		return nil, err
	}
	return c, nil
}

func filterBySelected(unfilteredItems []Item) []Item {
	s := make([]Item, 0)
	for _, item := range unfilteredItems {
		if item.Selected {
			s = append(s, item)
		}
	}
	return s
}

func sortByOrder(selectedItems []Item) []Item {
	sort.Slice(selectedItems, func(i, j int) bool {
		if selectedItems[i].ValidateOrder > selectedItems[j].ValidateOrder {
			return false
		}
		return true
	})
	return selectedItems
}

func BuildChecksumRequired(unfilteredItems []Item) string {
	// [font:order:x:y]
	formatter := "[%s:%d]"
	var builder strings.Builder
	builder.Grow(50)
	cpy := sortByOrder(filterBySelected(unfilteredItems))
	for _, item := range cpy {
		builder.WriteString(fmt.Sprintf(formatter, item.Font, item.ValidateOrder))
	}
	return builder.String()
}

func NewImageInteraction(
	accessKey string,
	validateMode int,
	validateCount int,
	unrepeatableRandFn func(start, end, count int) ([]int, error),
	hashFn func(string) (string, error),
) (*ImageInteraction, error) {
	fontLibLen, totalImageCount := len(RandomFontLib), validateMode*validateMode

	// generate the random coordinates, the count of these are equal to total image count
	clib, err := newCoordinates(validateMode, collection.SliceShuffle)
	if err != nil {
		return nil, err
	}

	// generate the random fonts, the count of these are equal to total image count
	fontRandIndexes := make([]int, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < totalImageCount; i++ {
		fontRandIndexes = append(fontRandIndexes, r.Intn(fontLibLen))
	}
	fs := getRandomFonts(fontRandIndexes)

	// generate items by upper fonts and coordinates
	items := make([]Item, 0)
	for i := 0; i < totalImageCount; i++ {
		items = append(items, Item{
			Font:          fs[i],
			Coordinate:    clib[i],
			Selected:      false,
			ValidateOrder: -1,
		})
	}

	// check data is valid
	g := Group{
		Items:         items,
		ValidateCount: validateCount,
		ValidateMode:  validateMode,
		ImageSize:     0,
	}
	if !g.IsValid() {
		return nil, errors.New("an incorrect image interaction group created")
	}

	// select n items as key targets will be validated
	selectedIndexes, err := unrepeatableRandFn(0, totalImageCount-1, validateCount)
	if err != nil {
		return nil, err
	}
	for i := 0; i < validateCount; i++ {
		items[selectedIndexes[i]].SetAsSelected(i)
	}

	// hash checksum
	checksum, err := hashFn(BuildChecksumRequired(g.Items))
	if err != nil {
		return nil, errors.New("checksum generation failed")
	}

	return &ImageInteraction{
		AccessKey:      accessKey,
		ImageInfoGroup: g,
		Checksum:       checksum,
	}, nil
}

// Package names provides a random name generator.
package names

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type NameGenerator struct {
	rnd *rand.Rand
	use map[string]struct{}
	mu  sync.Mutex
}

func NewNameGenerator() *NameGenerator {
	return &NameGenerator{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
		use: make(map[string]struct{}),
	}
}

// GetUserName generates an unused user name.
func (ng *NameGenerator) GetUserName() string {
	return ng.generate(ng.generateUserName)
}

func (ng *NameGenerator) generateUserName() string {
	i := 'a' + rune(ng.rnd.Int()%len(people))
	j := ng.rnd.Int() % len(people[i])
	k := ng.rnd.Int() % len(people[i])
	return fmt.Sprintf("%s%s", strings.Title(people[i][j]), strings.Title(people[i][k]))
}

// GetRoomName generates an unused room name.
func (ng *NameGenerator) GetRoomName() string {
	return ng.generate(ng.generateRoomName)
}

func (ng *NameGenerator) generateRoomName() string {
	i := ng.rnd.Int() % len(adjectives)
	j := ng.rnd.Int() % len(animals)
	return fmt.Sprintf("%s%s", strings.Title(adjectives[i]), strings.Title(animals[j]))
}

// Release allows the given name to be reused.
func (ng *NameGenerator) Release(name string) {
	ng.mu.Lock()
	defer ng.mu.Unlock()

	if _, ok := ng.use[name]; ok {
		delete(ng.use, name)
	}
}

func (ng *NameGenerator) generate(f func() string) string {
	ng.mu.Lock()
	defer ng.mu.Unlock()

	name := f()
	for _, ok := ng.use[name]; ok; _, ok = ng.use[name] {
		name = f()
	}
	ng.use[name] = struct{}{}
	return name
}

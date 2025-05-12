package main

import (
	"math"
	"testing"
	"unicode/utf8"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		if utf8.RuneCountInString(name) > 42 {
			return
		}
		for i := range person.name {
			person.name[i] = 0
		}

		for i, letter := range []rune(name) {
			if letter < 'A' || letter > 'z' {
				// или continue; или 2 раза проходить цикл, чтобы проверить все символы сначала; или через временный массив
				return
			}
			person.name[i] = byte(letter)
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		if !validateCoordinate(x) || !validateCoordinate(y) || !validateCoordinate(z) {
			return
		}
		person.coordX = int32(x)
		person.coordY = int32(y)
		person.coordZ = int32(z)
	}
}

func validateCoordinate(coord int) bool {
	if coord > math.MaxInt32 || coord < math.MinInt32 {
		return false
	}
	return true
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		if gold > math.MinInt32 && gold < 0 {
			return // хотя можно переписывать в max или min
		}
		person.gpStats3 = (person.gpStats3 & 0x80000000) | uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		if mana > 1_000 || mana < 0 {
			return // хотя можно переписывать в max или min
		}
		person.gpStats2 = (person.gpStats2 & 0xFFFFFC00) | uint32(mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		if health > 1_000 || health < 0 {
			return // хотя можно переписывать в max или min
		}
		person.gpStats2 = (person.gpStats2 & 0xFFF003FF) | uint32(health)<<10
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		if respect > 10 || respect < 0 {
			return // хотя можно переписывать в max или min
		}
		person.gpStats2 = (person.gpStats2 & 0xFF0FFFFF) | uint32(respect)<<20
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		if strength > 10 || strength < 0 {
			return
		}
		person.gpStats2 = (person.gpStats2 & 0xF0FFFFFF) | uint32(strength)<<24
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		if experience > 10 || experience < 0 {
			// next level?
			return
		}
		person.gpStats2 = (person.gpStats2 & 0x0FFFFFFF) | uint32(experience)<<28
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		if level > 10 || level < 0 {
			return
		}
		person.gpStats = (person.gpStats & 0xF0) | byte(level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.gpStats3 = (person.gpStats3 & 0x7FFFFFFF) | 1<<31
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.gpStats = (person.gpStats & 0b11101111) | 1<<4
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.gpStats = (person.gpStats & 0b11011111) | 1<<5
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		if personType >= GamePersonTypeEndPlug {
			return
		}
		person.gpStats = (person.gpStats & 0b00111111) | byte(personType)<<6
	}
}

const (
	_ = iota // чтобы не было багов с нулевым значением
	BuilderGamePersonType
	BlacksmithGamePersonType
	WarriorGamePersonType
	GamePersonTypeEndPlug
)

type GamePerson struct {
	/*
	   LIMIT = 64 byte
	   **Данные пользователя:**

	   - __Имя пользователя__ \[0…42\] символов латиницы | 42 byte
	     - нельзя ссылаться на символы строки по указателю (*нужно мапить символы строки в объект структуры, чтобы они находились рядом с другими данными*)

	   - __Координата по оси X__ \[-2_000_000_000…2_000_000_000\] значений -> INT32 (32 bit) (4 byte)	|
	   - __Координата по оси Y__ \[-2_000_000_000…2_000_000_000\] значений -> INT32 (32 bit) (4 byte)	| 12 byte
	   - __Координата по оси Z__ \[-2_000_000_000…2_000_000_000\] значений -> INT32 (32 bit) (4 byte)	|

	   - __Магическая сила (мана)__ \[0…1000\] значений -> 10 bit	|
	   - __Здоровье__ \[0…1000\] значений -> 10 bit					|
	   - __Уважение__ \[0…10\] значений -> 4 bit					|	10+10+12=32 bit (4 byte)
	   - __Сила__ \[0…10\] значений -> 4 bit						|
	   - __Опыт__ \[0…10\] значений -> 4 bit						|


	   - __Золото__ \[0…2_000_000_000\] значений -> UINT32 (31 bit)	|	(4 byte)
	   - __Есть ли у игрока дом__ \[true/false\] значения -> 1 bit	|

	   - __Уровень__ \[0…10\] значений -> 4 bit							|
	   - __Есть ли у игрока оружие__ \[true/false\] значения -> 1 bit	|
	   - __Есть ли у игрока семья__ \[true/false\] значения -> 1 bit	| 1 byte
	   - __Тип игрока__ \[строитель/кузнец/воин\] значения -> 2 bit		|
	*/
	name                   [42]byte
	gpStats                byte // lvl, gun, family, gpType
	coordX, coordY, coordZ int32
	gpStats2               uint32 // mana, hp, resp, strength, exp
	gpStats3               uint32 // money, house
}

func NewGamePerson(options ...Option) GamePerson {
	newGamePerson := GamePerson{}
	for i := range options {
		options[i](&newGamePerson)
	}
	return newGamePerson
}

func (p *GamePerson) Name() string {
	lastIdx := len(p.name)
	for i := range p.name {
		if p.name[i] == 0 {
			lastIdx = i
			break
		}
	}
	return string(p.name[:lastIdx])
}

func (p *GamePerson) X() int {
	return int(p.coordX)
}

func (p *GamePerson) Y() int {
	return int(p.coordY)
}

func (p *GamePerson) Z() int {
	return int(p.coordZ)
}

func (p *GamePerson) Gold() int {
	return int(p.gpStats3 & 0x7FFFFFFF)
}

func (p *GamePerson) Mana() int {
	return int(p.gpStats2 & 0x3FF)
}

func (p *GamePerson) Health() int {
	return int((p.gpStats2 & 0xFFC00) >> 10)
}

func (p *GamePerson) Respect() int {
	return int((p.gpStats2 & 0xF00000) >> 20)
}

func (p *GamePerson) Strength() int {
	return int((p.gpStats2 & 0xF000000) >> 24)
}

func (p *GamePerson) Experience() int {

	return int((p.gpStats2 & 0xF0000000) >> 28)
}

func (p *GamePerson) Level() int {
	return int(p.gpStats & 0x0F)
}

func (p *GamePerson) HasHouse() bool {
	return (p.gpStats3 & 0x80000000) == 0x80000000
}

func (p *GamePerson) HasGun() bool {
	return (p.gpStats & 0b00010000) == 0b00010000
}

func (p *GamePerson) HasFamily() bool {
	return (p.gpStats & 0b00100000) == 0b00100000
}

func (p *GamePerson) Type() int {
	return int((p.gpStats & 0b11000000) >> 6)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}

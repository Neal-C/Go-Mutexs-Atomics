package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Player struct {
	healthMutex sync.RWMutex
	health int32
}

func NewPlayer() *Player {
	return &Player{
		health: 100,
	}
}


func (p *Player) getHealth() int32 {
	p.healthMutex.RLock();
	defer p.healthMutex.RUnlock();
	return p.health;
}

func (p *Player) takeDamage(damage int32){
	p.healthMutex.Lock();
	defer p.healthMutex.Unlock();
	p.health -= damage;
}

func (p *Player) getAtomicHealth() int {
	// int32(p.health)
	return int(atomic.LoadInt32(&p.health));
}

func (p *Player) takeAtomicDamage(damage int32){
	health := p.getAtomicHealth();
	atomic.StoreInt32(&p.health, int32(health-int(damage)));
}

func startUILoop(player *Player) {
	ticker := time.NewTicker(time.Second);

	for {
		fmt.Printf("player health: %d \r ", player.getHealth());
		<- ticker.C
	}
}

func startGameLoop(player *Player){
	ticker := time.NewTicker(time.Second);
	for {

		damage := int32(rand.Intn(40));
		player.takeDamage(damage);

		if (damage >= player.health){
			fmt.Println("oh no, Game over 😥");
			break;
		}
		player.health -= damage;
		if (player.health <= 0 ){
			fmt.Println("oh no, Game over 😥");
			break;
		}
		<- ticker.C
	}
}

func main() {
	player := NewPlayer();
	go startUILoop(player);
	startGameLoop(player);
}
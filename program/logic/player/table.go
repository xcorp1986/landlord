package player

import (
	"chessSever/program/logic/game/games"
	"sync"
	"chessSever/program/logic/game/poker"
	"chessSever/program/logic/game"
	"strconv"
	"fmt"
	"errors"
	"time"
)

/*
	定义游戏桌对象
*/
type Table struct {
	Key          string     				//桌子key,用于从room索引中查找桌子
	Players      [...]*Player  				//玩家数组
	Game         games.IGame 				//该桌玩的游戏
	sync.RWMutex            				//操作playNum以及player时加锁
	CurrPokerCards []*poker.PokerCard  		//当前出的牌
	CurrPalyerIndex int 					//当前出牌的玩家数组index
	IsPlaying    bool                       //是否正在游戏中
}
//创建桌子
func newTable(player *Player, gameName string) *Table {

	currGame := game.GetGame(gameName)
	table := Table{
		Game: currGame,
		Key:  "table" + strconv.Itoa(time.Now().Nanosecond()),//桌子的key要保证唯一且好找，所以用时间戳，
		Players:[currGame.GetPlayerNum()]*Player{},
		IsPlaying:false,
	}
	fmt.Println("创建新桌子"+"table" + strconv.Itoa(player.Id))
	//桌子加入房间
	table.joinRoom()
	//将创建者加入桌子
	table.addPlayer(player)
	return &table
}
//加入房间
func (t *Table) joinRoom() {
	getRoom().addTable(t.Key, t)
}
//销毁桌子
func (t *Table) destory() {
	t.Lock()
	if len(t.Players) >= 0 {
		for _, p := range t.Players {
			p.LeaveTable()
		}
	}
	getRoom().removeTable(t.Key)
	fmt.Println("桌子"+t.Key+"销毁")
	t.Unlock()
}
//增加玩家
func (t *Table) addPlayer(player *Player) error {
	t.Lock()
	defer t.Unlock()
	for i,p := range t.Players{
		if(p == nil){
			p = player
			return nil
		}else{
			if i == len(t.Players)-1 {
				return errors.New("该卓玩家已经满了")
			}
		}
	}
}
//移除玩家
func (t *Table) removePlayer(player *Player) {
	t.Lock()
	for i, p := range t.Players {
		if p == player {
			t.Players[i] = nil
			break
		}
	}
	fmt.Println("桌子"+t.Key+"移除玩家"+strconv.Itoa(player.Id)+"，当前玩家数是"+strconv.Itoa(len(t.Players)))
	t.Unlock()
}

func (t *Table) userReady(){
	userAllReady := false
	for _,p := range t.Players{
		if p != nil && p.IsReady{
			userAllReady = true
		}else{
			userAllReady = false
		}
	}
	//用户都准备好了，则发牌
	if userAllReady {
		t.Game.DealCards()
		t.dealCards()
	}
}

func (t *Table) dealCards(){

}




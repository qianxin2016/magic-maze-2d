<template>
  <div class="HolyGrail">
    <link href="https://cdn.jsdelivr.net/npm/animate.css@3.5.1" rel="stylesheet" type="text/css">
    <header></header>
    <div class="HolyGrail-body">
      <div class="HolyGrail-content">
        <playground>
          <grid>
            <gridcell v-for="cell in board" :key="cell" @click="step(cell)">
              <transition name="gem-transition" enter-active-class="animated flash">
                <div v-show='cell.showWall' class="cell-wall" :style="cell.wall" />
              </transition>
              <transition name="gem-transition" leave-active-class="animated zoomOutRight">
                <img v-show="cell.show" class="cell-img" :src="cell.src" />
              </transition>
              <transition name="player-transition" :enter-active-class="enterAnim" :leave-active-class="exitAnim">
                <img v-show="is_foothold(cell.id)" class="cell-player" :src="playerImg" />
              </transition>
            </gridcell>
          </grid>
        </playground>
        <controller>
          <bonus>
            <div class="tron-icon">
              <img class="tron-icon-img" src="../assets/tron.png" />
            </div>
            <div class="total-bonus">
              <div v-for="digit in bonus_digits" class="bonus-digit">{{digit}}</div>
            </div>
          </bonus>
          <collection>
            <div v-for="player in players" :key="player" class="score">
              <div class="avatar">
                <img class="cell-img" :src="player.avatar" />
              </div>
              <div class="gem-collection">
                <div v-for="gem in player.gems" class="gem-item">
                  <transition appear appear-active-class="animated flash">
                    <img class="cell-img" :src="gem" />
                  </transition>
                </div>
              </div>
            </div>
          </collection>
          <dice-panel>
            <img class="dice" :src="diceImg" />
            <div class="dice-points">
              <div v-for="p in dice_points" :key="p" class="circle">
                <transition enter-active-class="animated fadeIn" leave-active-class="animated zoomOut">
                  <img v-show="p.visible" class="circle-img" src="../assets/dice/circle.png" />
                </transition>
              </div>
            </div>
          </dice-panel>
          <panel>
            <div class="panel-btn roll-dice" @click="roll_dice"> Roll Dice </div>
            <div class="panel-btn quit" @click="quit"> Quit </div>
          </panel>
        </controller>
      </div>
      <div class="HolyGrail-left"></div>
      <div class="HolyGrail-right"></div>
    </div>
    <footer> ©Flamingo 2018. All Rights Reserved. </footer>
    <Dialog @close="close_dialog" :show="showDialog" :title="dialogTitle" :button_text="dialogButtonText" />
  </div>
</template>
<script>
import { Group, Cell } from 'vux'
import Dialog from './Dialog.vue'
import contract from '../utils/contract'
import backend from '../utils/backend'

export default {
  components: {
    Group,
    Cell,
    Dialog
  },
  mounted () {
    setTimeout(() => {
      this.login()
    }, 1000)
  },
  data () {
    return {
      // constants
      wallAnims: ['border-left: 5px solid red',
        'border-right: 5px solid red',
        'border-top: 5px solid red',
        'border-bottom: 5px solid red'
      ],
      moveAnims: ['animated slideInRight',
        'animated slideInLeft',
        'animated slideInUp',
        'animated slideInDown'
      ],
      // variables
      board: Array.apply(null, { length: 49 })
        .map(function (_, index) {
          return {
            id: index,
            src: (index % 2 === 0 && index % 14 >= 8)
              ? require('../assets/diamond/' + (Math.floor(Math.random() * 8 + 0.5)) + '.png') : '',
            show: true,
            wall: '',
            showWall: false
          }
        }),
      enterAnim: 'animated slideInRight',
      exitAnim: 'animated fadeOut',
      playerId: this.$route.query.playerId,
      playerImg: this.$route.query.playerId == 0 ?
                  require('../assets/player1.png') : require('../assets/player2.png'),
      foothold: this.$route.query.playerId == 0 ? 48 : 0,
      players: Array.apply(null, { length: 2 })
        .map(function (_, index) {
          return {
            id: index,
            avatar: require('../assets/player' + (index + 1).toString() + '.png'),
            gems: []
          }
        }),
      roomId: this.$route.query.id,
      roomName: this.$route.query.name,
      bonus: this.$route.query.bonus,
      diceImg: require('../assets/dice/dice5.jpg'),
      diceNum: 0,
      dicePoints: [{index: 1, visible: false},
        {index: 2, visible: false},
        {index: 3, visible: false},
        {index: 4, visible: false},
        {index: 5, visible: false},
        {index: 6, visible: false}],
      dialogTitle: 'You WIN!',
      dialogButtonText: 'Take Your Bonus',
      showDialog: false
    }
  },
  computed: {
    bonus_digits: function () {
      let digits = this.bonus.toString().split('')
      let len = digits.length
      if (digits.length < 4) {
        for (let i = 0; i < 4 - len; i++) {
          digits.splice(0, 0, '0')
        }
      }
      return digits
    },
    dice_points: function () {
      for (let i = 0; i < this.dicePoints.length; i++) {
        if (this.dicePoints[i].index <= this.diceNum) {
          this.dicePoints[i].visible = true
        } else {
          this.dicePoints[i].visible = false
        }
      }
      return this.dicePoints
    },
    account: {
      get: function () {
        return this.$store.state.account.name
      },
      set: function (v) {
        this.$store.commit('UPDATE_ACCOUNT', v)
      }
    }
  },
  methods: {
    login: function () {
      tronWeb.trx.getAccount()
        .then(account => {
          console.log(account)
          account.name = tronWeb.address.fromHex(account.address)
          this.account = account
        }).catch(e => {
          console.log(e)
        })
    },
    logout: function () {
      this.account = null
    },
    is_foothold: function (position) {
      return position === this.foothold
    },
    step: function (cell) {
      if (this.diceNum === 0) {
        return
      }

      let direction = -1
      switch (cell.id - this.foothold) {
        case -1:
          direction = 0
          break
        case 1:
          direction = 1
          break
        case -7:
          direction = 2
          break
        case 7:
          direction = 3
          break
      }

      if (direction >= 0) {
        this.$http.get(backend.server + '/step?id=' + this.roomId
          + '&from=' + this.foothold + '&to=' + cell.id + "&account=" + this.account)
          .then((response) => {
            if (response.status === 200 && response.data.status === 'OK') {
              if (response.data.hitWall) {
                this.hit_wall(direction)
              } else {
                this.move_to(direction, cell.id)
                this.collect_gem()
              }
            }
          })
      }
    },
    hit_wall: function (direction) {
      this.board[this.foothold].wall = this.wallAnims[direction]
      this.board[this.foothold].showWall = true
      setTimeout(() => {
        this.board[this.foothold].showWall = false
        this.enterAnim = 'animated rotateIn delay-1s'
        this.exitAnim = 'animated tada'
        this.foothold = (this.playerId == 0 ? 48 : 0)
      }, 800)
      this.diceNum = 0
    },
    move_to: function (direction, position) {
      this.enterAnim = this.moveAnims[direction]
      this.exitAnim = 'animated fadeOut'
      this.foothold = position
      if (this.diceNum > 0) {
        this.diceNum--
      }
    },
    collect_gem: function () {
      if (this.board[this.foothold].src !== '' && this.board[this.foothold].show) {
        this.board[this.foothold].show = false
        this.players[this.playerId].gems.push(this.board[this.foothold].src)
        if (this.players[this.playerId].gems.length === 5) {
          this.diceNum = 0
          this.show_dialog('You WIN!', 'Take Your Bonus')
        }
      }
    },
    roll_dice: function () {
      this.diceImg = require('../assets/dice/dicef.jpg')
      setTimeout(() => {
        this.diceImg = require('../assets/dice/dices.jpg')
        setTimeout(() => {
          this.diceImg = require('../assets/dice/dicet.jpg')
          setTimeout(() => {
            this.diceImg = require('../assets/dice/dices.jpg')
            setTimeout(() => {
              this.diceImg = require('../assets/dice/dicef.jpg')
              setTimeout(() => {
                this.diceNum = Math.ceil(Math.random() * 6)
                this.diceImg = require('../assets/dice/dice' + this.diceNum.toString() + '.jpg')
              }, 100)
            }, 100)
          }, 100)
        }, 100)
      }, 100)
    },
    show_dialog: function (title, text) {
      this.dialogTitle = title
      this.dialogButtonText = text
      this.showDialog = true
    },
    close_dialog: function () {
      this.showDialog = false

      this.$http.get(backend.server
        + '/claim?id=' + this.roomId
        + '&account=' + this.account)
      .then((response) => {
        if (response.status === 200 && response.data.status === 'OK') {
          this.takeBonus(response.data.mazeInfo)
        } else {
          this.show_dialog(response.data.status, 'OK')
        }
      })
    },
    takeBonus: function (mazeInfo) {
      let serviceFee = 1
      tronWeb.trx.getBalance(this.account)
        .then((balance) => {
          // check balance
          if (tronWeb.fromSun(balance) >= serviceFee) {
            tronWeb.contract().at(contract.account)
              .then((con) => {
                // send transaction to create maze
                con.takeBonus(this.roomId, mazeInfo).send({
                  callValue: tronWeb.toSun(serviceFee),
                  shouldPollResponse: false
                }).then((result) => {
                  this.quit()
                }).catch(e => {
                  this.show_dialog('Fail to take bonus', 'OK');
                })
              })
            } else {
              this.show_dialog('Insufficient balance', 'OK')
            }
        }).catch(e => {
          this.show_dialog('Fail to get balance', 'OK');
        })
    },
    quit: function () {
      this.$router.push({path: '/', query:{}})
    }
  }
}
</script>
<style scoped>
.HolyGrail {
  display: flex;
  min-height: 100vh;
  flex-direction: column;
}

header {
  flex: 1;
  background-color: #fc4a1a;
}

footer {
  flex: 1;
  background-color: #dfdce3;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #f78733;
  font-size: 12px;
}

.HolyGrail-body {
  display: flex;
  flex: 10;
}

.HolyGrail-content {
  flex: 1;
  display: flex;
  min-height: 100%;
  flex-direction: row;
}

.HolyGrail-left,
.HolyGrail-right {
  flex: 0 0 12em;
  background-color: #4abdac;
}

.HolyGrail-left {
  order: -1;
}

playground {
  flex: 2;
  background-color: #ffffcc;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
}

controller {
  flex: 1;
  background-color: #f78733;
  display: flex;
  flex-direction: column;
  padding-top: 5%;
  padding-bottom: 5%;
  padding-left: 2%;
  padding-right: 2%;
}

grid {
  width: 90%;
  height: 90%;
  display: flex;
  flex-flow: row wrap;
  align-content: flex-start;
}

gridcell {
  box-sizing: border-box;
  flex: 0 0 14.2%;
  height: 14.2%;
  border: 1px solid #4c4c4c;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
}

.cell-img {
  width: 60%;
  height: auto;
}

.cell-player {
  width: 60%;
  height: 60%;
  position: absolute;
}

.cell-wall {
  width: 100%;
  height: 100%;
  position: absolute;
}

.delay-1s {
  animation-delay: 1s;
}

bonus {
  flex: 1;
  display: flex;
  flex-direction: row;
}

.tron-icon {
  flex: 0 0 25%;
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin-left: 10px;
}

.tron-icon-img {
  width: 80%;
  height: auto;
}

.total-bonus {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-right: 10px;
  position: relative;
}

.bonus-digit {
  flex: 1;
  border-radius: 5px;
  color: white;
  font-size: 40px;
  background-color: #4c4c4c;
  margin-left: 2px;
  margin-right: 2px;
  display: flex;
  justify-content: center;
  box-shadow: 0px 5px 5px #0e0c0d;
}

collection {
  flex: 3;
  display: flex;
  flex-direction: column;
}

.score {
  flex: 0 0 25%;
  display: flex;
  flex-direction: row;
  background-color: #4abdac;
  margin-top: 2%;
  margin-bottom: 2%;
  border-radius: 10px;
  box-shadow: 0px 5px 10px #0e0c0d;
}

.avatar {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  border-right: 3px solid #f78733;
}

.gem-collection {
  flex: 4;
  display: flex;
  flex-direction: row;
}

.gem-item {
  width: 20%;
  display: flex;
  justify-content: center;
  align-items: center;
}

panel {
  flex: 1;
  display: flex;
  flex-direction: row;
}

.panel-btn {
  display: flex;
  justify-content: center;
  align-items: center;
  flex: 1;
  margin-top: 5%;
  margin-bottom: 5%;
  border-radius: 10px;
  color: white;
  font-size: 25px;
  box-shadow: 0px 5px 10px #0e0c0d;
}

.roll-dice {
  margin-right: 2%;
  background-color: #4abdac;
}

.quit {
  margin-left: 2%;
  background-color: #4abdac;
}

dice-panel {
  flex: 1;
  background-color: white;
  border-radius: 10px;
  box-shadow: 0px 5px 10px #0e0c0d;
  display: flex;
  flex-direction: row;
  align-items: center;
  padding-left: 5%;
}

.dice {
  width: 23%;
  height: auto;
}

.dice-points {
  flex: 0 0 65%;
  height: 100%;
  border-left: 3px solid #f78733;
  display: flex;
  flex-direction: row;
  align-items: cell-img;
  margin-left: 15px;
  padding-left: 10px;
}

.circle {
  flex: 0 0 15%;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-right: 2px;
}

.circle-img {
  width: 100%;
  height: auto;
}
</style>

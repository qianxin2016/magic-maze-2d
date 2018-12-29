<template>
  <div class="HolyGrail">
    <link href="https://cdn.jsdelivr.net/npm/animate.css@3.5.1" rel="stylesheet" type="text/css">
    <header></header>
    <div class="HolyGrail-body">
      <div class="HolyGrail-content">
        <hall>
          <room v-for="r in rooms" :key="r.id" :style="room_bg(r)" @click="enter(r)">
            <div class="title">{{r.name}}</div>
            <div class="players">
                <img v-show="r.creator" src="../assets/player1.png" :class="player_anim(r)" />
                <img v-show="r.challenger" src="../assets/player2.png" :class="player_anim(r)" />
            </div>
            <div class="bonus">
              <img src="../assets/tron.png" class="token" />
              <div class="bonus-digits">{{r.bonus}}</div>
            </div>
          </room>
        </hall>
        <panel>
          <div class="account-panel">
            <div class="account-info">
              <img src="../assets/account.png" class="account-info-img" @click="login" />
              <span class="account-info-text">{{account_short}}</span>
            </div>
            <div class="account-info">
              <img src="../assets/balance.png" class="account-info-img" />
              <span class="account-info-text">{{balance}} TRX</span>
            </div>
          </div>
          <div class="create-panel">
            <div class="user-input">
              <label>Name:</label>
              <input type="text" v-model="roomName" />
            </div>
            <div class="user-input">
              <label>Bonus:</label>
              <input type="number" v-model="roomBonus" />
              <span class="trx-unit">TRX</span>
            </div>
            <div class="create-room" @click="create_room">Create Room</div>
          </div>
        </panel>
        </div>
      <div class="HolyGrail-left"></div>
      <div class="HolyGrail-right"></div>
    </div>
    <footer> ©Flamingo 2018. All Rights Reserved. </footer>
    <Dialog @close="close_dialog" :show="showDialog" :title="dialogTitle" :button_text="dialogButtonText" :icon="require('../assets/info.png')" />
  </div>
</template>
<script>
import { Group, Cell } from 'vux'
import Dialog from './Dialog.vue'
import contract from '../utils/contract'
import backend from '../utils/backend'

export default {
  beforeCreate () {
    this.$http.get(backend.server + '/roomlist?page=0')
      .then((response) => {
        let roomlist = response.data
        for (let id in roomlist) {
          let room = roomlist[id]
          this.rooms.splice(0, 0, {
            id: id,
            name: room.Name,
            bonus: room.Bonus,
            creator: room.Creator.Account,
            challenger: room.Challenger ? room.Challenger.Account : null
          })
        }
      })
  },
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
      dialogTitle: 'Please login TronLink',
      dialogButtonText: 'OK',
      showDialog: false,
      txid: 'fake-txid',
      roomName: '',
      roomBonus: '',
      anims: [
        'bounce',
        'tada',
        'jello',
        'rubberBand',
        'shake'
      ],
      rooms: [
        // Below fake data is just for debug purpose
        {id: 1,  name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: 'honeybadger'},
        {id: 2,  name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: 'honeybadger'},
        {id: 3,  name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: null},
        {id: 4,  name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: 'honeybadger'},
        {id: 5,  name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: 'honeybadger'},
        {id: 6,  name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: null},
        {id: 7,  name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: 'honeybadger'},
        {id: 8,  name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: 'honeybadger'},
        {id: 9,  name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: 'honeybadger'},
        {id: 10, name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: null},
        {id: 11, name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: 'honeybadger'},
        {id: 12, name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: 'honeybadger'},
        {id: 13, name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: 'honeybadger'},
        {id: 14, name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: null},
        {id: 15, name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: 'honeybadger'},
        {id: 16, name: 'fake-room', bonus: 1024, creator: 'flamingo', challenger: 'honeybadger'}
      ]
    }
  },
  computed: {
    account: {
      get: function () {
        return this.$store.state.account.name
      },
      set: function (v) {
        this.$store.commit('UPDATE_ACCOUNT', v)
      }
    },
    balance: function () {
      return this.$store.state.account.balance ? this.$store.state.account.balance / 1000000 : '0'
    },
    account_short: function () {
      return this.$store.state.account.name ? this.$store.state.account.name.substr(0, 18) + '...' : '--'
    }
  },
  methods: {
    show_dialog: function (title, text) {
      this.dialogTitle = title
      this.dialogButtonText = text
      this.showDialog = true
    },
    close_dialog: function () {
      this.showDialog = false
    },
    room_bg: function (room) {
      return room.challenger ? 'background-color: #777777' : 'background-color: #f78733'
    },
    player_anim: function (room) {
      return room.challenger ? 'player'
        : 'player animated infinite ' + this.anims[Math.floor(Math.random() * (this.anims.length - 1) + 0.5)]
    },
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
    create_room: function () {
      if (this.account === '--') {
        this.show_dialog('Please login TronLink', 'OK')
      } else if (this.roomName === '') {
        this.show_dialog('No room name', 'OK')
      } else if (this.roomBonus === '') {
        this.show_dialog('No room bonus', 'OK')
      } else if (this.roomBonus > this.balance) {
        this.show_dialog('Insufficient balance', 'OK')
      } else {
        let id = new Date().getTime()
        this.$http.get(backend.server + '/newmaze?id=' + id)
          .then((response) => {
            if (response.status === 200 && response.data.status === 'OK') {
              this.send_tx(id, response.data.mazeHash)
            }
          })
      }
    },
    send_tx: function (id, mazeHash) {
      tronWeb.trx.getBalance(this.account)
        .then((balance) => {
          // check balance
          if (tronWeb.fromSun(balance) >= this.roomBonus) {
            tronWeb.contract().at(contract.account)
              .then((con) => {
                // send transaction to create maze
                con.create(id, this.roomName, mazeHash).send({
                  callValue: tronWeb.toSun(this.roomBonus),
                  shouldPollResponse: false
                }).then((result) => {
                  // notify backend & enter room
                  this.notify_creation(id, result, mazeHash)
                }).catch(e => {
                  this.show_dialog('Fail to create room', 'OK');
                })
              })
            } else {
              this.show_dialog('Insufficient balance', 'OK')
            }
        }).catch(e => {
          this.show_dialog('Fail to get balance', 'OK');
        })
    },
    notify_creation: function (id, txid, mazeHash) {
      this.$http.get(backend.server
        + '/create?id=' + id
        + '&name=' + this.roomName
        + '&txid=' + txid
        + '&bonus='+ this.roomBonus
        + '&account=' + this.account
        + '&mazeHash=' + mazeHash)
      .then((response) => {
        // enter room
        if (response.status === 200 && response.data.status === 'OK') {
          this.$router.push({path: '/playground', query: {
            id: id,
            name: this.roomName,
            bonus: this.roomBonus,
            playerId: 0,
            account: this.account
          }})
        } else {
          this.show_dialog('Creation failed...', 'OK')
        }
      })
    },
    enter: function (room) {
      if (room.challenger) {
        this.show_dialog('This room is full', 'OK')
      } else if (room.creator === 'flamingo') {
        this.show_dialog('Just for debug', 'OK')
      } else if (this.account === room.creator) {
        this.$router.push({path: '/playground', query: {
          id: room.id,
          name: room.name,
          bonus: room.bonus,
          playerId: 0,
          account: room.creator
        }})
      } else {
        this.join(room)
      }
    },
    join: function (room) {
      tronWeb.trx.getBalance(this.account)
        .then((balance) => {
          // check balance
          if (tronWeb.fromSun(balance) >= room.bonus) {
            tronWeb.contract().at(contract.account)
              .then((con) => {
                // send transaction to create maze
                con.challenge(room.id).send({
                  callValue: tronWeb.toSun(room.bonus),
                  shouldPollResponse: false
                }).then((result) => {
                  // notify backend & enter room
                  this.notify_join(room, result)
                }).catch(e => {
                  this.show_dialog('Fail to enter room', 'OK');
                })
              })
            } else {
              this.show_dialog('Insufficient balance', 'OK')
            }
        }).catch(e => {
          this.show_dialog('Fail to get balance', 'OK');
        })
    },
    notify_join: function (room, txid) {
      this.$http.get(backend.server
        + '/join?id=' + room.id
        + '&txid=' + txid
        + '&account=' + this.account)
      .then((response) => {
        // enter room
        if (response.status === 200 && response.data.status === 'OK') {
          this.$router.push({path: '/playground', query: {
            id: room.id,
            name: room.name,
            bonus: room.bonus,
            playerId: 1,
            account: this.account
          }})
        } else {
          this.show_dialog('Failed to join...', 'OK')
        }
      })
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
  display: flex;
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

hall {
  flex: 3;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: flex-start;
  overflow-y: auto;
  overflow-x: hidden;
  background-color: #ffffcc;
  padding-left: 2%;
}

room {
  flex: 0 0 16%;
  height: 28%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border-radius: 10px;
  box-shadow: 0px 5px 10px #0e0c0d;
  background-color: #f78733;
  margin-left: 1.5%;
  margin-right: 1.5%;
  margin-top: 10px;
  margin-bottom: 10px;
  transition: all 0.2s;
}

room:hover {
  transform: scale(1.05);
}

panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
}

input {
  background-color: transparent;
  border: none;
  outline: 0px;
  border-bottom: 3px solid #d5f800;
  margin-left: 5px;
  flex: 1;
  color: white;
  font-size: 16px;
  width: 60%;
}

input::-webkit-input-placeholder {
  color: #d5f800;
}

input::-moz-input-placeholder {
  color: #d5f800;
}

input::-ms-input-placeholder {
  color: #d5f800;
}

input[type=number] {
    -moz-appearance:textfield;
}

input[type=number]::-webkit-inner-spin-button,
input[type=number]::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
}

label {
  color: #d5f800;
  font-size: 16px;
  font-family: bold;
  flex: 0 0 30%;
}

.trx-unit {
  color: #d5f800;
  font-size: 16px;
  font-family: bold;
  flex: 0 0 2%;  
}

.create-panel {
  flex: 0 0 30%;
  width: 90%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #4abdac;
  border-radius: 10px;
  box-shadow: 0px 5px 10px #0e0c0d;
  padding-top: 5%;
}

.user-input {
  flex: 1;
  width: 80%;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
}

.create-room {
  flex: 2;
  width: 70%;
  max-height: 25%;
  border-radius: 10px;
  background-color: #f78733;
  box-shadow: 0px 5px 10px #0e0c0d;
  color: white;
  font-size: 20px;
  font-family: bold;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 20px;
  margin-bottom: 20px;
}

.account-panel {
  flex: 0 0 15%;
  width: 90%;
  background-color: #4abdac;
  box-shadow: 0px 5px 10px #0e0c0d;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-start;
  margin-top: 10%;
  margin-bottom: 5%;
}

.account-info {
  flex: 1;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  margin-left: 10%;
}

.account-info-img {
  flex: 0 0 15%;
  width: 15%;
  height: auto;
}

.account-info-text {
  color: #d5f800;
  font-size: 12px;
  font-family: bold;
  flex: 1;
  margin-left: 10px;
}

.title {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: flex-end;
  color: #d5f800;
  border-bottom: 2px solid #d5f800;
  width: 80%;
}

.bonus {
  flex: 1;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 5px;
}

.token {
  flex: 0 0 20%;
  width: 20%;
  height: auto;
  margin-left: 10px;
}

.bonus-digits {
  flex: 0 0 50%;
  height: 80%;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 10px;
  color: white;
  font-size: 14px;
  background-color: #4c4c4c;
  margin-left: 5px;
  margin-right: 10px;
}

.players {
  flex: 3;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
}

.player {
  width: 30%;
  height: auto;
  margin-left: 2px;
  margin-right: 2px;

}

.walk1 {
  animation: 1s walk1 infinite alternate-reverse;
}

.walk2 {
  animation: 1s walk2 infinite alternate-reverse;
}

.delay-1s {
  animation-delay: 1s;
}

@keyframes walk1 {
  0% {
    transform: translateX(-10px);
  },
  25% {
    transform: translateX(0px);
  },
  50% {
    transform: translateX(10px);
  },
  100% {
    transform: translateX(0px);
  }
}

@keyframes walk2 {
  0% {
    transform: translateX(10px);
  },
  25% {
    transform: translateX(0px);
  },
  50% {
    transform: translateX(-10px);
  },
  100% {
    transform: translateX(0px);
  }
}

.debug {
  border: 1px solid red;
}

</style>

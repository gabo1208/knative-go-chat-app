import React from 'react'
import './App.css'
import { Dashboard } from './components/dashboard/Dashboard'
import userImg from './static/img/user.png'
import {
  FirstUserConnection,
  NewUserConnected,
  UserReconnected,
  UserDisconnected,
  GetConnectedUsers
} from './utils/Helper'
import ReconnectingWebSocket from 'reconnecting-websocket'
export default class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      menuBarClass: "change",
      username: "",
      open: false,
      sock: null,
      selectedUsername: "",
      connectedUsers: {},
      // This while I change to typescript
      // userModel: {username: string, messages: [{content: string, mine: bool}], connected: bool}
      width: 0
    }
    this.updateWindowDimensions = this.updateWindowDimensions.bind(this);
    this.handleMenuClick = this.handleMenuClick.bind(this)
  }

  usernameCallback = (username) => {
    if (username !== "" && this.state.sock) {
      console.log("sending username to sock")
      this.state.sock.send(JSON.stringify({ username: username.toLowerCase().trim() }))
    }
  }

  updateselectedUsernameCallback = (selectedUsername) => {
    this.setState(state => ({ ...state, selectedUsername }))
  }

  updateChatMessagesCallback = (message) => {
    if (message !== "" && message.length < 257) {
      let to = this.state?.selectedUsername
      let msg = JSON.stringify({ message, to })
      console.log("sending message to sock", msg)
      // update own messages
      this.state.connectedUsers[to].messages.push({ content: message, mine: true, viewed: true })
      this.setState(state => ({ ...state }), () => this.state.sock.send(msg))
    }
  }

  handleMenuClick() {
    this.setState(state => ({ ...state, menuBarClass: this.getNextChatClass() }))
  }

  getNextChatClass() {
    let menuClass = ""
    if (this.state.menuBarClass === "" || !this.state.username) {
      menuClass = "change"
    }

    return menuClass
  }

  renderAppMenu() {
    return <div
      className={`app-container ${this.state.menuBarClass}`}
      onClick={this.handleMenuClick}
    >
      <div className={`container ${this.state.username ? "cursor-pointer" : "blocked"}`}>
        <div className="bar1"></div>
        <div className="bar2"></div>
        <div className="bar3"></div>
      </div>
      <img className="header-img" alt="user-img.png" src={userImg} />
      <div className="app-title">What'sUp</div>
    </div>
  }

  renderDesktopHeader() {
    return <>
      <div className="side-bar-header">
        {this.renderAppMenu()}
      </div>
      <div className="chat-header"></div>
    </>
  }

  renderMobileHeader() {
    <div className="mobile-chat-header">
      {this.renderAppMenu()}
    </div>
  }

  componentDidMount() {
    this.updateWindowDimensions();
    window.addEventListener('resize', this.updateWindowDimensions);

    if (this.state.sock) {
      return
    }

    console.log("Protocol: " + window.location.protocol)
    let wsURL = "ws://" + document.location.host + "/ws"
    if (window.location.protocol === 'https:') {
      wsURL = "wss://" + document.location.host + "/ws"
    }

    console.log("WS URL: " + wsURL)

    let sock = new ReconnectingWebSocket(wsURL)
    sock.onopen = function () {
      console.log("connected to " + wsURL + " " + sock, this.state)
      if (this.state.username) {
        sock.send(JSON.stringify({ username: this.state.username, reconnecting: true }))
      }
    }.bind(this)

    sock.onclose = function (e) {
      console.log("connection closed (" + e.code + ")")
    }

    sock.onmessage = (e) => {
      let ev = JSON.parse(e?.data)

      if (ev.type) {
        onCloudEvent(ev)
      } else {
        this.state.connectedUsers[ev.from].messages.push({ content: ev.message })
        this.setState(state => ({ ...state }))
      }
    }


    this.setState(state => ({ ...state, sock }))

    const showError = () => (this.setState(state => ({ ...state, open: !state.open })))

    const onCloudEvent = (event) => {
      console.log(event)
      let menuBarClass = event.type === GetConnectedUsers
        ? this.state.menuBarClass
        : ''

      switch (event.type) {
        case GetConnectedUsers:
        case UserReconnected:
        case FirstUserConnection:
          this.setState(state => ({
            ...state,
            menuBarClass,
            username: event.data.username || state.username,
            connectedUsers: {
              ...state.connectedUsers,
              ...event.data.connectedUsers.reduce((acc, username) => {
                if (!state.connectedUsers[username]) {
                  acc[username] = {
                    messages: state.connectedUsers[username]?.messages || [],
                    connected: true
                  }
                }
                return acc
              }, {})
            }
          }))
          break
        case NewUserConnected:
          if (!this.state.connectedUsers[event.data]) {
            this.setState(state => ({
              ...state,
              menuBarClass: '',
              connectedUsers: {
                [event.data]: { messages: [], connected: true },
                ...state.connectedUsers
              }
            }))
          }
          break
        case UserDisconnected:
          this.setState(state => (
            state.connectedUsers[event.data].connected = false
          ))
          break
        default:
          console.log("error unnexpected event", event)
          break
      }
    }
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.updateWindowDimensions);
  }

  updateWindowDimensions() {
    this.setState({ width: window.innerWidth });
  }

  render() {
    return <div className="App">
      <header>
        <div className="header-container fade-in">
          {this.state.width < 767
            ? this.renderMobileHeader()
            : this.renderDesktopHeader()}
        </div>
      </header >
      <Dashboard
        menuBarStatus={this.state.menuBarClass !== ""}
        username={this.state.username}
        selectedUsername={this.state.selectedUsername}
        connectedUsers={this.state.connectedUsers}
        usernameCallback={this.usernameCallback}
        updateChatMessagesCallback={this.updateChatMessagesCallback}
        updateselectedUsernameCallback={this.updateselectedUsernameCallback}
      />
    </div >
  }
}

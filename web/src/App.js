import React, { useState, useCallback, useEffect } from 'react'
import './App.css'
import { Dashboard } from './components/dashboard/Dashboard'
import userImg from './static/img/user.png'
import { DetectSmallScreenWidth, NewUserConnected, UserDisconnected } from './utils/Helper'
import ReconnectingWebSocket from 'reconnecting-websocket'

var sock

function App() {
  const [appState, setAppState] = useState({
    menuBarClass: "change",
    username: "",
    open: false,
    selectedUsername: "",
    connectedUsers: {}
    // This while I change to typescript
    // userModel: {username: string, messages: [{content: string, mine: bool}]}
  })

  const usernameCallback = useCallback(username => {
    if (username !== "" && sock) {
      console.log("sending message to sock")
      sock.send(JSON.stringify({ username: username.toLowerCase().trim() }))
    }
  }, [])

  const updateselectedUsernameCallback = useCallback(selectedUsername => {
    setAppState({ ...appState, selectedUsername })
  }, [setAppState, appState])

  const updateChatMessagesCallback = useCallback(message => {
    if (message !== "" && message.length < 257) {
      let to = appState.connectedUsers[appState.selectedUsername].username
      console.log("sending message to sock")
      // update own messages
      appState.connectedUsers[appState.selectedUsername].messages.push({ content: message, mine: true })
      if (to === appState.username) {
        appState.connectedUsers[appState.selectedUsername].messages.push({ content: message })
      } else {
        // send message to the other user
        sock.send(JSON.stringify({ message, to }))
      }

      setAppState({ ...appState })
    }
  }, [appState])

  const handleMenuClick = () => {
    setAppState({ ...appState, menuBarClass: getNextChatClass() })
  }

  const getNextChatClass = () => {
    let menuClass = ""
    if (appState.menuBarClass === "" || !appState.username) {
      menuClass = "change"
    }

    return menuClass
  }

  const renderAppMenu = () => (
    <div
      className={`app-container ${appState.menuBarClass}`}
      onClick={handleMenuClick}
    >
      <div className={`container ${appState.username ? "cursor-pointer" : "blocked"}`}>
        <div className="bar1"></div>
        <div className="bar2"></div>
        <div className="bar3"></div>
      </div>
      <img className="header-img" alt="user-img.png" src={userImg} />
      <div className="app-title">What'sUp</div>
    </div>
  )

  const renderDesktopHeader = () => (
    <>
      <div className="side-bar-header">
        {renderAppMenu()}
      </div>
      <div className="chat-header"></div>
    </>
  )

  const renderMobileHeader = () => (
    <div className="mobile-chat-header">
      {renderAppMenu()}
    </div>
  )

  useEffect(() => {
    if (sock) {
      return
    }

    console.log("Protocol: " + window.location.protocol)
    let wsURL = "ws://" + document.location.host + "/ws"
    if (window.location.protocol === 'https:') {
      wsURL = "wss://" + document.location.host + "/ws"
    }

    console.log("WS URL: " + wsURL)

    let ws = new ReconnectingWebSocket(wsURL)
    sock = ws
    sock.onopen = function () {
      console.log("connected to " + wsURL + " " + sock)
    }

    sock.onclose = function (e) {
      console.log("connection closed (" + e.code + ")")
    }

    sock.onmessage = (e) => {
      let ev = JSON.parse(e?.data)

      if (ev.type) {
        onCloudEvent(ev)
      } else {
        appState.connectedUsers[ev.from].messages.push({ content: ev.message })
        setAppState({ ...appState })
      }
    }

    const showError = () => (this.setState({ open: !appState.open }))

    const onCloudEvent = (event) => {
      switch (event.type) {
        case NewUserConnected:
          setAppState({
            ...appState,
            menuBarClass: '',
            username: event.data.username,
            connectedUsers: {
              ...event.data.connectedUsers.reduce((acc, username) => {
                acc[username] = { username, messages: [] }
                return acc
              }, {})
            }
          })
          break
        case UserDisconnected:
          delete appState.connectedUsers[event.data.username]
          setAppState({
            ...appState,
            connectedUsers: { ...appState.connectedUsers }
          })
          break
        default:
          console.log("error unnexpected event")
          break
      }
    }
  })

  return (
    <div className="App">
      <header>
        <div className="header-container fade-in">
          {DetectSmallScreenWidth()
            ? renderMobileHeader()
            : renderDesktopHeader()}
        </div>
      </header >
      <Dashboard
        menuBarStatus={appState.menuBarClass !== ""}
        username={appState.username}
        selectedUsername={appState.selectedUsername}
        connectedUsers={appState.connectedUsers}
        usernameCallback={usernameCallback}
        updateChatMessagesCallback={updateChatMessagesCallback}
        updateselectedUsernameCallback={updateselectedUsernameCallback}
      />
    </div >
  )
}

export default App

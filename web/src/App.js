import React, { useState, useCallback } from 'react'
import './App.css'
import { Dashboard } from './components/dashboard/Dashboard'
import userImg from './static/img/user.png'
import { detectSmallScreen } from './utils/Helper'

function App() {
  const [appState, setAppState] = useState({
    menuBarClass: "change",
    username: ""
  })

  const usernameCallback = useCallback(username => {
    if (username !== "") {
      setAppState({ menuBarClass: '', username })
    }
  }, [setAppState])


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

  const renderDesktopHeader = () => (
    <div className="header-container">
      <div className="side-bar-header">
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
      </div>
      <div className="chat-header"></div>
    </div>
  )

  const renderMobileHeader = () => (
    <div className="mobile-chat-header"></div>
  )
  return (
    <div className="App">
      <header>
        {detectSmallScreen()
          ? renderMobileHeader()
          : renderDesktopHeader()}
      </header >
      <Dashboard
        menuBarStatus={appState.menuBarClass !== ""}
        username={appState.username}
        usernameCallback={usernameCallback}
      />
    </div >
  )
}

export default App

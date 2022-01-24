import React, { useState, useCallback } from 'react'
import './Dashboard.css'
import { Contacts } from './contacts/Contacts'
import { Chat } from './chat/Chat'
import { Menu } from './menu/Menu'
import { DetectSmallScreenWidth } from '../../utils/Helper'

export function Dashboard(props) {
  const [chatState, setChatState] = useState({
    selectedChatIndex: -1,
    connectedUsers: []
    // This while I implement typescript
    // userModel: {username: string, messages: {content: string, mine: bool}}
  })

  const updateChatMessagesCallback = useCallback(message => {
    let connectedUsersAux = chatState.connectedUsers
    connectedUsersAux[chatState.selectedChatIndex].message.append({ mine: true, message })
    setChatState({
      selectedChatIndex: chatState.selectedChatIndex,
      connectedUsers: connectedUsersAux
    })
  }, [setChatState, chatState])

  const updateSelectedChatIndexCallback = useCallback(selectedChatIndex => {
    setChatState({ ...chatState, selectedChatIndex })
  }, [setChatState, chatState])

  const renderMenuOrContacts = () => {
    return <div className="fade-in">
      {props.menuBarStatus
        ? <Menu
          usernameSetter={props.usernameCallback}
          username={props.username}
        />
        : <Contacts
          username={props.username}
          connectedUsers={chatState.connectedUsers}
          selectedChatIndex={chatState.selectedChatIndex}
          selectedChat={chatState.connectedUsers[chatState.selectedChatIndex]}
          updateSelectedChatIndexCallback={updateSelectedChatIndexCallback}
        />}
    </div>
  }

  const renderDesktopDashboard = () => (
    <>
      <div className="side-bar fade-in">
        {renderMenuOrContacts()}
      </div>
      <div className="chat fade-in">
        <Chat
          username={props.username}
          selectedChat={chatState.connectedUsers[chatState.selectedChat]}
          connectedUsersNumber={chatState.connectedUsers.length}
          updateChatMessagesCallback={updateChatMessagesCallback}
        />
      </div>
    </>
  )

  const renderMobileDashboard = () => {
    return (<div className="mobile-dashboard">
      {chatState.selectedChatIndex < 0
        ? renderMenuOrContacts()
        : <Chat
          username={props.username}
          selectedChat={chatState.connectedUsers[chatState.selectedChat]}
          connectedUsersNumber={chatState.connectedUsers.length}
          updateChatMessagesCallback={updateChatMessagesCallback}
        />}
    </div>)
  }

  return (
    <div className="dashboard">
      <div className="app-container bg">
        {DetectSmallScreenWidth()
          ? renderMobileDashboard()
          : renderDesktopDashboard()}
      </div>
    </div>
  )
}
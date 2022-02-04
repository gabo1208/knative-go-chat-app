import React from 'react'
import './Dashboard.css'
import { Contacts } from './contacts/Contacts'
import { Chat } from './chat/Chat'
import { Menu } from './menu/Menu'
import { DetectSmallScreenWidth } from '../../utils/Helper'

export function Dashboard(props) {
  const renderMenuOrContacts = () => {
    return <div className="fade-in">
      {props.menuBarStatus
        ? <Menu
          usernameSetter={props.usernameCallback}
          username={props.username}
        />
        : <Contacts
          username={props.username}
          connectedUsers={props.connectedUsers}
          selectedUsername={props.selectedUsername}
          updateselectedUsernameCallback={props.updateselectedUsernameCallback}
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
          selectedUsername={props.selectedUsername}
          connectedUsers={props.connectedUsers}
          updateChatMessagesCallback={props.updateChatMessagesCallback}
          selectedChat={props.connectedUsers[props.selectedUsername]}
        />
      </div>
    </>
  )

  const renderMobileDashboard = () => {
    return (<div className="mobile-dashboard">
      {!props.selectedUsername
        ? renderMenuOrContacts()
        : <Chat
          username={props.username}
          selectedUsername={props.selectedUsername}
          connectedUsers={props.connectedUsers}
          updateChatMessagesCallback={props.updateChatMessagesCallback}
          selectedChat={props.connectedUsers[props.selectedUsername]}
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
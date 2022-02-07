import React, { useState } from 'react'
import sendBtn from '../../../static/img/send-btn.png'
import './Chat.css'

export function Chat(props) {
  const [messageInput, setMessageInput] = useState('')

  const chatInputHandler = (e) => {
    setMessageInput(e.target.value)
  }

  const handleChatInputKeyDown = (e) => {
    if (e.key === 'Enter') {
      sendButton()
    }
  }

  const sendButton = () => {
    if (messageInput !== '') {
      props.updateChatMessagesCallback(messageInput)
      setMessageInput('')
    }
  }

  const chatDefaultMsg = () => {
    if (!props.username) {
      return "Please choose an username to start chatting!"
    }

    let connectedUsersNumber = Object.keys(props.connectedUsers)?.length

    if (!connectedUsersNumber) {
      return "Wait for someone to come online and start chatting :)!"
    }

    let userChatMsg = props.selectedUsername === props.username
      ? "Here is a space where you can talk to yourself and 'no one' is going to jugde you!"
      : `Start chatting with ${props.selectedUsername}!`
    return props.selectedChat ? userChatMsg
      : "Click on one of your contacts and start chatting!"
  }

  const renderMessages = () => {
    return props.selectedChat?.messages?.map((message, i) => {
      return <div key={i} className={"fade-in " +
        (message.mine ? "user" : "contact") + "-message"
      }>
        {message.content}
      </div>
    })
  }

  let connectedUsersNumber = Object.keys(props.connectedUsers)?.length
  return (
    <div className="chat-column">
      <div className="messages">
        <div className="messages-bubble">
          <div className="bubbles-bg">
            {props.selectedChat?.messages.length
              ? renderMessages()
              : <div className="empty-chat">
                <h3 className="chat-default-msg">{chatDefaultMsg()}</h3>
              </div>}
          </div>
        </div>
      </div>
      <div
        className={
          `chat-section ${(!connectedUsersNumber || !props.selectedUsername) && "blocked"}`
        }
      >
        <div className="chat-bubble">
          <div className="bubbles-bg no-overflow">
            <img
              src={sendBtn}
              className={
                `send-btn ${connectedUsersNumber && props.selectedUsername && "cursor-pointer"}
              `}
              alt="send-btn.png"
              onClick={sendButton}
            />
            <input
              type="text"
              name="message"
              className="chat-input"
              placeholder="Type a message"
              value={messageInput}
              onChange={chatInputHandler}
              onKeyDown={handleChatInputKeyDown}
              maxLength="256"
            />
          </div>
        </div>
      </div>
    </div>
  )
}
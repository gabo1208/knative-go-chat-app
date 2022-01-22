import React, { useState } from 'react'
import usernameIcon from '../../../static/img/username-icon.png'
import logoutIcon from '../../../static/img/logout-icon.png'
import './Menu.css'

export function Menu(props) {
  const [username, setUsername] = useState(props.username)

  const handleUsernameInputKeyDown = (e) => {
    if (e.key === 'Enter') {
      props.usernameSetter(username)
    }
  }

  const usernameInputHandler = (e) => {
    setUsername(e.target.value)
  }

  return (
    <div className="menu-column">
      <div className="option-card">
        <div className="option-img-wrapper">
          <img src={usernameIcon} className="option-img" alt="option-img.png" />
        </div>
        <div className="option-info">
          <div className="option">
            {props.username
              ? <h3 className="option-name">
                <b>Hey {props.username}</b>
              </h3>
              : <input
                type="text"
                name="username"
                className="chat-input menu"
                placeholder="Type your username"
                onChange={usernameInputHandler}
                value={username}
                onKeyDown={handleUsernameInputKeyDown}
              />
            }
          </div>
        </div>
      </div>
      <div className="option-card cursor-pointer">
        <div className="option-img-wrapper">
          <img src={logoutIcon} className="option-img logout" alt="option-img.png" />
        </div>
        <div className="option-info">
          <div className="option">
            <h3 className="option-name"><b>Logout</b></h3>
          </div>
        </div>
      </div>
    </div>
  )
}
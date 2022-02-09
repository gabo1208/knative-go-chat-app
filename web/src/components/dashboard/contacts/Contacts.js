import React, { useState } from 'react'
import './Contacts.css'
import userImg from '../../../static/img/user.png'

export function Contacts(props) {
  const [contactsFilter, setContactsFilter] = useState('')

  const contactsFilterInputHandler = (e) => {
    setContactsFilter(e.target.value)
  }

  const renderContacts = (contactsList) => {
    let contactsListKeys = Object.keys(contactsList)
    return contactsListKeys?.length
      ? contactsListKeys.filter(key =>
        !contactsFilter || key.includes(contactsFilter)
      ).map((key, i) => {
        return (
          <div
            key={i}
            className={
              `contact-card cursor-pointer ${key === props.selectedUsername && "selected"}`
            }
            onClick={() => props.updateselectedUsernameCallback(key)}
          >
            <div className="contact-img">
              <img className="user-img" src={userImg} alt="user-img.png" />
            </div>
            <div className="contact-info">
              <div className="contact-name">
                <h3 className="user-name capitalized"><b>
                  {key !== props.username
                    ? key
                    : "You"}
                </b></h3>
              </div>
              <div className="contact-last-seen">
                Connected
              </div>
              <div className="contact-last-msg">
                {contactsList[key].messages?.length
                  ? contactsList[key].messages[contactsList[key].messages.length - 1].content
                  : "Start Chatting with " + key}
              </div>
            </div>
          </div>
        )
      })
      : <div className="no-contacts-msg">Connected Users: 0</div>
  }

  return (
    <div className="contacts-column">
      <div className="filter">
        <input
          type="text"
          name="contacts-filter"
          className="chat-input contacts-filter"
          placeholder="Search for a Contact"
          onChange={contactsFilterInputHandler}
          value={contactsFilter}
        />
      </div>
      {renderContacts(props.connectedUsers)}
    </div>
  )
}

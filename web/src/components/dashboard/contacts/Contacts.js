import React, { useState } from 'react'
import './Contacts.css'
import userImg from '../../../static/img/user.png'

export function Contacts(props) {
  const [contactsFilter, setContactsFilter] = useState('')

  const contactsFilterInputHandler = (e) => {
    setContactsFilter(e.target.value)
  }

  const renderContacts = (contactsList) => {
    return contactsList.length
      ? contactsList.filter(contact =>
        !contactsFilter || contact.username.includes(contactsFilter)
      ).map((contact, i) => {
        return (
          <div
            className={
              `contact-card cursor-pointer ${i === props.selectedChat && "selected"}`
            }
          >
            <div className="contact-img">
              <img className="user-img" src={userImg} alt="user-img.png" />
            </div>
            <div className="contact-info">
              <div className="contact-name">
                <h3 className="user-name"><b>{contact.username}</b></h3>
              </div>
              <div className="contact-last-seen">
                Connected
              </div>
              <div className="contact-last-msg">
                {contact.messages[-1]}
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
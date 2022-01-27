import React, { useState, useCallback, useEffect } from 'react'
import './Dashboard.css'
import { Contacts } from './contacts/Contacts'
import { Chat } from './chat/Chat'
import { Menu } from './menu/Menu'
import { DetectSmallScreenWidth } from '../../utils/Helper'
import ReconnectingWebSocket from 'reconnecting-websocket';

export function Dashboard(props) {
  const [chatState, setChatState] = useState({
    selectedChatIndex: -1,
    connectedUsers: [],
    open: false
    // This while I change to typescript
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

  const ConnectToWebSocket = useEffect(() => {
    console.log("Protocol: " + window.location.protocol);
    let wsURL = "ws://" + document.location.host + "/ws";
    if (window.location.protocol === 'https:') {
      wsURL = "wss://" + document.location.host + "/ws";
    }

    console.log("WS URL: " + wsURL);

    let sock = new ReconnectingWebSocket(wsURL);
    sock.onopen = function () {
      console.log("connected to " + wsURL);
      //let fab = document.getElementById("fab");
      //fab.setAttribute("sockeye-connected", "true");
    };
    sock.onclose = function (e) {
      console.log("connection closed (" + e.code + ")");
      //fab.setAttribute("sockeye-connected", "false");
    };
    sock.onmessage = function (e) {
      window.dispatchEvent(new Event('cloudevent'));
      let t = JSON.parse(JSON.parse(e.data)); // at the moment the ws sends down a double encoded thing.

      console.log(t)
      onCloudEvent(t)
    }
  })

  const showError = () => (this.setState({ open: !this.state.open }))

  const onCloudEvent = (event) => {
    let data = { id: event.id };

    Object.keys(event).forEach(key => {
      if (key === "data") {
        data[key] = JSON.stringify(event[key]);
        return;
      }
      data[key] = event[key];
    });

    let al = [...this.state.events];

    if (this.state.revert) {
      if (data["data"] != null) {
        al.push(data);
        this.setState({
          events: al.reverse()
        });
        return;
      }

      if (data["data"] === undefined) {
        showError();
        console.log("More information on the invalid event: ", event);
        return;
      }
    }

    if (!this.state.revert) {
      if (data["data"] != null) {
        al.push(data);
        this.setState({
          events: al
        });
        return;
      }

      if (data["data"] === undefined) {
        showError();
        console.log("More information on the invalid event: ", event);
        return;
      }
    }
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
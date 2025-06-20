import * as apiClient from "../api-client";
import { useEffect, useState } from "react";
import { useAppContext } from "../contexts/AppContext";
import { useDispatch, useSelector } from "react-redux";
import { useWebSocket } from "../contexts/WebSocketContext";
import "./Messanger.scss";
import { FaTimes, FaTelegramPlane, FaPlus, FaMinus } from "react-icons/fa";
import { RootState } from "../contexts/store";
import CanvasAvatar from "../components/CanvasAvatar";
import {UserMessage} from "..//contexts/slices/userMessagesSlice";


const NewChat = ({ setOpenNewChat }) => {
  const socket = useWebSocket();

  const [emailOrNickname, setEmailOrNickname] = useState("");
  const [startMessage, setStartMessage] = useState("");
  const sendFirstMessage = () => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      // socket.send(JSON.stringify({ type: "ping" }));
      socket.send(
        JSON.stringify({
          type: "initial_message",
          to: emailOrNickname,
          body: startMessage,
          msgType: "string",
        })
      );
    }
    setOpenNewChat(false);
  };
  return (
    <div className="new-chat">
      <div className="new-chat-container">
        <input
          placeholder="Create new chat with (email \ nickname)"
          value={emailOrNickname}
          onChange={(e) => setEmailOrNickname(e.target.value)}
        />
        <div className="new-chat-start">
          <input
            placeholder="Start your conversation"
            value={startMessage}
            onChange={(e) => setStartMessage(e.target.value)}
          />
          <button onClick={() => sendFirstMessage()}>
            <FaTelegramPlane />
          </button>
        </div>
      </div>
    </div>
  );
};

const Chat = ({ activeToUserId, selectedChat }) => {
  const socket = useWebSocket();
  const user = useSelector((state: RootState) => state.auth.user);
  const [messageInputContent, setMessageInputContent] = useState("");
  useEffect(() => {
    console.log(1111)
  }, [selectedChat])
  const sendMessage = () => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(
        JSON.stringify({
          type: "incoming_message",
          to: activeToUserId,
          body: messageInputContent,
          msgType: "string",
        })
      );
      setMessageInputContent("")
    }
  };

  return (
    <div className="chat-combiner">
      <div className="chat-body">
        <div className="chat-body-scroll">
          {selectedChat.map((mess) => (
            <div
              className={`message-body ${
                mess.MessageSender === user.id
                  ? "right-message"
                  : "left-message"
              }`}
              key={mess.ID}
            >
              {mess.ID + mess.MessageBody}
            </div>
          ))}
        </div>
      </div>
      <div className="chat-control">
        <div className="chat-control-metadata"></div>
        <div className="chat-control-input">
          <input
            placeholder="Message"
            value={messageInputContent}
            onChange={(e) => setMessageInputContent(e.target.value)}
            type="text"
          />
          <button onClick={() => sendMessage()}>
            <FaTelegramPlane />
          </button>
        </div>
      </div>
    </div>
  );
};

const Messanger = () => {
  const socket = useWebSocket();
  const users = useSelector(
    (state: RootState) => state.usermessagesusers.groupedUsers
  );
  const messages = useSelector(
    (state: RootState) => state.usermessages.groupedMessages
  );
  const [displayUsersArray, setDisplayUsersArray] = useState<number[]>([]);
  useEffect(() => {
    console.log(222333, messages)
    function getSortedMessageGroupIds(messages: any) {
      const messageGroups = Object.entries(messages);

      const groupsWithTime = messageGroups.map(([id, msgs]: any) => {
        const lastMessage = msgs[msgs.length - 1]; 
        return [parseInt(id), new Date(lastMessage.CreatedAt).getTime()];
      });
      // const userMessages: any = messages[userId] || [];

      groupsWithTime.sort((a, b) => b[1] - a[1]);
      return groupsWithTime.map(([id]) => id);
    }
    setSelectedChat(messages[activeToUserId])
    setDisplayUsersArray(getSortedMessageGroupIds(messages));
  }, [users, messages]);
  const [openNewChat, setOpenNewChat] = useState(false);
  const [selectedChat, setSelectedChat] = useState<UserMessage[]>(false);
  const [activeToUserId, setActiveToUserId] = useState(0);
  const deleteConversation = (userID:number) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(
        JSON.stringify({
          type: "delete_conversation",
          user_id: userID,
        })
      );
    }
  };

  const [userSearch, setUserSearch] = useState("");
  return (
    <div className="messanger">
      <div className="messanger-sidebar">
        <div className="messanger-sidebar-action">
          <input
            value={userSearch}
            onChange={(e) => setUserSearch(e.target.value)}
            type="text"
          />
          <button onClick={() => setOpenNewChat(!openNewChat)}>
            {openNewChat ? <FaMinus /> : <FaPlus />}
          </button>
        </div>
        <div className="messanger-sidebar-chats">
          {displayUsersArray
            .filter((userId) => {
              const user: any = users.find((u) => u.ID === userId);
              if (!user) return false; // skip if no user found

              const searchLower = userSearch.toLowerCase();
              return (
                user.Nickname.toLowerCase().includes(searchLower) ||
                user.Email.toLowerCase().includes(searchLower)
              );
            })
            .map((userId) => {
              // Find user by id
              const user: any = users.find((u) => u.ID === userId);
              if (!user) return null; // safety check
              // Get last message from messages[userId]
              const userMessages: any = messages[userId] || [];
              const lastMessage =
                userMessages.length > 0
                  ? userMessages[userMessages.length - 1]
                  : null;

              return (
                <div
                  key={userId}
                  onClick={() => {
                    setSelectedChat(userMessages);
                    setOpenNewChat(false);
                    setActiveToUserId(userId);
                  }}
                  className={
                    "chat-item " +
                    (userId === activeToUserId ? "active-chat" : "")
                  }
                >
                  <CanvasAvatar avatar={user.Avatar} />

                  <div className="chat-item-body">
                    <div className="chat-user-info">
                      <strong>{user.Nickname}</strong> ({user.Email})
                    </div>
                    <div className="chat-last-message">
                      {lastMessage
                        ? lastMessage.MessageBody
                        : "No messages yet"}
                    </div>
                  </div>
                  <button onClick={() => deleteConversation(user.ID)}><FaTimes/></button>
                </div>
              );
            })}
        </div>
      </div>
      {openNewChat ? (
        <NewChat setOpenNewChat={setOpenNewChat} />
      ) : selectedChat ? (
        <Chat
          activeToUserId={activeToUserId}
          selectedChat={selectedChat}
        ></Chat>
      ) : (
        <></>
      )}
    </div>
  );
};

export default Messanger;
